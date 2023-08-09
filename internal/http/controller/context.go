package controller

import (
	"github.com/rainu/r-ray/internal/processor"
	"net/http"
	"regexp"
)

type context struct {
	response http.ResponseWriter
	request  *http.Request

	input  *processor.Input
	output *processor.Output

	processError error

	forwardRequestExpressions  []*regexp.Regexp
	forwardResponseExpressions []*regexp.Regexp
}
