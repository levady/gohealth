package httphandlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/levady/gohealth/internal/sitehealthchecker"
)

func TestHomepage(t *testing.T) {
	// Mocking
	implementedPath := homepageTplPath
	defer func() {
		homepageTplPath = implementedPath
	}()
	homepageTplPath = "../../../web/templates/homepage.html"

	// Request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Data preparation
	str := sitehealthchecker.NewStore()
	s := sitehealthchecker.Site{URL: "https://google.com"}
	str.Add(s)

	// Routing
	rr := httptest.NewRecorder()
	shh := SiteHealthHandler{
		SiteStore:           &str,
		HealtchCheckTimeout: 800 * time.Millisecond,
	}
	http.HandlerFunc(shh.Homepage).ServeHTTP(rr, req)
	resp := rr.Result()

	// Expectations
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}

	exp := string(`<li>https://google.com - </li><br />`)
	_, e := regexp.MatchString(exp, rr.Body.String())
	if e != nil {
		t.Errorf("Unexpected body %v", err)
	}
}

func TestHomepage_NotFound(t *testing.T) {
	// Request
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Routing
	str := sitehealthchecker.NewStore()
	shh := SiteHealthHandler{
		SiteStore:           &str,
		HealtchCheckTimeout: 800 * time.Millisecond,
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(shh.Homepage).ServeHTTP(rr, req)
	resp := rr.Result()

	// Expectations
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}
}

func TestSave(t *testing.T) {
	// Request
	form := url.Values{}
	form.Add("url", "http://zempag.com")
	req, err := http.NewRequest("POST", "/sites/save", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Routing
	rr := httptest.NewRecorder()
	str := sitehealthchecker.NewStore()
	shh := SiteHealthHandler{
		SiteStore:           &str,
		HealtchCheckTimeout: 800 * time.Millisecond,
	}
	http.HandlerFunc(shh.Save).ServeHTTP(rr, req)
	resp := rr.Result()

	// Expectations
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}

	exp := string(`<li>https://zempag.com - </li><br />`)
	_, e := regexp.MatchString(exp, rr.Body.String())
	if e != nil {
		t.Errorf("Unexpected body %v", rr.Body.String())
	}
}

func TestSave_Fail(t *testing.T) {
	// Mocking
	implementedPath := homepageTplPath
	defer func() {
		homepageTplPath = implementedPath
	}()
	homepageTplPath = "../../../web/templates/homepage.html"

	// Request
	form := url.Values{}
	form.Add("url", "saranghae")
	req, err := http.NewRequest("POST", "/sites/save", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Routing
	rr := httptest.NewRecorder()
	str := sitehealthchecker.NewStore()
	shh := SiteHealthHandler{
		SiteStore:           &str,
		HealtchCheckTimeout: 800 * time.Millisecond,
	}
	http.HandlerFunc(shh.Save).ServeHTTP(rr, req)
	resp := rr.Result()

	// Expectations
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}

	exp := string(`Site URL is not valid`)
	_, e := regexp.MatchString(exp, rr.Body.String())
	if e != nil {
		t.Errorf("Unexpected body %v", rr.Body.String())
	}
}

func TestHealthChecks(t *testing.T) {
	// Mocking
	implementedHealthChecks := runHealthChecks
	defer func() {
		runHealthChecks = implementedHealthChecks
	}()
	runHealthChecks = func(_ *sitehealthchecker.Store, _ time.Duration) {}

	// Request
	req, err := http.NewRequest("GET", "/ajax/sites/check", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Data preparations
	str := sitehealthchecker.NewStore()
	s := sitehealthchecker.Site{URL: "https://google.com"}
	str.Add(s)

	// Routing
	rr := httptest.NewRecorder()
	shh := SiteHealthHandler{
		SiteStore:           &str,
		HealtchCheckTimeout: 800 * time.Millisecond,
	}
	http.HandlerFunc(shh.HealthChecks).ServeHTTP(rr, req)
	resp := rr.Result()

	// Expectations
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}

	exp := string(`[{"id":1,"url":"https://google.com","healthy":null}]`)
	if body := rr.Body.String(); exp != body {
		t.Errorf("Unexpected body %v", body)
	}
}
