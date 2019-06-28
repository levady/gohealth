package sitehealthchecker

import (
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
