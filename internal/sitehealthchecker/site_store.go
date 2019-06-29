package sitehealthchecker

import (
	"errors"
	"net/url"
	"sort"
	"sync"
)

// Site represents Site data
type Site struct {
	ID      int64       `json:"id"`
	URL     string      `json:"url"`
	Healthy interface{} `json:"healthy"`
}

// Store represent data store for sites
type Store struct {
	sites     map[int64]*Site
	idTracker int64
	sync.RWMutex
}

// NewStore construct a new Store
func NewStore() Store {
	return Store{
		sites:     make(map[int64]*Site),
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
func (str *Store) UpdateHealth(siteID int64, status bool) {
	str.Lock()
	defer str.Unlock()

	if s, ok := str.sites[siteID]; ok {
		s.Healthy = status
	}
}

// HealthyIsNotNil returns false if Healthy attribute is not nil
func (s *Site) HealthyIsNotNil() bool {
	return s.Healthy != nil
}
