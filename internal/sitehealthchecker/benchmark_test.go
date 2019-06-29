package sitehealthchecker

// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// go test -run none -bench . -benchtime 3s

// Basic benchmark test.

import (
	"testing"
	"time"

	"github.com/levady/gohealth/internal/platform/sitestore"
)

func BenchmarkParallelHealthChecks(b *testing.B) {
	site1 := sitestore.Site{URL: "https://google.com"}
	site2 := sitestore.Site{URL: "https://golang.org/doc/articles/wiki/#tmp_7"}
	site3 := sitestore.Site{URL: "http://stat.us/200?sleep=40000"}
	site4 := sitestore.Site{URL: "https://smartystreets.com/blog/2015/02/go-testing-part-1-vanillla"}
	site5 := sitestore.Site{URL: "http://stat.us/200?sleep=10000"}

	store := sitestore.NewStore()
	store.Add(site1)
	store.Add(site2)
	store.Add(site3)
	store.Add(site4)
	store.Add(site5)

	for i := 0; i < b.N; i++ {
		ParallelHealthChecks(&store, 800*time.Millisecond)
	}
}

func BenchmarkSerialHealthChecks(b *testing.B) {
	site1 := sitestore.Site{URL: "https://google.com"}
	site2 := sitestore.Site{URL: "https://golang.org/doc/articles/wiki/#tmp_7"}
	site3 := sitestore.Site{URL: "http://stat.us/200?sleep=40000"}
	site4 := sitestore.Site{URL: "https://smartystreets.com/blog/2015/02/go-testing-part-1-vanillla"}
	site5 := sitestore.Site{URL: "http://stat.us/200?sleep=10000"}

	store := sitestore.NewStore()
	store.Add(site1)
	store.Add(site2)
	store.Add(site3)
	store.Add(site4)
	store.Add(site5)

	for i := 0; i < b.N; i++ {
		SerialHealthChecks(&store, 800*time.Millisecond)
	}
}
