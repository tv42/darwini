package darwini_test

import (
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/tv42/darwini"
)

func TestMethodBadMethod(t *testing.T) {
	m := darwini.Method{
		GET:    DoNotCall,
		POST:   DoNotCall,
		PUT:    DoNotCall,
		PATCH:  DoNotCall,
		DELETE: DoNotCall,
		Custom: map[string]http.HandlerFunc{"FROB": DoNotCall},
	}
	resp := DoRequest(m, "BAD", "", nil)
	if resp.Code != http.StatusMethodNotAllowed {
		t.Fatalf("Method must error on nil method: %v %v", resp.Code, resp.Body)
	}
	allow := resp.HeaderMap[http.CanonicalHeaderKey("Allow")]
	sort.Sort(sort.StringSlice(allow))
	got := strings.Join(allow, ", ")
	if g, e := got, "DELETE, FROB, GET, PATCH, POST, PUT"; g != e {
		t.Errorf("bad Allow header: %q != %q", g, e)
	}
}

func TestMethodGet(t *testing.T) {
	var seen bool
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Method{
		GET: h,
	}
	resp := DoRequest(m, "GET", "", nil)
	if !seen {
		t.Errorf("never saw GET: %v %v", resp.Code, resp.Body)
	}
}

func TestMethodPost(t *testing.T) {
	var seen bool
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Method{
		POST: h,
	}
	resp := DoRequest(m, "POST", "", nil)
	if !seen {
		t.Errorf("never saw POST: %v %v", resp.Code, resp.Body)
	}
}

func TestMethodPut(t *testing.T) {
	var seen bool
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Method{
		PUT: h,
	}
	resp := DoRequest(m, "PUT", "", nil)
	if !seen {
		t.Errorf("never saw PUT: %v %v", resp.Code, resp.Body)
	}
}

func TestMethodPatch(t *testing.T) {
	var seen bool
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Method{
		PATCH: h,
	}
	resp := DoRequest(m, "PATCH", "", nil)
	if !seen {
		t.Errorf("never saw PATCH: %v %v", resp.Code, resp.Body)
	}
}

func TestMethodDelete(t *testing.T) {
	var seen bool
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Method{
		DELETE: h,
	}
	resp := DoRequest(m, "DELETE", "", nil)
	if !seen {
		t.Errorf("never saw DELETE: %v %v", resp.Code, resp.Body)
	}
}

func TestMethodCustom(t *testing.T) {
	var seen bool
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Method{
		Custom: map[string]http.HandlerFunc{"FROB": h},
	}
	resp := DoRequest(m, "FROB", "", nil)
	if !seen {
		t.Errorf("never saw FROB: %v %v", resp.Code, resp.Body)
	}
}

func TestMethodSlash(t *testing.T) {
	var seen bool
	var seenPath string
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
		seenPath = req.URL.Path
	}
	m := darwini.Method{
		GET: h,
	}
	resp := DoRequest(m, "GET", "/", nil)
	if !seen {
		t.Errorf("never saw GET with a slash: %v %v", resp.Code, resp.Body)
	}
	if g, e := seenPath, "/"; g != e {
		t.Errorf("Method should not change the URL path: %q != %q", g, e)
	}
}

func TestMethodChild(t *testing.T) {
	var seen bool
	var seenPath string
	h := func(w http.ResponseWriter, req *http.Request) {
		seen = true
		seenPath = req.URL.Path
	}
	m := darwini.Method{
		GET: h,
	}
	resp := DoRequest(m, "GET", "/child", nil)
	if !seen {
		t.Errorf("never saw GET with a slash: %v %v", resp.Code, resp.Body)
	}
	if g, e := seenPath, "/child"; g != e {
		t.Errorf("Method should not change the URL path: %q != %q", g, e)
	}
}
