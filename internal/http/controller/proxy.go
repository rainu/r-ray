package controller

import (
	"errors"
	"fmt"
	"github.com/rainu/r-ray/internal/processor"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type proxy struct {
	headerPrefix string

	processor Processor
}

func NewProxy(headerPrefix string, processor Processor) *proxy {
	return &proxy{
		headerPrefix: strings.ToLower(headerPrefix),
		processor:    processor,
	}
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	input := processor.Input{
		Header: map[string][]string{},
		Method: r.Method,
		Body:   r.Body,
	}

	parsedUrl, err := url.ParseRequestURI(r.URL.Query().Get("url"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, fmt.Errorf("invalid url: %w", err))
		return
	}

	if parsedUrl.Host == "" || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, fmt.Errorf("invalid url"))
		return
	}

	input.URL = parsedUrl.String()

	var ok bool
	input.User.Username, input.User.Password, ok = r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	for name, values := range r.Header {
		if !strings.HasPrefix(strings.ToLower(name), p.headerPrefix) {
			//skip headers with no prefix
			continue
		}

		header := name[len(p.headerPrefix):]
		input.Header[header] = values
	}

	// do proxy call
	output, err := p.processor.Process(input)
	if err != nil && errors.Is(err, processor.ErrUnauthorized) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return
	}

	defer func() {
		if output.Body != nil {
			if err := output.Body.Close(); err != nil {
				logrus.WithError(err).Warn("Unable to close request body.")
			}
		}
	}()

	for name, values := range output.Header {
		w.Header()[p.headerPrefix+name] = values
	}
	w.Header()[p.headerPrefix+`Status-Line`] = []string{output.StatusLine}

	if _, err := io.Copy(w, output.Body); err != nil {
		logrus.WithError(err).Warn("Unable to copy body content.")
	}
}

func writeError(w http.ResponseWriter, err error) (int, error) {
	return w.Write([]byte(fmt.Sprintf(`{"error":%q}`, err.Error())))
}
