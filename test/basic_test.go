package test

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestSimpleHappy(t *testing.T) {
	mock = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Del("Date")
		w.Header().Set("Hello", "FROM TARGET")

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`<content from target>`))
	}
	resp := do(validRequest(http.MethodGet, "/", nil, nil))

	assert.Equal(t, http.StatusOK, resp.StatusCode, "The proxies return code should always be 200 in happy case.")
	assert.Equal(t, http.Header(map[string][]string{
		"Content-Length":                {"21"},                        //origin: proxy (because the target body is less than 512b)
		headerPrefix + "Content-Length": {"21"},                        //origin: target (because the target body is less than 512b)
		"Content-Type":                  {"text/plain; charset=utf-8"}, //origin: proxy (because the target body is less than 512b)
		headerPrefix + "Content-Type":   {"text/plain; charset=utf-8"}, //origin: target (because the target body is less than 512b)
		headerPrefix + "Hello":          {"FROM TARGET"},               //origin: target
		headerPrefix + "Status-Line":    {"HTTP/1.1 202 Accepted"},     //origin: proxy (concatenation of targets response code)
	}), resp.Header)

	rBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `<content from target>`, string(rBody))
}
