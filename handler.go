package darwini

import (
	"net/http"

	"golang.org/x/net/context"
)

// Handler is a HTTP handler that also takes a Context.
type Handler interface {
	ServeHTTP(ctx context.Context, w http.ResponseWriter, req *http.Request)
}

// HandlerFunc is an adapter that allows the use of functions or
// methods as Handlers.
type HandlerFunc func(ctx context.Context, w http.ResponseWriter, req *http.Request)

var _ Handler = HandlerFunc(nil)

func (fn HandlerFunc) ServeHTTP(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	fn(ctx, w, req)
}
