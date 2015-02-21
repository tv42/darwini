package darwini

import (
	"net/http"

	"golang.org/x/net/context"
)

// Var multiplexes dynamically based on the next path segment. For Var
// at /path, /path will be forbidden, /path/ is served by Index or
// forbidden if nil, /path/seg and /path/seg/anything are served by
// the handler Child returns for seg, or not found if Child is nil or
// returns nil.
type Var struct {
	Index Handler
	Child func(seg string) Handler
}

var _ Handler = Var{}

func (v Var) ServeHTTP(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "" {
		// Var does not serve "itself", ever. Use a wrapper Dir if you
		// want that.
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	seg, rest := segment(req.URL.Path)
	req.URL.Path = rest

	if seg == "" {
		if v.Index == nil {
			// No index resource.
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		v.Index.ServeHTTP(ctx, w, req)
		return
	}

	if v.Child == nil {
		// No child resource.
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	h := v.Child(seg)
	if h == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	h.ServeHTTP(ctx, w, req)
}
