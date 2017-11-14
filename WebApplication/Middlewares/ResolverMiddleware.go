package Middlewares

import (
	"context"
	"net/http"
	"reflect"

	"github.com/erezlevip/Ubotrade/API/AuthenticationServiceAPI"
	"github.com/erezlevip/Ubotrade/API/BotServiceAPI"
	"github.com/erezlevip/Ubotrade/API/ServiceAPIFactory"
	"github.com/erezlevip/Ubotrade/API/UserServiceAPI"
)

type ResolverMiddleware struct {
}

func (m *ResolverMiddleware) Next(w http.ResponseWriter, req *http.Request, ctx context.Context) (http.ResponseWriter, *http.Request, context.Context) {
	if req != nil {
		ctx = AddServiceToContext(BotServiceAPI.BotServiceAPI{}, ctx)
		ctx = AddServiceToContext(AuthenticationServiceAPI.AuthenticationServiceAPI{}, ctx)
		ctx = AddServiceToContext(UserServiceAPI.UserServiceAPI{}, ctx)
	}
	return w, req, ctx
}

func AddServiceToContext(service interface{}, ctx context.Context) context.Context {
	t := reflect.TypeOf(service)
	svcInstance := ServiceAPIFactory.GetServiceInstance(ctx, t)
	return context.WithValue(ctx, t, svcInstance)
}

func (m *ResolverMiddleware) Register() {
	Middlewares = append(Middlewares, m)
}

func ResolverMiddlewareMake() *ResolverMiddleware {
	return &ResolverMiddleware{}
}
