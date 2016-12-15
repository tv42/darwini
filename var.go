package darwini

import (
	"net/http"
)

// Var multiplexes dynamically based on the next path segment. For Var
// at /path, /path will be forbidden, /path/ is served by Index or
// forbidden if nil, /path/seg and /path/seg/anything are served by
// the handler Child returns for seg, or not found if Child is nil or
// returns nil.
type Var struct {
	Index http.Handler
	Child func(seg string) http.Handler
}

var _ http.Handler = Var{}

func (v Var) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
		v.Index.ServeHTTP(w, req)
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
	h.ServeHTTP(w, req)
}
