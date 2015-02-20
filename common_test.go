package darwini_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/tv42/darwini"
	"golang.org/x/net/context"
)

func MustNewRequest(method, urlStr string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		panic(err)
	}
	return req
}

func DoRequest(h darwini.Handler, method, urlStr string, body io.Reader) *httptest.ResponseRecorder {
	ctx := context.Background()
	resp := httptest.NewRecorder()
	req := MustNewRequest(method, urlStr, body)
	h.ServeHTTP(ctx, resp, req)
	return resp
}

func DoNotCall(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	http.Error(w, "not expecting to see this", http.StatusInternalServerError)
}
