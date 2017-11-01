package BotConsumer

import (
	"context"
	"fmt"
	"local/UbotTrade/Global"
	"time"
	"github.com/gaillard/go-online-linear-regression/v1"
)

type TradeDecisionMaker interface {
	Make(enterPriceDifference float64, config Global.TradingConfig)
	ShouldBuy(ctx context.Context,p float64, high float64, lastPrices []float64) bool
	ShouldSell(ctx context.Context, p float64, buyPrice float64, config Global.TradingConfig, last20prices []float64) (bool, float64)
}

type BasicTradeDecisionMaker struct {
	priceDifference float64
	tradingConfig   Global.TradingConfig
}

func (strategy *BasicTradeDecisionMaker) Make(enterPriceDifference float64, config Global.TradingConfig) {
	strategy.priceDifference = 1 - enterPriceDifference
	strategy.tradingConfig = config
}

func (strategy *BasicTradeDecisionMaker) ShouldBuy(ctx context.Context, p float64, high float64, lastSlopResults []float64) bool {
	if p == 0.0 {
		return false
	}
	positiveSlop := true
	if len(lastSlopResults) > 0 {
		slope := strategy.getLinearRegressionOnSlope(lastSlopResults)
		fmt.Println(time.Now(), "BotConsumer", strategy.tradingConfig.BotNumber, "slope:", slope)
		positiveSlop = slope > 0
	}
	fmt.Println(time.Now(), "BotConsumer", strategy.tradingConfig.BotNumber, "p:", p, "high:", high, "target:", (strategy.priceDifference * high), "price dif:", ((1 - p/high) * 100), '%')
	var targetLowPrice = strategy.priceDifference * high

	InsertTickerData(ctx, strategy.tradingConfig.BotNumber, strategy.tradingConfig.Currency, Global.BotTickerDataModel{
		P:        p,
		H:        high,
		Action:   "Buy",
		Currency: strategy.tradingConfig.Currency,
		Stairs:   []float64{targetLowPrice},
	})

	if p > targetLowPrice || !positiveSlop {
		fmt.Println(time.Now(), "BotConsumer", strategy.tradingConfig.BotNumber, "The current price gives low propability of a profit, retry in 1 minute.")
		return false
	}
	return true
}

func (strategy *BasicTradeDecisionMaker) ShouldSell(ctx context.Context, p float64, buyPrice float64, config Global.TradingConfig, lastSlopResults []float64) (shouldSell bool, priceAfterComission float64) {

	shouldSell = false
	priceAfterComission = 0.0

	if p < buyPrice*(1-config.Fallback) {
		buyingAfterCommission := buyPrice * (1 + config.Fallback + config.BaseCommission/100)
		return true, (p - buyingAfterCommission)
	}

	slope := strategy.getLinearRegressionOnSlope(lastSlopResults)

	stairsPrices := make([]float64, 0)
	for stairIndex := len(config.Stairs) - 1; stairIndex >= 0; stairIndex-- {
		buyingAfterCommission := buyPrice * (1 + config.Stairs[stairIndex].Ratio + config.BaseCommission/100)
		stairsPrices = append(stairsPrices, buyingAfterCommission)
		fmt.Println(time.Now(), "BotConsumer", strategy.tradingConfig.BotNumber, "stair and commission", (1 + config.Stairs[stairIndex].Ratio + config.BaseCommission/100))
		fmt.Println(time.Now(), "BotConsumer", strategy.tradingConfig.BotNumber, "p", p, ">=", "buyingAfterCommission", buyingAfterCommission)
		if p >= buyingAfterCommission && slope < 0 {
			fmt.Println(time.Now(), "BotConsumer", strategy.tradingConfig.BotNumber, "profitable")
			shouldSell = true
			priceAfterComission = (p - buyingAfterCommission)
		}
	}

	var high = 0.0
	InsertTickerData(ctx, strategy.tradingConfig.BotNumber, strategy.tradingConfig.Currency, Global.BotTickerDataModel{
		P:        p,
		H:        high,
		Action:   "Sell",
		Currency: strategy.tradingConfig.Currency,
		Stairs:   stairsPrices,
	})
	return
}

func (strategy *BasicTradeDecisionMaker) getSlope(current float64, Previous float64) float64 {
	x2 := 2.0
	x1 := 1.0
	y2 := current
	y1 := Previous
	return (y2 - y1) / (x2 - x1)
}

func (strategy *BasicTradeDecisionMaker) getLinearRegressionOnSlope(lastSlopResults []float64) float64 {
	if len(lastSlopResults) > 0 {
		r := regression.New(10)
		x := 0
		for y := len(lastSlopResults) - 1; y >= 0; y-- {
			r.Add(float64(x), lastSlopResults[y])
			x++
		}
		slope, _ := r.Calculate()
		return slope
	}
	return 0.0
}




