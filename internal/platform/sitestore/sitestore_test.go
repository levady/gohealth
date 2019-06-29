package sitestore

import (
	"fmt"
	"testing"
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
		fmt.Printf("Site1 is %v \n", s)
		t.Errorf("Expected the first site in the array to be google.com, but it was %v.", s.URL)
	}

	if s := sites[len(sites)-1]; s.URL != "http://stat.us/200?sleep=10000" {
		fmt.Printf("Site1 is %v \n", s)
		t.Errorf("Expected the first site in the array to be stat.us, but it was %v.", s.URL)
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
		name  string
		input bool
		exp   bool
	}{
		{
			name:  "Update to health to true",
			input: true,
			exp:   true,
		},
		{
			name:  "Update to health to false",
			input: false,
			exp:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := str.sites[1]
			str.UpdateHealth(s.ID, tc.input)

			if s.Healthy != tc.exp {
				t.Errorf("Expected site to be updated to %v but got %v.", tc.exp, s.Healthy)
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
