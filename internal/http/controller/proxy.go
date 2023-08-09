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
	headerPrefix          string
	forwardRequestHeader  string
	forwardResponseHeader string

	processor Processor
}

func NewProxy(headerPrefix, forwardRequestHeader, forwardResponseHeader string, processor Processor) *proxy {
	return &proxy{
		headerPrefix:          headerPrefix,
		forwardRequestHeader:  forwardRequestHeader,
		forwardResponseHeader: forwardResponseHeader,

		processor: processor,
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

	// TODO: forward request headers...
	//if frhValues := r.Header.Values(p.forwardRequestHeader); len(frhValues) != 0 {
	//	//the client wants that the given header from this request will be forwarded to the target
	//	for _, frhRegex := range frhValues {
	//		for name, values := range r.Header {
	//		if values := r.Header.Values(fwdHeader); len(values) != 0 {
	//			input.Header[fwdHeader] = values
	//		}
	//	}
	//}

	// TODO: make small peaces of code

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

	// TODO: forward response headers...

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
