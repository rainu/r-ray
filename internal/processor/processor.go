package processor

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type processor struct {
	httpClient httpClient

	userStore UserStore
}

func New(userStore UserStore) *processor {
	return &processor{
		httpClient: requestLogger{&http.Client{}},
		userStore:  userStore,
	}
}

func (p *processor) Process(input Input) (Output, error) {
	if !p.userStore.IsValid(input.User.Username, input.User.Password) {
		return Output{}, ErrUnauthorized
	}

	log := logrus.WithField("user", input.User.Username)

	req, err := http.NewRequest(input.Method, input.URL, input.Body)
	if err != nil {
		return Output{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = input.Header

	log.
		WithField("header", req.Header).
		WithField("method", req.Method).
		WithField("url", req.URL).
		Debug("Process request.")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return Output{}, fmt.Errorf("unable to do request: %w", err)
	}

	log.
		WithField("header", resp.Header).
		WithField("method", input.Method).
		WithField("url", input.URL).
		Debug("Process request done.")

	return Output{
		StatusLine: fmt.Sprintf("%s %s", resp.Proto, resp.Status),
		Body:       resp.Body,
		Header:     resp.Header,
	}, nil
}
