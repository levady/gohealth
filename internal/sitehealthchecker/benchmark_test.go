package sitehealthchecker

// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// go test -run none -bench . -benchtime 3s

// Basic benchmark test.

import (
	"testing"
	"time"
)

func BenchmarkParallelHealthChecks(b *testing.B) {
	site1 := Site{URL: "https://google.com"}
	site2 := Site{URL: "https://golang.org/doc/articles/wiki/#tmp_7"}
	site3 := Site{URL: "http://stat.us/200?sleep=40000"}
	site4 := Site{URL: "https://smartystreets.com/blog/2015/02/go-testing-part-1-vanillla"}
	site5 := Site{URL: "http://stat.us/200?sleep=10000"}

	shc := New(800 * time.Millisecond)
	shc.AddSite(site1)
	shc.AddSite(site2)
	shc.AddSite(site3)
	shc.AddSite(site4)
	shc.AddSite(site5)

	for i := 0; i < b.N; i++ {
		shc.ParallelHealthChecks()
	}
}

func BenchmarkSerialHealthChecks(b *testing.B) {
	site1 := Site{URL: "https://google.com"}
	site2 := Site{URL: "https://golang.org/doc/articles/wiki/#tmp_7"}
	site3 := Site{URL: "http://stat.us/200?sleep=40000"}
	site4 := Site{URL: "https://smartystreets.com/blog/2015/02/go-testing-part-1-vanillla"}
	site5 := Site{URL: "http://stat.us/200?sleep=10000"}

	shc := New(800 * time.Millisecond)
	shc.AddSite(site1)
	shc.AddSite(site2)
	shc.AddSite(site3)
	shc.AddSite(site4)
	shc.AddSite(site5)

	for i := 0; i < b.N; i++ {
		shc.SerialHealthChecks()
	}
}
