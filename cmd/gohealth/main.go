package main

import (
	"context"
	"expvar"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/levady/gohealth/cmd/gohealth/httphandlers"
	"github.com/levady/gohealth/internal/platform/sitestore"
	"github.com/levady/gohealth/internal/sitehealthchecker"
)

type config struct {
	Host string
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

	cfg := config{
		Host: host,
	}

	// =========================================================================
	// Initializaing site memory store store

	log.Printf("main : Initializing site memory store")
	str := sitestore.NewStore()

	// =========================================================================
	// App Starting

	// Print the build version for our logs.
	expvar.NewString("build").Set(build)
	log.Printf("main : Started : Application initializing : version %q", build)
	defer log.Println("main : Completed")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	app := http.Server{
		Addr:    cfg.Host,
		Handler: handler(log, &str),
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : App listening on %s", app.Addr)
		serverErrors <- app.ListenAndServe()
	}()

	// Run a ticker that will check the health of all sites every 15 seconds
	ticker := time.NewTicker(15 * time.Second)

	go func() {
		log.Printf("main : Site health checker running")
		for {
			<-ticker.C
			log.Printf("main : ticker : Run health checks")
			sitehealthchecker.ParallelHealthChecks(&str, 800*time.Millisecond)
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
		err := app.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", 5*time.Second, err)
			err = app.Close()
		}

		// TODO: SIGSTOP does not work on Windows
		// Log the status of this shutdown.
		// switch {
		// case sig == syscall.SIGSTOP:
		// 	return errors.New("integrity issue caused shutdown")
		// case err != nil:
		// 	return err
		// }
	}

	return nil
}

func handler(logger *log.Logger, str *sitestore.Store) http.Handler {
	handler := http.DefaultServeMux

	shh := httphandlers.SiteHealthHandler{SiteStore: str}

	handler.HandleFunc("/", shh.Homepage)
	handler.HandleFunc("/sites/save", shh.Save)
	handler.HandleFunc("/ajax/sites/check", shh.HealthChecks)
	handler.HandleFunc("/ajax/sites/delete/", shh.Delete)

	return handler
}
