package darwini

import (
	"net/http"

	"golang.org/x/net/context"
)

// Map multiplexes requests to children based on a map lookup. For Map
// at /path, /path will be forbidden, /path/seg and /path/seg/anything
// are served by the map entry for seg, or not found if nil.
//
// As a special case, missing /path/ is forbidden instead of not
// found, to avoid a situation where /path/foo exists but its parent
// does not.
type Map map[string]Handler

var _ Handler = Map(nil)

func (c Map) ServeHTTP(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	seg, rest := segment(req.URL.Path)
	req.URL.Path = rest

	var child Handler
	if c != nil {
		child = c[seg]
	}
	if child == nil {
		if seg == "" {
			// Special case trailing slash; it's too easy to assume
			// that if "/foo/" 404s, "/foo/bar" won't exist.
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.NotFound(w, req)
		return
	}
	child.ServeHTTP(ctx, w, req)
}
