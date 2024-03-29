package sitestore

import (
	"errors"
	"net/url"
	"sort"
	"sync"
	"time"
)

const (
	// Unknown indicate the that site status is unknown
	Unknown = iota
	// Healthy indicate that the site is healthy
	Healthy
	// Unhealthy indicate that the site is not healthy
	Unhealthy
)

// Site represents Site data
type Site struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Status    int       `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store represent data store for sites
type Store struct {
	sites     map[int]*Site
	idTracker int
	sync.RWMutex
}

// NewStore construct a new Store
func NewStore() Store {
	return Store{
		sites:     make(map[int]*Site),
		idTracker: 0,
	}
}

// List returns a collection of sites
func (str *Store) List() []Site {
	str.RLock()
	defer str.RUnlock()

	sites := make([]Site, 0)
	for _, site := range str.sites {
		sites = append(sites, *site)
	}

	sort.Slice(sites, func(i, j int) bool {
		return sites[i].ID < sites[j].ID
	})

	return sites
}

// ListFilter returns a collection of sites filtered by their last updated at in seconds
func (str *Store) ListFilter(lookbackPeriod int) []Site {
	str.RLock()
	defer str.RUnlock()

	filter := time.Now().Add(time.Duration(-lookbackPeriod) * time.Second)
	sites := make([]Site, 0)
	for _, site := range str.sites {
		if site.UpdatedAt.IsZero() || site.UpdatedAt.Before(filter) {
			sites = append(sites, *site)
		}
	}

	sort.Slice(sites, func(i, j int) bool {
		return sites[i].ID < sites[j].ID
	})

	return sites
}

// Add adds a single site to the store
func (str *Store) Add(st Site) error {
	// Validate URL
	u, err := url.ParseRequestURI(st.URL)
	if err != nil {
		return errors.New("Site URL is not valid")
	} else if u.Scheme == "" || u.Host == "" {
		return errors.New("Site URL must be an absolute URL")
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("Site URL must begin with http or https")
	}

	// Validate duplicate URL
	duplicate := false
	for _, site := range str.List() {
		if site.URL == st.URL {
			duplicate = true
			break
		}
	}

	if !duplicate {
		str.Lock()
		{
			str.idTracker = str.idTracker + 1
			st.ID = str.idTracker
			str.sites[str.idTracker] = &st
		}
		str.Unlock()
	}

	return nil
}

// UpdateHealth update the health status of a site
func (str *Store) UpdateHealth(siteID int, status int) error {
	str.Lock()
	defer str.Unlock()

	s, found := str.sites[siteID]

	if !found {
		return errors.New("Site does not exist")
	}

	s.Status = status
	s.UpdatedAt = time.Now()
	return nil
}

// Delete deletes a site from the store
func (str *Store) Delete(siteID int) error {
	str.Lock()
	defer str.Unlock()

	if _, ok := str.sites[siteID]; !ok {
		return errors.New("Site does not exist")
	}

	delete(str.sites, siteID)
	return nil
}
