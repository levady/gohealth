package sitehealthchecker

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/levady/gohealth/internal/platform/sitestore"
)

func TestParallelHealthChecks(t *testing.T) {
	// Mocking
	implementedSiteChecker := siteChecker
	defer func() {
		siteChecker = implementedSiteChecker
	}()

	site1 := sitestore.Site{URL: "https://zempag.com"}
	site2 := sitestore.Site{URL: "https://www.google.com"}
	site3 := sitestore.Site{URL: "https://koprol.com"}

	store := sitestore.NewStore()
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

	ParallelHealthChecks(&store, 800*time.Millisecond, 0)

	sites := store.List()

	if s := sites[0]; s.Healthy == true {
		t.Errorf("Expected Site1 %v to timeout.", s.URL)
	}

	if s := sites[1]; s.Healthy == true {
		t.Errorf("Expected Site2 %v to return 500.", s.URL)
	}

	if s := sites[2]; s.Healthy == false {
		t.Errorf("Expected Site3 %v to be healthy.", s.URL)
	}
}

func TestParallelHealthChecksWithLookbackPeriod(t *testing.T) {
	// Mocking
	implementedSiteChecker := siteChecker
	defer func() {
		siteChecker = implementedSiteChecker
	}()

	site1 := sitestore.Site{URL: "https://zempag.com", UpdatedAt: time.Now().Add(time.Duration(-12) * time.Second)}
	site2 := sitestore.Site{URL: "https://www.google.com", UpdatedAt: time.Now().Add(time.Duration(-27) * time.Second)}
	site3 := sitestore.Site{URL: "https://koprol.com"}

	store := sitestore.NewStore()
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

	ParallelHealthChecks(&store, 800*time.Millisecond, 15)

	sites := store.List()

	s1 := sites[0]
	if s1.HealthyIsNotNil() {
		t.Errorf("Expected Site1 'Health' not to be updated but got %v", s1.Healthy)
	}

	if s1.UpdatedAt != site1.UpdatedAt {
		t.Errorf("Expected Site1 'UpdatedAt' not to be updated but got %v", s1.UpdatedAt)
	}

	if s := sites[1]; s.Healthy == true {
		t.Errorf("Expected Site2 %v to return 500.", s.URL)
	}

	if s := sites[2]; s.Healthy == false {
		t.Errorf("Expected Site3 %v to be healthy.", s.URL)
	}
}
