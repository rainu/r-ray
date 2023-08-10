package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestValidationBadUrl(t *testing.T) {
	cases := []struct {
		url         string
		expectedErr string
	}{
		{
			url:         appBaseUrl,
			expectedErr: `{"error":"invalid url: parse \"\": empty url"}`,
		},
		{
			url:         appBaseUrl + "//",
			expectedErr: `{"error":"invalid url"}`,
		},
		{
			url:         appBaseUrl + "/brOken",
			expectedErr: `{"error":"invalid url"}`,
		},
		{
			url:         appBaseUrl + "http://host:port/test",
			expectedErr: `{"error":"invalid url: parse \"http://host:port/test\": invalid port \":port\" after host"}`,
		},
	}

	mock = nil
	for i, tt := range cases {
		t.Run(fmt.Sprintf("TestValidationBadUrl_%d", i), func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			if err != nil {
				panic(err)
			}
			resp := do(req)

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			rBody, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tt.expectedErr, string(rBody))
		})
	}
}
