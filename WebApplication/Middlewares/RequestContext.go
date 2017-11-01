package Middlewares

import (
	"context"
	"net/http"
	"reflect"
)

type requestHandler func(w http.ResponseWriter, req *http.Request, ctx context.Context)

func BeginRequest(fn requestHandler, skipAuth bool) func(w http.ResponseWriter, req *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		var ctx = context.Background()
		for _, middleWare := range Middlewares {
			if skipAuth || !(reflect.TypeOf(middleWare) == reflect.TypeOf(AuthMiddleware{})) {
				w, req, ctx = middleWare.Next(w, req, ctx)
			}
		}

		fn(w, req, ctx)
	}
}
