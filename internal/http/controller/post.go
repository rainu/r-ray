package controller

import (
	"errors"
	"github.com/rainu/r-ray/internal/processor"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

const (
	StatusLineHeaderSuffix = `Status-Line`
)

func (p *proxy) checkProcessingErrors(ctx *context) bool {
	if ctx.processError != nil && errors.Is(ctx.processError, processor.ErrUnauthorized) {
		ctx.response.WriteHeader(http.StatusUnauthorized)
		return false
	} else if ctx.processError != nil {
		ctx.response.WriteHeader(http.StatusInternalServerError)
		writeError(ctx.response, ctx.processError)
		return false
	}

	return true
}

func (p *proxy) transferResponseHeader(ctx *context) bool {
	for name, values := range ctx.output.Header {
		ctx.response.Header()[p.headerPrefix+name] = values
	}

	return true
}

func (p *proxy) transferForwardResponseHeader(ctx *context) bool {
	// the client can define headers which should be transfer 1:1 from the target server
	// for this the client can use the ForwardResponseHeader
	// the values of this header are interpreted as regular expression

	if len(ctx.forwardResponseExpressions) == 0 {
		return true
	}

	for hName, hValues := range ctx.output.Header {
		lhName := strings.ToLower(hName)

		for _, expr := range ctx.forwardResponseExpressions {
			if expr.MatchString(lhName) {
				//this header should be forwarded
				ctx.response.Header()[hName] = hValues

				//remove the suffixed header
				ctx.response.Header().Del(p.headerPrefix + hName)

				break //dont need to check other expr.
			}
		}
	}

	return true
}

func (p *proxy) transferStatusCode(ctx *context) bool {
	if ctx.forwardResponseStatus {
		ctx.response.WriteHeader(ctx.output.StatusCode)
	} else {
		ctx.response.Header()[p.headerPrefix+StatusLineHeaderSuffix] = []string{ctx.output.StatusLine}
	}

	return true
}

func (p *proxy) copyBody(ctx *context) bool {
	if _, err := io.Copy(ctx.response, ctx.output.Body); err != nil {
		logrus.WithError(err).Warn("Unable to copy body content.")
	}

	return true
}
