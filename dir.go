package darwini

import (
	"net/http"
)

// Dir multiplexes requests between a parent resource and children.
// For Dir at /path, Parent will serve requests to /path, while Child
// will serve requests to /path/anything.
type Dir struct {
	Parent http.Handler
	Child  http.Handler
}

var _ http.Handler = Dir{}

func (d Dir) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "" {
		if d.Parent == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		d.Parent.ServeHTTP(w, req)
		return
	}
	if d.Child == nil {
		http.NotFound(w, req)
		return
	}
	d.Child.ServeHTTP(w, req)
}
