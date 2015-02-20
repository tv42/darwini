package darwini_test

import (
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/tv42/darwini"
	"golang.org/x/net/context"
)

func TestMethodBadSlash(t *testing.T) {
	m := darwini.Method{
		GET: DoNotCall,
	}
	resp := DoRequest(m, "GET", "/", nil)
	if resp.Code != http.StatusNotFound {
		t.Errorf("Method must not handle children: %v %v", resp.Code, resp.Body)
	}
}

func TestMethodBadChild(t *testing.T) {
	m := darwini.Method{
		GET: DoNotCall,
	}
	resp := DoRequest(m, "GET", "/bad", nil)
	if resp.Code != http.StatusNotFound {
		t.Errorf("Method must not handle children: %v %v", resp.Code, resp.Body)
	}
}

func TestMethodBadMethod(t *testing.T) {
	m := darwini.Method{
		GET:    DoNotCall,
		POST:   DoNotCall,
		PUT:    DoNotCall,
		PATCH:  DoNotCall,
		DELETE: DoNotCall,
		Custom: map[string]darwini.HandlerFunc{"FROB": DoNotCall},
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
	h := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
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
	h := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
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
	h := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
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
	h := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
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
	h := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
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
	h := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Method{
		Custom: map[string]darwini.HandlerFunc{"FROB": h},
	}
	resp := DoRequest(m, "FROB", "", nil)
	if !seen {
		t.Errorf("never saw FROB: %v %v", resp.Code, resp.Body)
	}
}
