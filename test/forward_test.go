package test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRequestHeaderForward(t *testing.T) {
	mock = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, `hello to proxy`, r.Header.Get("Hello"))

		w.WriteHeader(http.StatusAccepted)
	}
	resp := do(validRequest(http.MethodGet, "/", map[string][]string{
		"Hello":                  {"hello to proxy"}, //this should be forwarded because of forward header (not header-prefix)
		fwdReqHeaderPrefix + "0": {`^he.*`},          //regex which should match
	}, nil))

	assert.Equal(t, http.StatusOK, resp.StatusCode, "The proxies return code should always be 200 in happy case.")
	assert.Equal(t, http.Header(map[string][]string{
		"Content-Length":                {"0"},                     //origin: proxy (because the target body is less than 512b)
		headerPrefix + "Content-Length": {"0"},                     //origin: target (because the target body is less than 512b)
		headerPrefix + "Status-Line":    {"HTTP/1.1 202 Accepted"}, //origin: proxy (concatenation of targets response code)
	}), resp.Header)
}

func TestRequestHeaderForward_LowerThanDedicated(t *testing.T) {
	mock = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, `hello to target`, r.Header.Get("Hello"))

		w.WriteHeader(http.StatusAccepted)
	}
	resp := do(validRequest(http.MethodGet, "/", map[string][]string{
		"Hello":                  {"hello to proxy"},  //this should not be forwarded because there is a dedicated header
		fwdReqHeaderPrefix + "0": {`^he.*`},           //regex which should match
		headerPrefix + "Hello":   {"hello to target"}, //this header should be used!
	}, nil))

	assert.Equal(t, http.StatusOK, resp.StatusCode, "The proxies return code should always be 200 in happy case.")
	assert.Equal(t, http.Header(map[string][]string{
		"Content-Length":                {"0"},                     //origin: proxy (because the target body is less than 512b)
		headerPrefix + "Content-Length": {"0"},                     //origin: target (because the target body is less than 512b)
		headerPrefix + "Status-Line":    {"HTTP/1.1 202 Accepted"}, //origin: proxy (concatenation of targets response code)
	}), resp.Header)
}

func TestResponseHeaderForward(t *testing.T) {
	mock = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}
	resp := do(validRequest(http.MethodGet, "/", map[string][]string{
		fwdRespHeaderPrefix + "0": {`^content-.*`},
	}, nil))

	assert.Equal(t, http.StatusOK, resp.StatusCode, "The proxies return code should always be 200 in happy case.")
	assert.Equal(t, http.Header(map[string][]string{
		"Content-Length":             {"0"},                     //origin: target
		headerPrefix + "Status-Line": {"HTTP/1.1 202 Accepted"}, //origin: proxy (concatenation of targets response code)
	}), resp.Header)
}
