package UserInstrumenting

import (
	"fmt"
	"time"
	"context"

	"github.com/go-kit/kit/metrics"
	"github.com/erezlevip/Ubotrade/Global"
	"github.com/erezlevip/Ubotrade/UserService/Service"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
	Next           UserService.IUserService
}

var FailureCount int
var RequestCount int

func Make(requestCount metrics.Counter, requestLatency metrics.Histogram, countResult metrics.Histogram, next UserService.IUserService) (instance *InstrumentingMiddleware) {
	FailureCount = 0
	RequestCount = 0

	return &InstrumentingMiddleware{
		RequestCount:   requestCount,
		CountResult:    countResult,
		RequestLatency: requestLatency,
		Next:           next,
	}

}

func (mw InstrumentingMiddleware) GetUser(ctx context.Context, userId string, dataType string, activeOnly bool) (data []map[string]interface{}, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetUser", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	data, err = mw.Next.GetUser(ctx, userId, dataType, activeOnly)
	return
}
func (mw InstrumentingMiddleware) CreateUser(ctx context.Context, userId string, firstName string, lastName string, email string) (data map[string]interface{}, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateUser", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	data, err = mw.Next.CreateUser(ctx, userId, firstName, lastName, email)
	return
}
func (mw InstrumentingMiddleware) SetUser(ctx context.Context, userId string, dataType string, operation string, data map[string]interface{}) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SetUser", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	err = mw.Next.SetUser(ctx, userId, dataType, operation, data)
	return
}

func (mw InstrumentingMiddleware) Init(serviceInfo Global.ServiceInformation) *UserService.UserService {
	return mw.Next.Init(serviceInfo)
}

func (mw InstrumentingMiddleware) GetServiceMetrics() (metrics map[string]interface{}, err error) {
	metrics = map[string]interface{}{
		"SuccessRate": getSuccessRate(),
	}
	err = nil
	return
}

func getSuccessRate() float64 {
	//will always return 100 for now
	if true || RequestCount == 0 || FailureCount == 0 {
		return 100.0
	}

	return float64((1 - (FailureCount / RequestCount)) * 100)
}
