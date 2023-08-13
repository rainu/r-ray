package http

import (
	"encoding/json"
	"github.com/rainu/r-ray/internal/config"
	"github.com/rainu/r-ray/internal/http/controller"
	"net/http"
)

type metaMiddleware struct {
	delegate http.Handler

	meta []byte
}

func NewMetaMiddleware(cfg *config.Config, delegate http.Handler) http.Handler {
	meta, err := json.Marshal(map[string]interface{}{
		"headerPrefix":                cfg.RequestHeaderPrefix,
		"forwardRequestHeaderPrefix":  cfg.ForwardRequestHeaderPrefix,
		"forwardResponseHeaderPrefix": cfg.ForwardResponseHeaderPrefix,
		"forwardResponseStatusHeader": cfg.ForwardResponseStatusHeader,
		"statusHeader":                cfg.RequestHeaderPrefix + controller.StatusLineHeaderSuffix,
	})
	if err != nil {
		panic(err)
	}

	return &metaMiddleware{
		delegate: delegate,
		meta:     meta,
	}
}

func (c metaMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/.meta" {
		c.delegate.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(c.meta)
}
