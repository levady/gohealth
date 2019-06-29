package sitehealthchecker

import (
	"net/http"
	"runtime"
	"time"

	"github.com/levady/gohealth/internal/platform/sitestore"
)

var siteChecker = checkSiteWithTimeout

// SerialHealthChecks run health checks on all stored Sites in serial
func SerialHealthChecks(store *sitestore.Store, timeout time.Duration) {
	for _, s := range store.List() {
		resp, err := siteChecker(s.URL, timeout)
		if err != nil || resp.StatusCode != 200 {
			store.UpdateHealth(s.ID, false)
		} else {
			store.UpdateHealth(s.ID, true)
		}
	}
}

// ParallelHealthChecks run health checks on all stored Sites in parallel
func ParallelHealthChecks(store *sitestore.Store, timeout time.Duration) {
	sites := store.List()
	sitesLen := len(sites)
	resultCh := make(chan bool, len(sites))

	grs := runtime.NumCPU()
	batchCh := make(chan bool, grs)

	for idx := 0; idx < sitesLen; idx++ {
		go func(i int) {
			batchCh <- true
			{
				s := sites[i]
				resp, err := siteChecker(s.URL, timeout)
				if err != nil || resp.StatusCode != 200 {
					store.UpdateHealth(s.ID, false)
				} else {
					store.UpdateHealth(s.ID, true)
				}
				resultCh <- true
			}
			<-batchCh
		}(idx)
	}

	for sitesLen > 0 {
		<-resultCh
		sitesLen--
	}
}

func checkSiteWithTimeout(url string, timeout time.Duration) (*http.Response, error) {
	client := http.Client{Timeout: timeout}
	return client.Get(url)
}
