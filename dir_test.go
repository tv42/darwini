package darwini_test

import (
	"net/http"
	"testing"

	"github.com/tv42/darwini"
)

func TestDirParent(t *testing.T) {
	var seen bool
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Map{
		"foo": darwini.Dir{
			Parent: http.HandlerFunc(h),
		},
	}
	resp := DoRequest(m, "GET", "/foo", nil)
	if !seen {
		t.Errorf("never saw GET: %v %v", resp.Code, resp.Body)
	}
}

func TestDirChild(t *testing.T) {
	var seen bool
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Map{
		"foo": darwini.Dir{
			Child: darwini.Map{
				"bar": http.HandlerFunc(h),
			},
		},
	}
	resp := DoRequest(m, "GET", "/foo/bar", nil)
	if !seen {
		t.Errorf("never saw GET: %v %v", resp.Code, resp.Body)
	}
}

func TestDirChildSlashOnly(t *testing.T) {
	var seen bool
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Map{
		"foo": darwini.Dir{
			Child: darwini.Map{
				"": http.HandlerFunc(h),
			},
		},
	}
	resp := DoRequest(m, "GET", "/foo/", nil)
	if !seen {
		t.Errorf("never saw GET: %v %v", resp.Code, resp.Body)
	}
}
