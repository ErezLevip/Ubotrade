package RegistryServiceLogging

import (
	"time"

	"github.com/go-kit/kit/log"

	"fmt"
	"local/UbotTrade/DataHandlers/Redis"
	"local/UbotTrade/Global"
	"local/UbotTrade/RegistryService/Service"
	"context"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   RegistryService.IRegistryService
}

func (mw LoggingMiddleware) Register(ctx context.Context, serviceInformation Global.ServiceInformation) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "Register",
			"input", fmt.Sprintf("%b", serviceInformation),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.Next.Register(ctx, serviceInformation)
	return
}

func (mw LoggingMiddleware) DeRegister(ctx context.Context ,serviceInformation Global.ServiceInformation) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "DeRegister",
			"input", fmt.Sprintf("%b", serviceInformation),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.Next.DeRegister(ctx, serviceInformation)
	return
}

func (mw LoggingMiddleware) GetService(ctx context.Context, serviceName string) (serviceInformation Global.ServiceInformation, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetService",
			"input", serviceName,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	serviceInformation, err = mw.Next.GetService(ctx, serviceName)
	return
}

func (mw LoggingMiddleware) Init(config DataHandlers.RedisConfiguration) *RegistryService.RegistryService {
	return mw.Next.Init(config)
}
