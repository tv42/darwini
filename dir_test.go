package darwini_test

import (
	"net/http"
	"testing"

	"github.com/tv42/darwini"
	"golang.org/x/net/context"
)

func TestDirParent(t *testing.T) {
	var seen bool
	h := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Map{
		"foo": darwini.Dir{
			Parent: darwini.HandlerFunc(h),
		},
	}
	resp := DoRequest(m, "GET", "/foo", nil)
	if !seen {
		t.Errorf("never saw GET: %v %v", resp.Code, resp.Body)
	}
}

func TestDirChild(t *testing.T) {
	var seen bool
	h := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Map{
		"foo": darwini.Dir{
			Child: darwini.Map{
				"bar": darwini.HandlerFunc(h),
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
	h := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Map{
		"foo": darwini.Dir{
			Child: darwini.Map{
				"": darwini.HandlerFunc(h),
			},
		},
	}
	resp := DoRequest(m, "GET", "/foo/", nil)
	if !seen {
		t.Errorf("never saw GET: %v %v", resp.Code, resp.Body)
	}
}
