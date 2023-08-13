package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func (p *proxy) validateRequest(ctx *context) bool {
	targetUrl, err := url.PathUnescape(strings.TrimPrefix(ctx.request.RequestURI, "/"))
	if err != nil {
		ctx.response.WriteHeader(http.StatusBadRequest)
		writeError(ctx.response, fmt.Errorf("invalid url: %w", err))
		return false
	}

	parsedUrl, err := url.ParseRequestURI(targetUrl)
	if err != nil {
		ctx.response.WriteHeader(http.StatusBadRequest)
		writeError(ctx.response, fmt.Errorf("invalid url: %w", err))
		return false
	}

	if parsedUrl.Host == "" || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		ctx.response.WriteHeader(http.StatusBadRequest)
		writeError(ctx.response, fmt.Errorf("invalid url"))
		return false
	}

	ctx.input.URL = parsedUrl.String()

	return true
}

func (p *proxy) compileForwardExpressions(ctx *context) bool {
	ctx.forwardRequestExpressions = make([]*regexp.Regexp, 0, 5)
	for name, value := range ctx.request.Header {
		if !strings.HasPrefix(strings.ToLower(name), strings.ToLower(p.forwardRequestHeader)) {
			//skip headers with no prefix
			continue
		}

		expr, err := regexp.Compile(strings.Join(value, " "))
		if err != nil {
			ctx.response.WriteHeader(http.StatusBadRequest)
			writeError(ctx.response, fmt.Errorf("invalid forward request header expression: %w", err))
			return false
		}
		ctx.forwardRequestExpressions = append(ctx.forwardRequestExpressions, expr)
	}

	ctx.forwardResponseExpressions = make([]*regexp.Regexp, 0, 5)
	for name, value := range ctx.request.Header {
		if !strings.HasPrefix(strings.ToLower(name), strings.ToLower(p.forwardResponseHeader)) {
			//skip headers with no prefix
			continue
		}

		expr, err := regexp.Compile(strings.Join(value, " "))
		if err != nil {
			ctx.response.WriteHeader(http.StatusBadRequest)
			writeError(ctx.response, fmt.Errorf("invalid forward response header expression: %w", err))
			return false
		}
		ctx.forwardResponseExpressions = append(ctx.forwardResponseExpressions, expr)
	}

	return true
}

func (p *proxy) extractAuthorization(ctx *context) bool {
	ctx.input.User.Username, ctx.input.User.Password, _ = ctx.request.BasicAuth()

	return true
}

func (p *proxy) transferRequestHeader(ctx *context) bool {
	for name, values := range ctx.request.Header {
		if !strings.HasPrefix(strings.ToLower(name), strings.ToLower(p.headerPrefix)) {
			//skip headers with no prefix
			continue
		}

		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(p.forwardRequestHeader)) {
			//skip special header
			continue
		}
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(p.forwardResponseHeader)) {
			//skip special header
			continue
		}

		if strings.ToLower(name) == strings.ToLower(p.forwardRequestHeader) ||
			strings.ToLower(name) == strings.ToLower(p.forwardResponseHeader) {
			//those are special headers and should not be transmitted
			continue
		}

		header := name[len(p.headerPrefix):]
		ctx.input.Header[header] = values
	}

	return true
}

func (p *proxy) transferForwardRequestHeader(ctx *context) bool {
	// the client can define headers which should be transfer 1:1 to the target server
	// for this the client can use the ForwardRequestHeader
	// the values of this header are interpreted as regular expression

	if len(ctx.forwardRequestExpressions) == 0 {
		return true
	}

	for hName, hValues := range ctx.request.Header {
		lhName := strings.ToLower(hName)

		if lhName == strings.ToLower(p.forwardRequestHeader) ||
			lhName == strings.ToLower(p.forwardResponseHeader) {
			//those are special headers and should be ignored
			continue
		}

		for _, expr := range ctx.forwardRequestExpressions {
			if expr.MatchString(lhName) {
				//this header should be forwarded
				ctx.input.Header[hName] = hValues
				break //dont need to check other expr.
			}
		}
	}

	return true
}
