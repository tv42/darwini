package darwini

import (
	"net/http"
)

// Method multiplexes requests based on the HTTP method. The handler
// for a method is set either with the predefined fields, or for
// custom methods, with the Custom map.
//
// The values here are HandlerFuncs and not Handlers, as it is common
// to make them be methods on the same value.
type Method struct {
	GET    http.HandlerFunc
	POST   http.HandlerFunc
	PUT    http.HandlerFunc
	PATCH  http.HandlerFunc
	DELETE http.HandlerFunc
	Custom map[string]http.HandlerFunc
}

var _ http.Handler = Method{}

func (m Method) get(method string) http.HandlerFunc {
	switch method {
	case "GET":
		return m.GET
	case "POST":
		return m.POST
	case "PUT":
		return m.PUT
	case "PATCH":
		return m.PATCH
	case "DELETE":
		return m.DELETE
	default:
		return m.Custom[method]
	}
}

func (m Method) err(w http.ResponseWriter, req *http.Request) {
	if m.GET != nil {
		w.Header().Add("Allow", "GET")
	}
	if m.POST != nil {
		w.Header().Add("Allow", "POST")
	}
	if m.PUT != nil {
		w.Header().Add("Allow", "PUT")
	}
	if m.PATCH != nil {
		w.Header().Add("Allow", "PATCH")
	}
	if m.DELETE != nil {
		w.Header().Add("Allow", "DELETE")
	}
	for k, v := range m.Custom {
		if v != nil {
			w.Header().Add("Allow", k)
		}
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (m Method) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h := m.get(req.Method)
	if h == nil {
		m.err(w, req)
		return
	}
	h.ServeHTTP(w, req)
}
