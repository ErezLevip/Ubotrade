package AuthenticationLogging

import (
	"time"
	"fmt"
	"context"

	"github.com/go-kit/kit/log"
	"github.com/erezlevip/Ubotrade/Global"
	"github.com/erezlevip/Ubotrade/AuthenticationService/Service"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   AuthenticationService.IAuthenticationService
}



func (mw LoggingMiddleware)  Login(ctx context.Context, clientId string, firstName string, lastName string, email string, sessionId string) (userId string,err error)  {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "Login",
			"input", fmt.Sprintf("%b", clientId,sessionId,firstName,lastName,email),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	userId, err = mw.Next.Login(ctx,clientId,firstName,lastName,email,sessionId)
	return
}

func (mw LoggingMiddleware)  GetToken(ctx context.Context, clientId string) (token string,err error)  {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetToken",
			"input", fmt.Sprintf("%b", clientId),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	token, err = mw.Next.GetToken(ctx, clientId)
	return
}

func (mw LoggingMiddleware) ValidateToken(ctx context.Context, token string) (IsValid bool, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "ValidateToken",
			"input", fmt.Sprintf("%b", token),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	IsValid, err = mw.Next.ValidateToken(ctx, token)
	return
}

func (mw LoggingMiddleware) Init(serviceInfo Global.ServiceInformation) *AuthenticationService.AuthenticationService {
	return mw.Next.Init(serviceInfo)
}

func (mw LoggingMiddleware) GetServiceMetrics() (map[string]interface{}, error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetServiceMetrics",
			"input", nil,
			"err", nil,
			"took", time.Since(begin),
		)
	}(time.Now())
	return make(map[string]interface{}), nil
}
