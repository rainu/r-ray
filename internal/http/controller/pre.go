package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func (p *proxy) validateRequest(ctx *context) bool {
	parsedUrl, err := url.ParseRequestURI(ctx.request.URL.Query().Get("url"))
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
	rawAllowExpressions := ctx.request.Header.Values(p.forwardRequestHeader)
	ctx.forwardRequestExpressions = make([]*regexp.Regexp, len(rawAllowExpressions))

	for i, expr := range rawAllowExpressions {
		var err error
		ctx.forwardRequestExpressions[i], err = regexp.Compile(expr)
		if err != nil {
			ctx.response.WriteHeader(http.StatusBadRequest)
			writeError(ctx.response, fmt.Errorf("invalid forward request header expression: %w", err))
			return false
		}
	}

	rawAllowExpressions = ctx.request.Header.Values(p.forwardResponseHeader)
	ctx.forwardResponseExpressions = make([]*regexp.Regexp, len(rawAllowExpressions))

	for i, expr := range rawAllowExpressions {
		var err error
		ctx.forwardResponseExpressions[i], err = regexp.Compile(expr)
		if err != nil {
			ctx.response.WriteHeader(http.StatusBadRequest)
			writeError(ctx.response, fmt.Errorf("invalid forward response header expression: %w", err))
			return false
		}
	}

	return true
}

func (p *proxy) checkAuthorization(ctx *context) bool {
	var ok bool
	ctx.input.User.Username, ctx.input.User.Password, ok = ctx.request.BasicAuth()
	if !ok {
		ctx.response.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}

func (p *proxy) transferRequestHeader(ctx *context) bool {
	for name, values := range ctx.request.Header {
		if !strings.HasPrefix(strings.ToLower(name), strings.ToLower(p.headerPrefix)) {
			//skip headers with no prefix
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
