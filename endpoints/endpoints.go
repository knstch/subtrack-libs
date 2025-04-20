package endpoints

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/knstch/subtrack-libs/middleware"
	"github.com/knstch/subtrack-libs/transport"

	"github.com/go-chi/chi/v5"
)

type Endpoint struct {
	Method  string
	Path    string
	Handler endpoint.Endpoint
	Decoder httptransport.DecodeRequestFunc
	Encoder httptransport.EncodeResponseFunc
	Req     interface{}
	Res     interface{}
	Mdw     []middleware.Middleware
	Opts    []httptransport.ServerOption
}

func InitHttpEndpoints(endpoints []Endpoint) http.Handler {
	r := chi.NewRouter()

	for _, ep := range endpoints {
		handler := ep.Handler
		for _, mw := range ep.Mdw {
			handler = mw(handler)
		}

		opts := append(ep.Opts,
			httptransport.ServerErrorEncoder(transport.EncodeError),
			httptransport.ServerBefore(httptransport.PopulateRequestContext),
		)

		r.Method(ep.Method, ep.Path, httptransport.NewServer(
			handler,
			ep.Decoder,
			ep.Encoder,
			opts...,
		))
	}

	return r
}
