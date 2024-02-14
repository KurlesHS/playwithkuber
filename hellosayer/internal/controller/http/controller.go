package http

import (
	"context"
	"net/http"
)

type HelloSayerUserCase interface {
	SayHello(ctx context.Context, name string) string
}

type HandlerRegistrar interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

type Controller struct {
	HelloSayerUserCase
	mux HandlerRegistrar
}

func New(us HelloSayerUserCase, mux HandlerRegistrar) *Controller {
	res := &Controller{
		HelloSayerUserCase: us,
		mux:                mux,
	}
	mux.HandleFunc("/hello/{name}", res.SayHello)
	return res
}
