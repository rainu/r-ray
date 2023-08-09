package controller

import (
	"fmt"
	"github.com/rainu/r-ray/internal/processor"
	"net/http"
	"net/url"
	"strings"
)

func (p *proxy) validateRequest(w http.ResponseWriter, r *http.Request, i *processor.Input) bool {
	parsedUrl, err := url.ParseRequestURI(r.URL.Query().Get("url"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, fmt.Errorf("invalid url: %w", err))
		return false
	}

	if parsedUrl.Host == "" || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, fmt.Errorf("invalid url"))
		return false
	}

	i.URL = parsedUrl.String()

	return true
}

func (p *proxy) checkAuthentication(w http.ResponseWriter, r *http.Request, i *processor.Input) bool {
	var ok bool
	i.User.Username, i.User.Password, ok = r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}

func (p *proxy) transferRequestHeader(w http.ResponseWriter, r *http.Request, i *processor.Input) bool {
	for name, values := range r.Header {
		if !strings.HasPrefix(strings.ToLower(name), strings.ToLower(p.headerPrefix)) {
			//skip headers with no prefix
			continue
		}

		if strings.ToLower(name) == strings.ToLower(p.forwardRequestHeader) ||
			strings.ToLower(name) == strings.ToLower(p.forwardResponseHeader) {
			//those are special headers and should not be transmitted
			continue
		}

		header := name[len(p.headerPrefix):]
		i.Header[header] = values
	}

	return true
}
