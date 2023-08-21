package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	})
	lambda.StartWithOptions(handle, lambda.WithContext(context.Background()))
}

func handle(ctx context.Context, request events.APIGatewayProxyRequest) (any, error) {
	logrus.Info("Hello from lambda")

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}
