package controller

import (
	"fmt"
	"github.com/rainu/r-ray/internal/processor"
	"github.com/sirupsen/logrus"
	"net/http"
)

type processingStep func(*context) bool

type proxy struct {
	headerPrefix          string
	forwardRequestHeader  string
	forwardResponseHeader string

	preProcessing  []processingStep
	processor      Processor
	postProcessing []processingStep
}

func NewProxy(headerPrefix, forwardRequestHeader, forwardResponseHeader string, processor Processor) *proxy {
	result := &proxy{
		headerPrefix:          headerPrefix,
		forwardRequestHeader:  forwardRequestHeader,
		forwardResponseHeader: forwardResponseHeader,

		processor: processor,
	}
	result.preProcessing = []processingStep{
		result.validateRequest,
		result.checkAuthentication,
		result.compileForwardExpressions,
		result.transferForwardRequestHeader,
		result.transferRequestHeader,
	}
	result.postProcessing = []processingStep{
		result.checkProcessingErrors,
		result.transferResponseHeader,
		result.transferForwardResponseHeader,
		result.transferStatusCode,
		result.copyBody,
	}

	return result
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	input := processor.Input{
		Header: map[string][]string{},
		Method: r.Method,
		Body:   r.Body,
	}
	processingCtx := &context{
		response: w,
		request:  r,
		input:    &input,
	}

	for _, preProcessingStep := range p.preProcessing {
		if !preProcessingStep(processingCtx) {
			//the current step stops the processing
			return
		}
	}

	// do proxy call
	var output processor.Output
	output, processingCtx.processError = p.processor.Process(input)
	processingCtx.output = &output

	defer func() {
		if output.Body != nil {
			if err := output.Body.Close(); err != nil {
				logrus.WithError(err).Warn("Unable to close request body.")
			}
		}
	}()

	for _, postProcessingStep := range p.postProcessing {
		if !postProcessingStep(processingCtx) {
			//the current step stops the processing
			return
		}
	}
}

func writeError(w http.ResponseWriter, err error) (int, error) {
	return w.Write([]byte(fmt.Sprintf(`{"error":%q}`, err.Error())))
}
