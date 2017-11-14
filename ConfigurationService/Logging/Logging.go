package ConfigurationServiceLogging

import (
	"time"
	"context"

	"github.com/go-kit/kit/log"
	"github.com/erezlevip/Ubotrade/ConfigurationService/Service"
	"github.com/erezlevip/Ubotrade/Global"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   ConfigurationService.IConfigurationService
}

func (mw LoggingMiddleware) GetConfiguration(ctx context.Context, key string) (configuration map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetConfiguration",
			"input", key,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	configuration, err = mw.Next.GetConfiguration(ctx, key)
	return
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

func (mw LoggingMiddleware) Init(staticConfigPath string, serviceInfo Global.ServiceInformation) *ConfigurationService.ConfigurationService {
	return mw.Next.Init(staticConfigPath, serviceInfo)
}
