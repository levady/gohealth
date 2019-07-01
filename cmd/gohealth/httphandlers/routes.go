package httphandlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/levady/gohealth/internal/platform/sitestore"
	"github.com/levady/gohealth/internal/platform/sse"
)

// Middleware is the base type for all handlers
type Middleware struct {
	logger *log.Logger
}

// Routes return application routes handlers
func Routes(logger *log.Logger, str *sitestore.Store, broker *sse.Broker, sse bool) http.Handler {
	router := http.DefaultServeMux

	shh := SiteHealthHandler{SiteStore: str, SSE: sse}
	router.HandleFunc("/", shh.Homepage)
	router.HandleFunc("/sites/save", shh.Save)
	router.HandleFunc("/ajax/sites/check", shh.HealthChecks)
	router.HandleFunc("/ajax/sites/delete/", shh.Delete)

	if sse {
		router.HandleFunc("/sse", broker.SSE)
	}

	mw := Middleware{logger: logger}

	return mw.logging(router)
}

func (m *Middleware) logging(hdlr http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			m.logger.Println(requestID(), r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), time.Since(start))
		}(time.Now())

		hdlr.ServeHTTP(w, r)
	})
}

func requestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
