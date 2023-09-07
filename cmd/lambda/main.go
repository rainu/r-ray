package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"time"
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

	bb := bytes.NewBuffer([]byte("popel"))

	go func() {
		time.Sleep(1 * time.Second)

		_, err := bb.WriteString("penis")
		if err != nil {
			fmt.Println(err)
		}
	}()

	return bb, nil
}
