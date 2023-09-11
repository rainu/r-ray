package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rainu/r-ray/internal/config"
	ihttp "github.com/rainu/r-ray/internal/http"
	"github.com/rainu/r-ray/internal/http/controller"
	"github.com/rainu/r-ray/internal/processor"
	"github.com/rainu/r-ray/internal/store"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.ReadConfig()
	if err != nil {
		logrus.WithError(err).Error("Error while reading config.")
		os.Exit(1)
		return
	}

	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	userStore := store.NewUser()
	for _, credential := range cfg.Credentials {
		userStore.Add(credential.UsernameAndPassword())
	}

	p := processor.New(userStore)
	proxy = ihttp.CorsMiddleware{
		Delegate: ihttp.NewMetaMiddleware(cfg, controller.NewProxy(
			cfg.RequestHeaderPrefix,
			cfg.ForwardRequestHeaderPrefix,
			cfg.ForwardResponseHeaderPrefix,
			cfg.ForwardResponseStatusHeader,
			p,
		)),

		Origins: cfg.CorsAllowOrigin,
		Methods: cfg.CorsAllowMethods,
		Headers: cfg.CorsAllowHeaders,
		MaxAge:  cfg.CorsAllowMaxAge,
	}
	lambda.StartWithOptions(handle, lambda.WithContext(context.Background()))
}

var proxy http.Handler

func handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (any, error) {
	var body io.Reader = strings.NewReader(request.Body)
	if request.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, body)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		request.RequestContext.HTTP.Method,
		fmt.Sprintf("https://lambda.aws.com%s?%s", request.RawPath, request.RawQueryString),
		body,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	for hName, hValues := range request.Headers {
		httpReq.Header.Add(hName, hValues)
	}
	for _, value := range request.Cookies {
		httpReq.Header.Add("Cookie", value)
	}

	httpResp := newRecorder()
	proxy.ServeHTTP(httpResp, httpReq)
	httpResp.b64.Close()

	response := events.APIGatewayV2HTTPResponse{
		StatusCode:        httpResp.status,
		MultiValueHeaders: httpResp.Header(),
		Body:              httpResp.buffer.String(),
		IsBase64Encoded:   true,
	}
	return response, nil
}

type responseRecorder struct {
	header http.Header
	status int
	buffer *bytes.Buffer
	b64    io.WriteCloser
}

func newRecorder() responseRecorder {
	result := responseRecorder{
		buffer: bytes.NewBuffer(nil),
	}
	result.b64 = base64.NewEncoder(base64.StdEncoding, result.buffer)

	return result
}

func (r responseRecorder) Header() http.Header {
	return r.header
}

func (r responseRecorder) Write(bytes []byte) (int, error) {
	return r.b64.Write(bytes)
}

func (r responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
}
