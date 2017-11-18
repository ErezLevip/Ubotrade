package AuthenticationInstrumenting

import (
	"fmt"
	"time"
	"context"

	"github.com/go-kit/kit/metrics"
	"github.com/erezlevip/Ubotrade/Global"
	"github.com/erezlevip/Ubotrade/AuthenticationService/Service"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
	Next           AuthenticationService.IAuthenticationService
}

var FailureCount int
var RequestCount int

func Make(requestCount metrics.Counter, requestLatency metrics.Histogram, countResult metrics.Histogram, next AuthenticationService.IAuthenticationService) (instance *InstrumentingMiddleware) {
	FailureCount = 0
	RequestCount = 0

	return &InstrumentingMiddleware{
		RequestCount:   requestCount,
		CountResult:    countResult,
		RequestLatency: requestLatency,
		Next:           next,
	}

}

func (mw InstrumentingMiddleware) Login(ctx context.Context, clientId string, firstName string, lastName string, email string, sessionId string) (userId string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Login", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil && err.Error() != "UnAuthorized" {
			FailureCount++
		}
	}(time.Now())
	userId, err = mw.Next.Login(ctx, clientId, firstName, lastName, email, sessionId)
	return
}
func (mw InstrumentingMiddleware) GetToken(ctx context.Context, clientId string) (token string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetToken", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	token, err = mw.Next.GetToken(ctx, clientId)
	return
}
func (mw InstrumentingMiddleware) ValidateToken(ctx context.Context, token string) (IsValid bool, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ValidateToken", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	IsValid, err = mw.Next.ValidateToken(ctx, token)
	return
}

func (mw InstrumentingMiddleware) Init(serviceInfo Global.ServiceInformation) *AuthenticationService.AuthenticationService {
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
