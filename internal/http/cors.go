package http

import (
	"net/http"
	"strconv"
	"strings"
)

type CorsMiddleware struct {
	Delegate http.Handler

	Origins []string
	Methods []string
	Headers []string
	MaxAge  int
}

func (c CorsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(c.Origins) == 0 {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(c.Origins, ", "))
	}

	if len(c.Methods) == 0 {
		w.Header().Set("Access-Control-Allow-Methods", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.Methods, ", "))
	}

	if len(c.Headers) == 0 {
		w.Header().Set("Access-Control-Allow-Headers", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.Headers, ", "))
	}

	if c.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(c.MaxAge))
	}

	c.Delegate.ServeHTTP(w, r)
}
