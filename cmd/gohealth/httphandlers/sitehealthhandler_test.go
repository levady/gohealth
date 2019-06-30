package httphandlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/levady/gohealth/internal/platform/sitestore"
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
	str := sitestore.NewStore()
	s := sitestore.Site{URL: "https://google.com"}
	str.Add(s)

	// Routing
	rr := httptest.NewRecorder()
	shh := SiteHealthHandler{SiteStore: &str}
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
	var testCases = []struct {
		name          string
		route         string
		method        string
		expStatusCode int
	}{
		{
			name:          "POST request to `/`",
			route:         "/",
			method:        "POST",
			expStatusCode: 404,
		},
		{
			name:          "DELETE request to `/oneinamillion/twice/sumida`",
			route:         "/oneinamillion/twice/sumida",
			method:        "DELETE",
			expStatusCode: 404,
		},
		{
			name:          "GET request to `/blackpink/in/your/area`",
			route:         "/blackpink/in/your/area",
			method:        "GET",
			expStatusCode: 404,
		},
		{
			name:          "PUT request to `/spoonman`",
			route:         "/spoonman",
			method:        "PUT",
			expStatusCode: 404,
		},
	}

	for _, tc := range testCases {
		// Request
		req, err := http.NewRequest(tc.method, tc.route, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Routing
		str := sitestore.NewStore()
		shh := SiteHealthHandler{SiteStore: &str}
		rr := httptest.NewRecorder()
		http.HandlerFunc(shh.Homepage).ServeHTTP(rr, req)
		resp := rr.Result()

		// Expectations
		if tc.expStatusCode != http.StatusNotFound {
			t.Errorf("Unexpected status code %d", resp.StatusCode)
		}
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
	str := sitestore.NewStore()
	shh := SiteHealthHandler{SiteStore: &str}
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
	str := sitestore.NewStore()
	shh := SiteHealthHandler{SiteStore: &str}
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

func TestDelete(t *testing.T) {
	var testCases = []struct {
		name            string
		siteID          string
		expStatusCode   int
		expResponseBody string
		hasSite         bool
	}{
		{
			name:            "Deleting an existing site",
			siteID:          "1",
			expStatusCode:   http.StatusOK,
			expResponseBody: "{}",
			hasSite:         true,
		},
		{
			name:            "Deleting a non existing site",
			siteID:          "100",
			expStatusCode:   http.StatusNotFound,
			expResponseBody: "404 not found. \n",
			hasSite:         false,
		},
	}

	for _, tc := range testCases {
		str := sitestore.NewStore()
		// Data preparations
		if tc.hasSite {
			s := sitestore.Site{URL: "https://google.com"}
			str.Add(s)
		}

		// Request
		req, err := http.NewRequest("DELETE", "/ajax/sites/delete/"+tc.siteID, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Routing
		rr := httptest.NewRecorder()
		shh := SiteHealthHandler{SiteStore: &str}
		http.HandlerFunc(shh.Delete).ServeHTTP(rr, req)
		resp := rr.Result()

		// Expectations
		if tc.expStatusCode != resp.StatusCode {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}

		_, e := regexp.MatchString(tc.expResponseBody, rr.Body.String())
		if e != nil {
			t.Errorf("Unexpected body %v", rr.Body.String())
		}
	}
}

func TestHealthChecks(t *testing.T) {
	// Request
	req, err := http.NewRequest("GET", "/ajax/sites/check", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Data preparations
	str := sitestore.NewStore()
	s := sitestore.Site{URL: "https://google.com"}
	str.Add(s)

	// Routing
	rr := httptest.NewRecorder()
	shh := SiteHealthHandler{SiteStore: &str}
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
