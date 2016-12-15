package darwini_test

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func MustNewRequest(method, urlStr string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		panic(err)
	}
	return req
}

func DoRequest(h http.Handler, method, urlStr string, body io.Reader) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	req := MustNewRequest(method, urlStr, body)
	h.ServeHTTP(resp, req)
	return resp
}

func DoNotCall(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "not expecting to see this", http.StatusInternalServerError)
}
