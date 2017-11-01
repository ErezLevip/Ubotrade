package Middlewares

import (
	"context"
	"local/UbotTrade/API/AuthenticationServiceAPI"
	"local/UbotTrade/API/BotServiceAPI"
	"local/UbotTrade/API/ServiceAPIFactory"
	"local/UbotTrade/API/UserServiceAPI"
	"net/http"
	"reflect"
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
