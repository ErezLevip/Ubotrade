package RegistryServiceInstrumenting

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"local/UbotTrade/DataHandlers/Redis"
	"local/UbotTrade/Global"
	"local/UbotTrade/RegistryService/Service"
	"context"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
	Next           RegistryService.IRegistryService
}

func (mw InstrumentingMiddleware) Register(ctx context.Context, serviceInformation Global.ServiceInformation) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Register", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Next.Register(ctx, serviceInformation)
	return
}

func (mw InstrumentingMiddleware) DeRegister(ctx context.Context, serviceInformation Global.ServiceInformation) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "DeRegister", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Next.DeRegister(ctx, serviceInformation)
	return
}

func (mw InstrumentingMiddleware) GetService(ctx context.Context, serviceName string) (serviceInformation Global.ServiceInformation, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetService", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	serviceInformation, err = mw.Next.GetService(ctx, serviceName)
	return
}

func (mw InstrumentingMiddleware) Init(config DataHandlers.RedisConfiguration) *RegistryService.RegistryService {
	return mw.Next.Init(config)
}
