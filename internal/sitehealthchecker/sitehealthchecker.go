package sitehealthchecker

import (
	"errors"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"
)

// Site represents Site data
type Site struct {
	URL     string      `json:"url"`
	Healthy interface{} `json:"healthy"`
}

// SiteHealthChecker represents SiteHealthChecker service
type SiteHealthChecker struct {
	Sites   []*Site
	Timeout time.Duration
}

var mutex sync.Mutex

// New creates a new SiteHealthChecker
func New(timeout time.Duration) SiteHealthChecker {
	return SiteHealthChecker{
		Sites:   make([]*Site, 0),
		Timeout: timeout,
	}
}

// AddSite add a single site to the Sites slice
func (shc *SiteHealthChecker) AddSite(s Site) error {
	u, err := url.Parse(s.URL)
	if err != nil {
		return errors.New("Site URL is not valid")
	} else if u.Scheme == "" || u.Host == "" {
		return errors.New("Site URL must be an absolute URL")
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("Site URL must begin with http or https")
	}

	mutex.Lock()
	{
		shc.Sites = append(shc.Sites, &s)
	}
	mutex.Unlock()

	return nil
}

var siteChecker = checkSiteWithTimeout

// SerialHealthChecks run health checks on all stored Sites in serial
func (shc *SiteHealthChecker) SerialHealthChecks() {
	for idx := range shc.Sites {
		s := shc.Sites[idx]
		resp, err := siteChecker(s.URL, shc.Timeout)
		if err != nil || resp.StatusCode != 200 {
			s.Healthy = false
		} else {
			s.Healthy = true
		}
	}
}

// ParallelHealthChecks run health checks on all stored Sites in parallel
func (shc *SiteHealthChecker) ParallelHealthChecks() {
	sitesLen := len(shc.Sites)
	resultCh := make(chan bool, sitesLen)

	grs := runtime.NumCPU()
	batchCh := make(chan bool, grs)

	for idx := 0; idx < sitesLen; idx++ {
		go func(i int) {
			batchCh <- true
			{
				s := shc.Sites[i]
				resp, err := siteChecker(s.URL, shc.Timeout)
				mutex.Lock()
				{
					if err != nil || resp.StatusCode != 200 {
						s.Healthy = false
					} else {
						s.Healthy = true
					}
				}
				mutex.Unlock()
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
