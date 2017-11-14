package BotConsumer

import (
	"context"
	"strconv"
	"fmt"
	"time"
	"log"

	"github.com/erezlevip/Ubotrade/Global"
	"github.com/nu7hatch/gouuid"
)

type UbotTrader interface {
	Start(config Global.TradingConfig)
}

type BasicUbotTrader struct {
	ctx                  context.Context
	strategy             TradeDecisionMaker
	tradingApi           TradeApi
	config               Global.TradingConfig
	amount               float64
	p                    float64
	buyPrice             float64
	originalAmount       float64
	currencyAmount       float64
	liquidAmountCurrency float64
	liquidAmountUSD      float64
	lastActivity         Global.ActivityModel
	Id                   uuid.UUID
	Name                 string
	UserId               string
}

const ActivitiesCollection = "Activities"
const TickerDataCollection = "TickerData"
const BotsCollection = "FreshBots"

const BuyActivityType = "Buy"
const SellActivityType = "Sell"

func (bot *BasicUbotTrader) Make(ctx context.Context, strategy TradeDecisionMaker, tradingApi TradeApi, botData Global.BotInformation) {
	log.Println(time.Now(),"------------------------------------------starting bot ", botData.Name)
	bot.strategy = strategy
	bot.tradingApi = tradingApi
	bot.config = botData.Configuration
	bot.currencyAmount = botData.CurrencyAmount
	bot.liquidAmountCurrency = botData.LiquidAmountCurrency
	bot.liquidAmountUSD = botData.LiquidAmountUSD
	bot.originalAmount = botData.OriginalAmount
	bot.Id = botData.Id
	bot.Name = botData.Name
	bot.UserId = botData.UserId

	if (botData.LastActivityId != uuid.UUID{}) {
		bot.lastActivity = GetActivityById(ctx, botData.LastActivityId)
	}
}
func (bot *BasicUbotTrader) Start(ctx context.Context, usdAmount float64, monitor bool) (float64, float64) {
	bot.amount = usdAmount

	if !monitor {
		InsertBot(ctx, Global.BotInformation{
			Amount:               bot.amount,
			Configuration:        bot.config,
			CurrencyAmount:       bot.currencyAmount,
			LiquidAmountCurrency: bot.liquidAmountCurrency,
			LiquidAmountUSD:      bot.liquidAmountUSD,
			OriginalAmount:       bot.originalAmount,
			Name:                 bot.Name,
		})
		InsertNewNotification(ctx, "Bot "+strconv.Itoa(bot.config.BotNumber)+" Has started trading: "+bot.config.Currency)
		fmt.Println(time.Now(), "BotConsumer", bot.config.BotNumber, "starting bot using config:", bot.config)
	} else {
		InsertNewNotification(ctx, "Bot "+strconv.Itoa(bot.config.BotNumber)+" Re activated trading: "+bot.config.Currency)
		fmt.Println(time.Now(), "BotConsumer", bot.config.BotNumber, "reactivate bot using config:", bot.config)
	}
	high := 0.0

	if !monitor || bot.buyPrice == 0.0 && bot.lastActivity.ActivityType == "Sell" {
		for {
			fmt.Println(time.Now(), "BotConsumer", bot.config.BotNumber, "Current USD Balance", bot.liquidAmountUSD)
			bot.p, high = bot.tradingApi.GetTickerInfo()
			UpdateBotLastCheck(ctx)
			lastSlopPrices := make([]float64, 0)
			if bot.config.BuyOnPositiveSlop {
				lastSlopPrices = GetLastTickerData(ctx, bot.config.PriceBlockForSlop)
			}

			resp := bot.strategy.ShouldBuy(ctx, bot.p, high, lastSlopPrices)
			if !resp {
				time.Sleep(time.Duration(1) * time.Minute)
			} else {
				bot.buyPrice = bot.p
				break
			}
		}

		currencyAmount := float64(bot.amount) / bot.buyPrice
		activityResponse := bot.tradingApi.Buy(currencyAmount, bot.p)
		if activityResponse {
			bot.buyPrice = bot.p
			bot.currencyAmount = currencyAmount
			fmt.Println(time.Now(), "BotConsumer", bot.config.BotNumber, "you bought", bot.currencyAmount, bot.config.Currency, "for", bot.amount, "$", "original buying price:", bot.buyPrice)
			bot.liquidAmountUSD -= bot.amount
			InsertNewNotification(ctx, "Bot "+strconv.Itoa(bot.config.BotNumber)+" Bought "+fmt.Sprintf("%.2f", bot.currencyAmount)+" "+bot.config.Currency)
		}
		InsertActivity(ctx, Global.ActivityModel{
			ActivityPrice:        bot.p,
			ActivityType:         BuyActivityType,
			ActualAmountCurrency: bot.currencyAmount,
			ActualAmountUSD:      bot.liquidAmountUSD,
			BotId:                bot.Id,
			PriceDifference:      0,
			TimeStamp:            time.Now(),
		})

		bot.originalAmount = bot.currencyAmount
	}
	var lastPrice float64 = 0
	for {
		fmt.Println(time.Now(), "BotConsumer", bot.config.BotNumber, "Current USD Balance", bot.liquidAmountUSD)
		bot.p, high = bot.tradingApi.GetTickerInfo()
		UpdateBotLastCheck(ctx)
		if lastPrice <= 0.0 {
			fmt.Println(time.Now(), "BotConsumer", bot.config.BotNumber, "high", high, "current price", bot.p)
		} else if bot.p > lastPrice {
			fmt.Println(time.Now(), "BotConsumer", bot.config.BotNumber, "high", high, "current price", bot.p, "--------UP")
		} else {
			fmt.Println(time.Now(), "BotConsumer", bot.config.BotNumber, "high", high, "current price", bot.p, "--------DOWN")
		}
		lastPrice = bot.p
		if bot.p > 0 {

			sell, priceDif := bot.strategy.ShouldSell(ctx, bot.p, bot.buyPrice, bot.config, GetLastTickerData(ctx, bot.config.PriceBlockForSlop))
			if sell {
				actualAmountUSD := bot.tradingApi.Bid(bot.currencyAmount, bot.p)
				bot.liquidAmountUSD += actualAmountUSD
				bot.liquidAmountCurrency += actualAmountUSD / bot.p
				InsertActivity(ctx, Global.ActivityModel{
					ActivityPrice:        bot.p,
					ActivityType:         SellActivityType,
					ActualAmountCurrency: bot.currencyAmount,
					ActualAmountUSD:      bot.liquidAmountUSD,
					BotId:                bot.Id,
					PriceDifference:      priceDif,
					TimeStamp:            time.Now(),
				})

				if priceDif < 0 {
					InsertNewNotification(ctx, "Bot "+strconv.Itoa(bot.config.BotNumber)+" Lost: "+strconv.FormatFloat(priceDif, 'E', -1, 64))
					fmt.Println(time.Now(), "BotConsumer", bot.config.BotNumber, "Sold in loss: ", priceDif)
				} else {
					InsertNewNotification(ctx, "Bot "+strconv.Itoa(bot.config.BotNumber)+" Sold: "+strconv.FormatFloat(actualAmountUSD, 'E', -1, 64))
					fmt.Println(time.Now(), "BotConsumer", bot.config.BotNumber, "sold actual:", actualAmountUSD, "liquid Currency:", bot.liquidAmountCurrency, "liquidUSD", bot.liquidAmountUSD, "dif", priceDif)
				}
				return bot.liquidAmountUSD, (bot.liquidAmountUSD - float64(bot.amount))
			}
		}
		iterationTime,err := strconv.ParseInt(bot.config.TimeIterations,10,64)
		if(err != nil){
			log.Panic(err.Error())
		}

		time.Sleep(time.Duration(iterationTime))
	}
}
