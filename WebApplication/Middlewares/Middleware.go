package Middlewares

import (
	"context"
	"net/http"
)

type MiddleWare interface {
	Next(w http.ResponseWriter, req *http.Request, ctx context.Context) (http.ResponseWriter, *http.Request, context.Context)
	Register()
}

var Middlewares = make([]MiddleWare, 0)
