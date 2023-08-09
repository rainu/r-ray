package controller

import (
	"errors"
	"github.com/rainu/r-ray/internal/processor"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (p *proxy) checkProcessingErrors(w http.ResponseWriter, r *http.Request, i *processor.Input, o *processor.Output, err error) bool {
	if err != nil && errors.Is(err, processor.ErrUnauthorized) {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return false
	}

	return true
}

func (p *proxy) transferResponseHeader(w http.ResponseWriter, r *http.Request, i *processor.Input, o *processor.Output, _ error) bool {
	for name, values := range o.Header {
		w.Header()[p.headerPrefix+name] = values
	}

	return true
}

func (p *proxy) transferStatusCode(w http.ResponseWriter, r *http.Request, i *processor.Input, o *processor.Output, _ error) bool {
	w.Header()[p.headerPrefix+`Status-Line`] = []string{o.StatusLine}

	return true
}

func (p *proxy) copyBody(w http.ResponseWriter, r *http.Request, i *processor.Input, o *processor.Output, _ error) bool {
	if _, err := io.Copy(w, o.Body); err != nil {
		logrus.WithError(err).Warn("Unable to copy body content.")
	}

	return true
}
