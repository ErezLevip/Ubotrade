package ConfigurationServiceInstrumenting

import (
	"fmt"
	"time"
	"context"

	"github.com/go-kit/kit/metrics"
	"github.com/erezlevip/Ubotrade/ConfigurationService/Service"
	"github.com/erezlevip/Ubotrade/Global"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
	Next           ConfigurationService.IConfigurationService
}

var FailureCount int
var RequestCount int

func Make(requestCount metrics.Counter, requestLatency metrics.Histogram, countResult metrics.Histogram, next ConfigurationService.IConfigurationService) (instance *InstrumentingMiddleware) {
	FailureCount = 0
	RequestCount = 0
	return &InstrumentingMiddleware{
		RequestCount:   requestCount,
		CountResult:    countResult,
		RequestLatency: requestLatency,
		Next:           next,
	}
}

func (mw InstrumentingMiddleware) GetConfiguration(ctx context.Context, key string) (configuration map[string]interface{}, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetConfiguration", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	configuration, err = mw.Next.GetConfiguration(ctx, key)
	return
}

func (mw InstrumentingMiddleware) GetServiceMetrics() (metrics map[string]interface{}, err error) {
	metrics = map[string]interface{}{
		"SuccessRate":  getSuccessRate(),
		"RequestCount": RequestCount,
	}
	err = nil
	return
}

func (mw InstrumentingMiddleware) Init(staticConfigPath string, serviceInfo Global.ServiceInformation) *ConfigurationService.ConfigurationService {
	return mw.Next.Init(staticConfigPath, serviceInfo)
}

func getSuccessRate() float64 {
	if RequestCount == 0 || FailureCount == 0 {
		return 100.0
	}

	return float64((1 - (FailureCount / RequestCount)) * 100)
}
