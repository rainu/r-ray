package processor

import "net/http"

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type UserStore interface {
	IsValid(username, password string) bool
}
