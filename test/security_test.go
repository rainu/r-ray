package test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUnauthorized(t *testing.T) {
	req := validRequest(http.MethodGet, "/", nil, nil)

	//no auth given...
	req.Header.Del("Authorization")

	called := false
	mock = func(w http.ResponseWriter, r *http.Request) {
		called = true
	}

	resp := do(req)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "The application should cancel the request cause of missing authentication!")
	assert.False(t, called, "The target server should not be called!")
}

func TestBadCredentials(t *testing.T) {
	req := validRequest(http.MethodGet, "/", nil, nil)
	req.SetBasicAuth(username, password+"-invalid")

	called := false
	mock = func(w http.ResponseWriter, r *http.Request) {
		called = true
	}

	resp := do(req)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "The application should cancel the request cause of invalid credentials!")
	assert.False(t, called, "The target server should not be called!")
}
