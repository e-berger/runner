package probes

import "net/http"

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPClientMock struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (H HTTPClientMock) Do(r *http.Request) (*http.Response, error) {
	return H.DoFunc(r)
}
