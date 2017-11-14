package BotLogging

import (
	"time"
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/nu7hatch/gouuid"
	"github.com/erezlevip/Ubotrade/BotService/Service"
	"github.com/erezlevip/Ubotrade/Global"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   BotService.IBotService
}

func (mw LoggingMiddleware) GetBotInformation(ctx context.Context,  Id uuid.UUID, BotNumber int) (bots Global.BotInformation, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetBotInformation",
			"input", fmt.Sprintf("%b", Id, BotNumber),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	bots, err = mw.Next.GetBotInformation(ctx, Id, BotNumber)
	return
}

func (mw LoggingMiddleware) GetLastActivities(ctx context.Context,  Id uuid.UUID, BotNumber int, MaxResultCount int) (activities []Global.ActivityModel, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetBotActivities",
			"input", fmt.Sprintf("%b", Id, BotNumber, MaxResultCount),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	activities, err = mw.Next.GetLastActivities(ctx, Id, BotNumber, MaxResultCount)
	return
}

func (mw LoggingMiddleware) GetAllActiveBots(ctx context.Context,  userId string) (bots []Global.GeneralBotInfo, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetAllActiveBots",
			"input", fmt.Sprintf("%b",userId),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	bots, err = mw.Next.GetAllActiveBots(ctx, userId)
	return
}

func (mw LoggingMiddleware) GetBotTickerData(ctx context.Context, id uuid.UUID, botNumber int, maxResultCount int) (tickerData []Global.BotTickerDataModel, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetBotTickerData",
			"input", fmt.Sprintf("%b", maxResultCount, id, botNumber),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	tickerData, err = mw.Next.GetBotTickerData(ctx, id, botNumber, maxResultCount)
	return
}
func (mw LoggingMiddleware) GetBotProfits(ctx context.Context, id uuid.UUID, botNumber int, days int) (profits []Global.ActivityModel, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetBotProfits",
			"input", fmt.Sprintf("%b", days, id, botNumber),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	profits, err = mw.Next.GetBotProfits(ctx, id, botNumber, days)
	return
}

func (mw LoggingMiddleware) CreateNewBot(ctx context.Context, tradingConfiguration Global.TradingConfig, UserId string) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "CreateNewBot",
			"input", fmt.Sprintf("%b", tradingConfiguration),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.Next.CreateNewBot(ctx, tradingConfiguration, UserId)
	return
}

func (mw LoggingMiddleware) Init(serviceInfo Global.ServiceInformation) *BotService.BotService {
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
