package endpoints

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"

	"github.com/knstch/subtrack-libs/transport"
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

		if len(ep.Opts) != 0 {
			ep.Opts = append(ep.Opts, httptransport.ServerErrorEncoder(transport.EncodeError))
		} else {
			ep.Opts = []httptransport.ServerOption{httptransport.ServerErrorEncoder(transport.EncodeError)}
		}

		opts := append(ep.Opts, httptransport.ServerErrorEncoder(transport.EncodeError))

		r.Methods(ep.Method).Path(ep.Path).Handler(httptransport.NewServer(
			handler,
			ep.Decoder,
			ep.Encoder,
			opts...,
		))
	}

	return r
}
