package Middlewares

import (
	"context"
	"local/UbotTrade/API/AuthenticationServiceAPI"
	"local/UbotTrade/API/ServiceAPIFactory"
	"local/UbotTrade/Global"
	"net/http"
	"reflect"
)

type AuthMiddleware struct {
}

func (m *AuthMiddleware) Next(w http.ResponseWriter, req *http.Request, ctx context.Context) (http.ResponseWriter, *http.Request, context.Context) {
	if req != nil {
		var sessionId = req.Header.Get("X-AUTHORIZATION")
		if (sessionId != "") {
			authSvc := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(AuthenticationServiceAPI.AuthenticationServiceAPI{})).(*AuthenticationServiceAPI.AuthenticationServiceAPI)

			res, err := authSvc.Login(Global.LoginRequest{SessionId: sessionId})
			if res.UserId == "" || err != nil {
				w.WriteHeader(res.Base.Status)
			} else {
				ctx = context.WithValue(ctx, "UserId", res.UserId)
			}
		}
	}
	return w, req, ctx
}

func (m *AuthMiddleware) Register() {
	Middlewares = append(Middlewares, m)
}

func AuthMiddleWareMake() *AuthMiddleware {
	return &AuthMiddleware{}
}
