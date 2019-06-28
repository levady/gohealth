package sitehealthchecker

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	shc := New(15 * time.Second)

	if sCount := len(shc.Sites); sCount != 0 {
		t.Errorf("Expected Sites length of %d, but it was %d instead.", 0, sCount)
	}

	if timeout := shc.Timeout; timeout != 15*time.Second {
		t.Errorf("Expected Sites length of %d, but it was %v instead.", timeout, shc.Timeout)
	}
}

func TestAddSite(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shc := New(15 * time.Second)

			var err error
			for _, s := range tc.input {
				err = shc.AddSite(s)
			}

			if !tc.hasErr && err != nil {
				t.Errorf("Error is not expected. Got err: %v", err)
			}

			if sCount := len(shc.Sites); sCount != tc.exp {
				t.Errorf("Expected Sites length of %d, but it was %d instead.", tc.exp, sCount)
			}
		})
	}
}

func TestRunHealthChecks(t *testing.T) {
	// Mocking
	implementedSiteChecker := siteChecker
	defer func() {
		siteChecker = implementedSiteChecker
	}()

	site1 := Site{URL: "https://zempag.com"}
	site2 := Site{URL: "https://www.google.com"}
	site3 := Site{URL: "https://koprol.com"}

	shc := New(800 * time.Millisecond)
	shc.AddSite(site1)
	shc.AddSite(site2)
	shc.AddSite(site3)

	siteChecker = func(url string, _ time.Duration) (*http.Response, error) {
		switch url {
		case "https://zempag.com":
			return &http.Response{}, errors.New("Timeout")
		case "https://www.google.com":
			return &http.Response{StatusCode: 500}, nil
		default:
			return &http.Response{StatusCode: 200}, nil
		}
	}

	shc.RunHealthChecks()

	if shc.Sites[0].Healthy == true {
		t.Errorf("Expected Site1 %v to timeout.", shc.Sites[0].URL)
	}

	if shc.Sites[1].Healthy == true {
		t.Errorf("Expected Site2 %v to return 500.", shc.Sites[1].URL)
	}

	if shc.Sites[2].Healthy == false {
		t.Errorf("Expected Site3 %v to be healthy.", shc.Sites[2].URL)
	}
}
