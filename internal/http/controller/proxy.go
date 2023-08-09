package controller

import (
	"fmt"
	"github.com/rainu/r-ray/internal/processor"
	"github.com/sirupsen/logrus"
	"net/http"
)

type preProcessing func(http.ResponseWriter, *http.Request, *processor.Input) bool
type postProcessing func(http.ResponseWriter, *http.Request, *processor.Input, *processor.Output, error) bool

type proxy struct {
	headerPrefix          string
	forwardRequestHeader  string
	forwardResponseHeader string

	preProcessing  []preProcessing
	processor      Processor
	postProcessing []postProcessing
}

func NewProxy(headerPrefix, forwardRequestHeader, forwardResponseHeader string, processor Processor) *proxy {
	result := &proxy{
		headerPrefix:          headerPrefix,
		forwardRequestHeader:  forwardRequestHeader,
		forwardResponseHeader: forwardResponseHeader,

		processor: processor,
	}
	result.preProcessing = []preProcessing{
		result.validateRequest,
		result.checkAuthentication,
		result.transferRequestHeader,
	}
	result.postProcessing = []postProcessing{
		result.checkProcessingErrors,
		result.transferResponseHeader,
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

	for _, preProcessingStep := range p.preProcessing {
		if !preProcessingStep(w, r, &input) {
			//the current step stops the processing
			return
		}
	}

	// TODO: forward request headers...
	//if frhValues := r.Header.Values(p.forwardRequestHeader); len(frhValues) != 0 {
	//	//the client wants that the given header from this request will be forwarded to the target
	//	for _, frhRegex := range frhValues {
	//		for name, values := range r.Header {
	//		if values := r.Header.Values(fwdHeader); len(values) != 0 {
	//			input.Header[fwdHeader] = values
	//		}
	//	}
	//}

	// do proxy call
	output, err := p.processor.Process(input)
	defer func() {
		if output.Body != nil {
			if err := output.Body.Close(); err != nil {
				logrus.WithError(err).Warn("Unable to close request body.")
			}
		}
	}()

	for _, postProcessingStep := range p.postProcessing {
		if !postProcessingStep(w, r, &input, &output, err) {
			//the current step stops the processing
			return
		}
	}

	// TODO: forward response headers...
}

func writeError(w http.ResponseWriter, err error) (int, error) {
	return w.Write([]byte(fmt.Sprintf(`{"error":%q}`, err.Error())))
}
