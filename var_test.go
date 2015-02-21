package darwini_test

import (
	"net/http"
	"testing"

	"github.com/tv42/darwini"
	"golang.org/x/net/context"
)

func TestVarSelf(t *testing.T) {
	m := darwini.Map{
		"foo": darwini.Var{},
	}
	resp := DoRequest(m, "GET", "/foo", nil)
	if resp.Code != http.StatusForbidden {
		t.Errorf("Var must be 403 Forbidden for itself: %v %v", resp.Code, resp.Body)
	}
}

func TestVarIndex(t *testing.T) {
	var seen bool
	h := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		seen = true
	}
	m := darwini.Var{
		Index: darwini.HandlerFunc(h),
	}
	resp := DoRequest(m, "GET", "/", nil)
	if !seen {
		t.Errorf("never saw GET: %v %v", resp.Code, resp.Body)
	}
}

func TestVarIndexIsNil(t *testing.T) {
	m := darwini.Var{}
	resp := DoRequest(m, "GET", "/", nil)
	if resp.Code != http.StatusForbidden {
		t.Errorf("Var index must be 403 Forbidden if not set: %v %v", resp.Code, resp.Body)
	}
}

func TestVarChild(t *testing.T) {
	var seen bool
	itemFn := func(seg string) darwini.Handler {
		if g, e := seg, "foo"; g != e {
			t.Errorf("darwini.Var gave wrong string: %q != %q", g, e)
		}
		return darwini.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
			if g, e := req.URL.Path, ""; g != e {
				t.Errorf("darwini.Var child url path is wrong: %q != %q", g, e)
			}
			seen = true
		})
	}
	m := darwini.Var{
		Child: itemFn,
	}
	resp := DoRequest(m, "GET", "/foo", nil)
	if !seen {
		t.Errorf("never saw GET: %v %v", resp.Code, resp.Body)
	}
}

func TestVarChildReturnsNil(t *testing.T) {
	m := darwini.Var{
		Child: func(seg string) darwini.Handler { return nil },
	}
	resp := DoRequest(m, "GET", "/foo", nil)
	if resp.Code != http.StatusNotFound {
		t.Errorf("Var item must be 404 Not Found if function returns nil: %v %v", resp.Code, resp.Body)
	}
}

func TestVarChildIsNil(t *testing.T) {
	m := darwini.Var{}
	resp := DoRequest(m, "GET", "/foo", nil)
	if resp.Code != http.StatusNotFound {
		t.Errorf("Var item must be 404 Not Found if not set: %v %v", resp.Code, resp.Body)
	}
}

func TestVarChildDeep(t *testing.T) {
	var seen bool
	childFn := func(seg string) darwini.Handler {
		if g, e := seg, "foo"; g != e {
			t.Errorf("darwini.Var gave wrong string: %q != %q", g, e)
		}
		return darwini.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
			if g, e := req.URL.Path, "/bar/baz"; g != e {
				t.Errorf("darwini.Var child url path is wrong: %q != %q", g, e)
			}
			seen = true
		})
	}
	m := darwini.Var{
		Child: childFn,
	}
	resp := DoRequest(m, "GET", "/foo/bar/baz", nil)
	if !seen {
		t.Errorf("never saw GET: %v %v", resp.Code, resp.Body)
	}
}
