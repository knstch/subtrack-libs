package endpoints

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

type Middleware func(endpoint.Endpoint) endpoint.Endpoint

type Endpoint struct {
	Method  string
	Path    string
	Handler endpoint.Endpoint
	Decoder httptransport.DecodeRequestFunc
	Encoder httptransport.EncodeResponseFunc
	Req     interface{}
	Res     interface{}
	Mdw     []Middleware
	Opts    []httptransport.ServerOption
}

func InitHttpEndpoints(endpoints []Endpoint) http.Handler {
	r := mux.NewRouter()

	for _, ep := range endpoints {
		handler := ep.Handler
		for _, mw := range ep.Mdw {
			handler = mw(handler)
		}

		r.Methods(ep.Method).Path(ep.Path).Handler(httptransport.NewServer(
			handler,
			ep.Decoder,
			ep.Encoder,
			ep.Opts...,
		))
	}

	return r
}
