package darwini

import (
	"net/http"

	"golang.org/x/net/context"
)

// Dir multiplexes requests between a parent resource and children.
// For Dir at /path, Parent will serve requests to /path, while Child
// will serve requests to /path/anything.
type Dir struct {
	Parent Handler
	Child  Handler
}

var _ Handler = Dir{}

func (d Dir) ServeHTTP(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "" {
		if d.Parent == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		d.Parent.ServeHTTP(ctx, w, req)
		return
	}
	if d.Child == nil {
		http.NotFound(w, req)
		return
	}
	d.Child.ServeHTTP(ctx, w, req)
}
