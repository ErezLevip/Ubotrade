package BotInstrumenting

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/nu7hatch/gouuid"
	"local/UbotTrade/BotService/Service"
	"local/UbotTrade/Global"
	"context"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
	Next           BotService.IBotService
}

var FailureCount int
var RequestCount int

func Make(requestCount metrics.Counter, requestLatency metrics.Histogram, countResult metrics.Histogram, next BotService.IBotService) (instance *InstrumentingMiddleware) {
	FailureCount = 0
	RequestCount = 0

	return &InstrumentingMiddleware{
		RequestCount:   requestCount,
		CountResult:    countResult,
		RequestLatency: requestLatency,
		Next:           next,
	}

}

func (mw InstrumentingMiddleware) GetBotInformation(ctx context.Context, Id uuid.UUID, BotNumber int) (bots Global.BotInformation, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetBotInformation", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	bots, err = mw.Next.GetBotInformation(ctx, Id, BotNumber)
	return
}
func (mw InstrumentingMiddleware) GetLastActivities(ctx context.Context, Id uuid.UUID, BotNumber int, MaxResultCount int) (activities []Global.ActivityModel, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetLastActivities", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	activities, err = mw.Next.GetLastActivities(ctx, Id, BotNumber, MaxResultCount)
	return
}

func (mw InstrumentingMiddleware) GetAllActiveBots(ctx context.Context, userId string) (bots []Global.GeneralBotInfo, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetAllActiveBots", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	bots, err = mw.Next.GetAllActiveBots(ctx, userId)
	return
}

func (mw InstrumentingMiddleware) GetBotTickerData(ctx context.Context,  id uuid.UUID, botNumber int, maxResultCount int) (prices []Global.BotTickerDataModel, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetBotTickerData", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	prices, err = mw.Next.GetBotTickerData(ctx , id, botNumber, maxResultCount)
	return
}
func (mw InstrumentingMiddleware) GetBotProfits(ctx context.Context,  id uuid.UUID, botNumber int, days int) (prices []Global.ActivityModel, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetBotProfits", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	prices, err = mw.Next.GetBotProfits(ctx, id, botNumber, days)
	return
}

func (mw InstrumentingMiddleware) CreateNewBot(ctx context.Context, tradingConfiguration Global.TradingConfig, UserId string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateNewBot", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		RequestCount++
		if err != nil {
			FailureCount++
		}
	}(time.Now())
	err = mw.Next.CreateNewBot(ctx, tradingConfiguration, UserId)
	return
}

func (mw InstrumentingMiddleware) Init(serviceInfo Global.ServiceInformation) *BotService.BotService {
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
