package UserLogging

import (
	"time"
	"fmt"
	"context"

	"github.com/go-kit/kit/log"
	"github.com/erezlevip/Ubotrade/Global"
	"github.com/erezlevip/Ubotrade/UserService/Service"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   UserService.IUserService
}

func (mw LoggingMiddleware) GetUser(ctx context.Context, userId string, dataType string, activeOnly bool) (data []map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetUser",
			"input", fmt.Sprintf("%b", userId, dataType),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	data, err = mw.Next.GetUser(ctx, userId, dataType, activeOnly)
	return
}

func (mw LoggingMiddleware) CreateUser(ctx context.Context, userId string, firstName string, lastName string, email string) (data map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "CreateUser",
			"input", fmt.Sprintf("%b", userId, firstName, lastName, email),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	data, err = mw.Next.CreateUser(ctx, userId, firstName, lastName, email)
	return
}

func (mw LoggingMiddleware) SetUser(ctx context.Context, userId string, dataType string, operation string, data map[string]interface{}) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "SetUser",
			"input", fmt.Sprintf("%b", userId, dataType, data, operation),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.Next.SetUser(ctx, userId, dataType, operation, data)
	return
}

func (mw LoggingMiddleware) Init(serviceInfo Global.ServiceInformation) *UserService.UserService {
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
