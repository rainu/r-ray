package http

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type loggingMiddleware struct {
	delegate http.Handler
}

func (l loggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	defer func() {
		duration := time.Now().Sub(t)

		log := logrus.
			WithField("method", r.Method).
			WithField("url", r.URL).
			WithField("duration", duration)

		log.Info("Http serve done.")
	}()

	l.delegate.ServeHTTP(w, r)
}
