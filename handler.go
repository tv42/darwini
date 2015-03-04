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

// WithContext adapts a darwini.Handler to work as a net/http.Handler.
// The context will be canceled when the HTTP client disconnects.
func WithContext(h Handler) http.Handler {
	return withContext{h}
}

type withContext struct {
	Handler
}

var _ http.Handler = withContext{}

func (h withContext) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if n, ok := w.(http.CloseNotifier); ok {
		gone := n.CloseNotify()
		go func() {
			select {
			case <-ctx.Done():
				return
			case <-gone:
				cancel()
			}
		}()
	}

	h.Handler.ServeHTTP(ctx, w, req)
}

// NoContext adapts a net/http Handler to work as a darwini.Handler.
// Context deadlines or cancellation are not respected.
func NoContext(h http.Handler) Handler {
	return noContext{h}
}

type noContext struct {
	http.Handler
}

var _ Handler = noContext{}

func (h noContext) ServeHTTP(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	h.Handler.ServeHTTP(w, req)
}
