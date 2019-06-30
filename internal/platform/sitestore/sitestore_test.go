package sitestore

import (
	"fmt"
	"testing"
	"time"
)

var (
	site1 = Site{URL: "https://google.com"}
	site2 = Site{URL: "https://golang.org/doc/articles/wiki/#tmp_7"}
	site3 = Site{URL: "http://stat.us/200?sleep=40000"}
	site4 = Site{URL: "https://smartystreets.com/blog/2015/02/go-testing-part-1-vanillla"}
	site5 = Site{URL: "http://stat.us/200?sleep=10000"}
)

func TestNewStore(t *testing.T) {
	str := NewStore()

	if sCount := len(str.sites); sCount != 0 {
		t.Errorf("Expected Sites length of %d, but it was %d instead.", 0, sCount)
	}

	if str.idTracker != 0 {
		t.Errorf("Expected idTracker to start with 0, but it was %v instead.", str.idTracker)
	}
}

func TestList(t *testing.T) {
	str := NewStore()
	str.Add(site1)
	str.Add(site2)
	str.Add(site3)
	str.Add(site4)
	str.Add(site5)

	sites := str.List()

	if s := sites[0]; s.URL != "https://google.com" {
		t.Errorf("Expected the first site in the array to be google.com, but it was %v.", s.URL)
	}

	if s := sites[len(sites)-1]; s.URL != "http://stat.us/200?sleep=10000" {
		t.Errorf("Expected the first site in the array to be stat.us, but it was %v.", s.URL)
	}
}

func TestListFilter(t *testing.T) {
	site1.UpdatedAt = time.Now().Add(time.Duration(-12) * time.Second)
	site2.UpdatedAt = time.Now().Add(time.Duration(-15) * time.Second)
	site3.UpdatedAt = time.Time{}
	site4.UpdatedAt = time.Now().Add(time.Duration(-100) * time.Second)
	site5.UpdatedAt = time.Now()

	str := NewStore()

	str.Add(site1)
	str.Add(site2)
	str.Add(site3)
	str.Add(site4)
	str.Add(site5)

	sites := str.ListFilter(15)

	for _, s := range sites {
		fmt.Printf("======> copied site %v: %v \n", s.URL, s.updatedAt)
	}

	if len(sites) != 3 {
		t.Errorf("Expected result length to 3 but it was %v", len(sites))
	}

	if s := sites[0]; s.URL != "https://golang.org/doc/articles/wiki/#tmp_7" {
		t.Errorf("Expected the first site in the array to be golang.org, but it was %v.", s.URL)
	}

	if s := sites[len(sites)-1]; s.URL != "https://smartystreets.com/blog/2015/02/go-testing-part-1-vanillla" {
		t.Errorf("Expected the last site in the array to be smartystreets.com, but it was %v.", s.URL)
	}
}

func TestAdd(t *testing.T) {
	var testCases = []struct {
		name   string
		input  []Site
		exp    int
		hasErr bool
	}{
		{
			name: "Adding multiple Sites",
			input: []Site{
				Site{
					URL:     "https://golang.org/doc/articles/wiki/",
					Healthy: nil,
				},
				Site{
					URL:     "https://google.com/",
					Healthy: nil,
				},
			},
			exp:    2,
			hasErr: false,
		},
		{
			name: "Adding an empty URL",
			input: []Site{
				Site{
					URL:     "",
					Healthy: nil,
				},
			},
			exp:    0,
			hasErr: true,
		},
		{
			name: "Adding a non absolute URL",
			input: []Site{
				Site{
					URL:     "/sites/save",
					Healthy: nil,
				},
			},
			exp:    0,
			hasErr: true,
		},
		{
			name: "Adding URL that does not start with http or https",
			input: []Site{
				Site{
					URL:     "ftp://websiteaddress.com",
					Healthy: nil,
				},
			},
			exp:    0,
			hasErr: true,
		},
		{
			name: "Adding duplicate Sites",
			input: []Site{
				Site{
					URL:     "https://golang.org/doc/articles/wiki/",
					Healthy: nil,
				},
				Site{
					URL:     "https://google.com/",
					Healthy: nil,
				},
				Site{
					URL:     "https://golang.org/doc/articles/wiki/",
					Healthy: nil,
				},
			},
			exp:    2,
			hasErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			str := NewStore()

			var err error
			for _, s := range tc.input {
				err = str.Add(s)
			}

			if !tc.hasErr && err != nil {
				t.Errorf("Error is not expected. Got err: %v", err)
			}

			sCount := len(str.sites)
			if sCount != tc.exp {
				t.Errorf("Expected Sites length of %d, but it was %d instead.", tc.exp, sCount)
			}
		})
	}
}

func TestAdd_AutoIncrementID(t *testing.T) {
	str := NewStore()
	str.Add(site1)
	str.Add(site2)
	str.Add(site3)
	str.Add(site4)
	str.Add(site5)

	for i := 0; i < len(str.sites); i++ {
		_, ok := str.sites[int64(i+1)]
		if !ok {
			t.Errorf("Expected to have key %d but got nil", i+1)
		}
	}
}

func TestUpdateHealth(t *testing.T) {
	str := NewStore()
	str.Add(site1)

	var testCases = []struct {
		name   string
		siteID int64
		input  bool
		exp    interface{}
		hasErr bool
	}{
		{
			name:   "Update to health to true",
			siteID: 1,
			input:  true,
			exp:    true,
			hasErr: false,
		},
		{
			name:   "Update to health to false",
			siteID: 1,
			input:  false,
			exp:    false,
			hasErr: false,
		},
		{
			name:   "Updating a site that does not exist",
			siteID: 100,
			input:  true,
			exp:    nil,
			hasErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := str.UpdateHealth(tc.siteID, tc.input)

			if tc.hasErr && err == nil {
				t.Errorf("Expected to return an error but got nil")
			}

			s := str.sites[tc.siteID]
			if !tc.hasErr && s.Healthy != tc.exp {
				t.Errorf("Expected site to be updated to %v but got %v.", tc.exp, s.Healthy)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	str := NewStore()
	str.Add(site1)

	var testCases = []struct {
		name   string
		siteID int64
		exp    interface{}
	}{
		{
			name:   "Deleting an existing site",
			siteID: 1,
			exp:    nil,
		},
		{
			name:   "Deleting a site that does not exist",
			siteID: 100,
			exp:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			str.Delete(tc.siteID)
			if s, ok := str.sites[tc.siteID]; ok {
				t.Errorf("Expected site to be deleted but got %v.", s)
			}
		})
	}
}

func TestHealthyIsNotNil(t *testing.T) {
	var testCases = []struct {
		name  string
		input interface{}
		exp   bool
	}{
		{
			name:  "Healthy is nil",
			input: nil,
			exp:   false,
		},
		{
			name:  "Healthy is false",
			input: false,
			exp:   true,
		},
		{
			name:  "Healthy is true",
			input: true,
			exp:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := Site{Healthy: tc.input}

			if s.HealthyIsNotNil() != tc.exp {
				t.Errorf("Expected site healthy to return %v but got %v.", tc.exp, s.HealthyIsNotNil())
			}
		})
	}
}
