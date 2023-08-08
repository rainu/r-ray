package processor

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type requestLogger struct {
	delegate httpClient
}

func (r requestLogger) Do(req *http.Request) (res *http.Response, err error) {
	t := time.Now()
	defer func() {
		duration := time.Now().Sub(t)

		log := logrus.
			WithField("method", req.Method).
			WithField("url", req.URL).
			WithField("duration", duration)

		if res != nil {
			log = log.WithField("status", res.Status).
				WithField("length", res.ContentLength)
		}
		if err != nil {
			log = log.WithError(err)
		}

		log.Info("Http call done.")
	}()

	res, err = r.delegate.Do(req)
	return
}
