package sitehealthchecker

import (
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// Site represents Site data
type Site struct {
	URL     string
	Healthy *bool
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
		return errors.Wrapf(err, "Site URL is not valid")
	} else if u.Scheme == "" || u.Host == "" {
		return errors.New("Site URL must be an absolute URL")
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("Site URL must begin with http or https")
	}

	shc.Sites = append(shc.Sites, &s)

	return nil
}
