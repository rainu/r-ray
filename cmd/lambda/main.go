package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rainu/r-ray/internal/config"
	ihttp "github.com/rainu/r-ray/internal/http"
	"github.com/rainu/r-ray/internal/http/controller"
	"github.com/rainu/r-ray/internal/processor"
	"github.com/rainu/r-ray/internal/store"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
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
	//httpReq, err := http.NewRequestWithContext(
	//	ctx,
	//	request.RequestContext.HTTP.Method,
	//	path.Join("https://lambda.aws.com", request.RawPath),
	//	strings.NewReader(request.Body),
	//)
	//
	//if err != nil {
	//	return nil, fmt.Errorf("unable to build request: %w", err)
	//}
	//
	//for hName, hValues := range request.Headers {
	//	httpReq.Header.Set(hName, hValues)
	//	httpReq.Header.Add()
	//}
	//
	//
	//request.Headers
	//
	//
	//httpResp := httptest.NewRecorder()
	//
	//proxy.ServeHTTP(httpResp, httpReq)
	//
	//return events.ApiResp

	logrus.WithField("request", request).Info("Call")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
