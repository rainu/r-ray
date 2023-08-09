package test

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSimpleWithoutAnyTransfer(t *testing.T) {
	mock = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusAccepted)
	}
	resp := do(validRequest(http.MethodGet, "/", nil, nil))

	assert.Equal(t, http.StatusOK, resp.StatusCode, "The proxies return code should always be 200 in happy case.")
	assert.Equal(t, http.Header(map[string][]string{
		"Content-Length":                {"0"},                     //origin: proxy (because the target body is less than 512b)
		headerPrefix + "Content-Length": {"0"},                     //origin: target (because the target body is less than 512b)
		headerPrefix + "Status-Line":    {"HTTP/1.1 202 Accepted"}, //origin: proxy (concatenation of targets response code)
	}), resp.Header)
}

func TestSimpleWithRequestHeaderTransfer(t *testing.T) {
	mock = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "hello from client", r.Header.Get("Hello"))

		w.WriteHeader(http.StatusAccepted)
	}
	resp := do(validRequest(http.MethodGet, "/", map[string][]string{
		headerPrefix + "Hello": {"hello from client"},
	}, nil))

	assert.Equal(t, http.StatusOK, resp.StatusCode, "The proxies return code should always be 200 in happy case.")
	assert.Equal(t, http.Header(map[string][]string{
		"Content-Length":                {"0"},                     //origin: proxy (because the target body is less than 512b)
		headerPrefix + "Content-Length": {"0"},                     //origin: target (because the target body is less than 512b)
		headerPrefix + "Status-Line":    {"HTTP/1.1 202 Accepted"}, //origin: proxy (concatenation of targets response code)
	}), resp.Header)
}

func TestSimpleWithResponseHeaderTransfer(t *testing.T) {
	mock = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Hello", "FROM TARGET")

		w.WriteHeader(http.StatusAccepted)
	}
	resp := do(validRequest(http.MethodGet, "/", nil, nil))

	assert.Equal(t, http.StatusOK, resp.StatusCode, "The proxies return code should always be 200 in happy case.")
	assert.Equal(t, http.Header(map[string][]string{
		"Content-Length":                {"0"},                     //origin: proxy (because the target body is less than 512b)
		headerPrefix + "Content-Length": {"0"},                     //origin: target (because the target body is less than 512b)
		headerPrefix + "Hello":          {"FROM TARGET"},           //origin: target
		headerPrefix + "Status-Line":    {"HTTP/1.1 202 Accepted"}, //origin: proxy (concatenation of targets response code)
	}), resp.Header)
}

func TestSimpleWithBodyTransfer(t *testing.T) {
	mock = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		rBody, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, `body from client`, string(rBody))

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`body from target`))
	}
	resp := do(validRequest(http.MethodPost, "/", nil, strings.NewReader(`body from client`)))

	assert.Equal(t, http.StatusOK, resp.StatusCode, "The proxies return code should always be 200 in happy case.")
	assert.Equal(t, http.Header(map[string][]string{
		"Content-Length":                {"16"},                        //origin: proxy (because the target body is less than 512b)
		headerPrefix + "Content-Length": {"16"},                        //origin: target (because the target body is less than 512b)
		"Content-Type":                  {"text/plain; charset=utf-8"}, //origin: proxy (because the target body is less than 512b)
		headerPrefix + "Content-Type":   {"text/plain; charset=utf-8"}, //origin: target (because the target body is less than 512b)
		headerPrefix + "Status-Line":    {"HTTP/1.1 202 Accepted"},     //origin: proxy (concatenation of targets response code)
	}), resp.Header)

	rBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `body from target`, string(rBody))
}
