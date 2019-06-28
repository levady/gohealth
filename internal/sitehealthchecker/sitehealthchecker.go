package sitehealthchecker

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

// Site represents Site data
type Site struct {
	URL     string
	Healthy interface{}
}

// SiteHealthChecker represents SiteHealthChecker service
type SiteHealthChecker struct {
	Sites   []*Site
	Timeout time.Duration
}

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

	shc.Sites = append(shc.Sites, &s)

	return nil
}

var siteChecker = checkSiteWithTimeout

// RunHealthChecks run health checks on all stored Sites
func (shc *SiteHealthChecker) RunHealthChecks() {
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

func checkSiteWithTimeout(url string, timeout time.Duration) (*http.Response, error) {
	client := http.Client{Timeout: timeout}
	return client.Get(url)
}
