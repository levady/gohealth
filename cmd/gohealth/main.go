package main

import (
	"context"
	"encoding/json"
	"errors"
	"expvar"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/levady/gohealth/cmd/gohealth/httphandlers"
	"github.com/levady/gohealth/internal/platform/sitestore"
	"github.com/levady/gohealth/internal/platform/sse"
	"github.com/levady/gohealth/internal/sitehealthchecker"
)

type config struct {
	Host           string
	LookbackPeriod int
	SSE            bool
}

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {
	// =========================================================================
	// Logging

	log := log.New(os.Stdout, "GO HEALTH : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// Parse Configuration
	host := os.Getenv("HOST")
	if host == "" {
		host = ":8080"
	}

	lpCfg := 0
	if lp := os.Getenv("LOOKBACK_PERIOD"); lp != "" {
		var err error
		lpCfg, err = strconv.Atoi(lp)
		if err != nil {
			return errors.New("main : Failed parsing LOOKBACK_PERIOD config")
		}
	}

	sseCfg := false
	if sse := os.Getenv("SSE"); sse != "" {
		var err error
		sseCfg, err = strconv.ParseBool(sse)
		if err != nil {
			return errors.New("main : Failed parsing SSE config")
		}
	}

	cfg := config{
		Host:           host,
		LookbackPeriod: lpCfg,
		SSE:            sseCfg,
	}

	prettyCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Printf("main: Failed marshal indent app config %v", err)
	}
	log.Printf("main : Config :%v", string(prettyCfg))

	// =========================================================================
	// Initializaing site memory store

	log.Printf("main : Initializing site memory store")
	str := sitestore.NewStore()

	// =========================================================================
	// App Starting

	// Print the build version for our logs.
	expvar.NewString("build").Set(build)
	log.Printf("main : Started : Application initializing : version %q", build)
	defer log.Println("main : Completed")

	// Consruct a broker server
	broker := sse.NewServer(log)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	server := http.Server{
		Addr:    cfg.Host,
		Handler: httphandlers.Routes(log, &str, broker, cfg.SSE),
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : App listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Run a ticker that will check the health of all sites every 15 seconds
	ticker := time.NewTicker(15 * time.Second)

	go func() {
		log.Printf("main : Site health checker running")
		for {
			<-ticker.C
			log.Printf("main : ticker : Run health checks")
			sitehealthchecker.ParallelHealthChecks(&str, 800*time.Millisecond, cfg.LookbackPeriod)
			broker.Notifier <- []byte("done")
		}
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return err

	case sig := <-shutdown:
		log.Printf("main : %v : Shuttting down site health checker", sig)
		ticker.Stop()

		log.Printf("main : %v : Shuttting down app", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", 5*time.Second, err)
			err = server.Close()
		}

		if err != nil {
			return err
		}
	}

	return nil
}
