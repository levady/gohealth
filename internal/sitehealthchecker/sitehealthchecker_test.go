package sitehealthchecker

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestParallelHealthChecks(t *testing.T) {
	// Mocking
	implementedSiteChecker := siteChecker
	defer func() {
		siteChecker = implementedSiteChecker
	}()

	site1 := Site{URL: "https://zempag.com"}
	site2 := Site{URL: "https://www.google.com"}
	site3 := Site{URL: "https://koprol.com"}

	store := NewStore()
	store.Add(site1)
	store.Add(site2)
	store.Add(site3)

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

	ParallelHealthChecks(&store, 800*time.Millisecond)

	if s := store.sites[1]; s.Healthy == true {
		t.Errorf("Expected Site1 %v to timeout.", s.URL)
	}

	if s := store.sites[2]; s.Healthy == true {
		t.Errorf("Expected Site2 %v to return 500.", s.URL)
	}

	if s := store.sites[3]; s.Healthy == false {
		t.Errorf("Expected Site3 %v to be healthy.", s.URL)
	}
}
