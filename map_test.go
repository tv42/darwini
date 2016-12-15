package darwini_test

import (
	"net/http"
	"testing"

	"github.com/tv42/darwini"
)

func TestMapSimple(t *testing.T) {
	var seen bool
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Map{
		"foo": http.HandlerFunc(h),
	}
	resp := DoRequest(m, "GET", "/foo", nil)
	if !seen {
		t.Errorf("never saw GET: %v %v", resp.Code, resp.Body)
	}
}

func TestMapNotFound(t *testing.T) {
	m := darwini.Map{}
	resp := DoRequest(m, "GET", "/bar", nil)
	if resp.Code != http.StatusNotFound {
		t.Errorf("Map must not find non-existent children: %v %v", resp.Code, resp.Body)
	}
}

func TestMapSelf(t *testing.T) {
	m := darwini.Map{
		"foo": darwini.Map{},
	}
	resp := DoRequest(m, "GET", "/foo", nil)
	if resp.Code != http.StatusForbidden {
		t.Errorf("Map must serve 403 Forbidden for itself: %v %v", resp.Code, resp.Body)
	}
}

func TestMapSlash(t *testing.T) {
	var seen bool
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Map{
		"": http.HandlerFunc(h),
	}
	resp := DoRequest(m, "GET", "/", nil)
	if !seen {
		t.Errorf("never saw GET: %v %v", resp.Code, resp.Body)
	}
}

func TestMapSlashImplicit(t *testing.T) {
	m := darwini.Map{}
	resp := DoRequest(m, "GET", "/", nil)
	if resp.Code != http.StatusForbidden {
		t.Errorf("Map must serve 403 Forbidden without an index: %v %v", resp.Code, resp.Body)
	}
}
