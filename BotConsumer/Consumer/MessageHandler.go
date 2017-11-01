package BotConsumer

import (
	"context"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"local/UbotTrade/API/ConfigurationServiceAPI"
	"local/UbotTrade/ConfigurationService/Service"
	"local/UbotTrade/DataHandlers/MongoDB"
	"local/UbotTrade/DataHandlers/RabbitMQ"
	"local/UbotTrade/Global"
	"log"
	"time"
	"local/UbotTrade/API/ServiceAPIFactory"
	"reflect"
	"strconv"
)

// HandleMessage Triggerd by the consumer and invoke the bot from the message
func HandleMessage(message []byte) bool {
	msgObj := RabbitMQ.RabbitMessage{}
	asn1.Unmarshal(message, &msgObj)
	fmt.Println(time.Now(),"message begin -----------------------------------")
	fmt.Println(time.Now(),msgObj)
	fmt.Println(time.Now(),"message end -------------------------------------")

	botData := Global.BotInformation{}
	err := json.Unmarshal([]byte(msgObj.Payload), &botData)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	botData.UserId = msgObj.SenderId
	bot := BasicUbotTrader{}
	strategy := &BasicTradeDecisionMaker{}
	api := &BasicTradeApi{}

	isMonitor := msgObj.SenderId == "monitor"
	ctx := createBotContext(botData, isMonitor)

	if !isMonitor {
		botData.Configuration.BotNumber = ctx.Value("BotNumber").(int)
	}

	retryTimeDuration, err := strconv.ParseInt(botData.Configuration.RetryTimeDuration,10,64)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	api.Make(botData.Configuration.Currency,time.Duration(retryTimeDuration))
	// initialize bot strategy
	strategy.Make(botData.Configuration.PercentagePriceDifferenceOnBuy, botData.Configuration)

	//create the bot.
	bot.Make(ctx, strategy, api, botData)
	currentAmount := 100.0
	go func() {
		for {
			liquidAmountUSD, revenue := bot.Start(ctx, currentAmount, msgObj.Monitor)
			if liquidAmountUSD != 0.0 && revenue != 0.0 {
				fmt.Println(time.Now(), "BotConsumer", botData.Configuration.BotNumber, "liquidAmountUSD :", liquidAmountUSD, "$", "revenue: ", revenue, "$")
			} else {
				fmt.Println(time.Now(), "BotConsumer", botData.Configuration.BotNumber, "bot response is empty")
			}
		}
	}()
	return true
}

func createBotContext(botData Global.BotInformation, monitor bool) (ctx context.Context) {
	ctx = context.Background()
	dbHandler := &MongoDB.MongoHandler{}

	configService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(ConfigurationServiceAPI.ConfigurationServiceAPI{})).(*ConfigurationServiceAPI.ConfigurationServiceAPI)
	ConfigResponse, err := configService.GetConfiguration(ctx, Global.ConfigurationRequest{Key: ConfigurationService.ConnectionStringsConfigurationKey})
	if err != nil {
		log.Println(err.Error())
	}
	connectionStringsConfig := ConfigResponse.Configuration["Data"].(map[string]interface{})
	dbHandler.Init(map[string]string{
		MongoDB.ConnectionStringKey: connectionStringsConfig[ConfigurationService.MongoConnectionStringKey].(string),
		MongoDB.DatabaseKey:         connectionStringsConfig[ConfigurationService.MongoDatabaseKey].(string),
	})

	ctx = context.WithValue(ctx, "MongoHandler", dbHandler)

	var botNumber int
	if (monitor) {
		botNumber = botData.Configuration.BotNumber
	} else {
		botNumber = GetBotNumber(ctx)
	}

	ctx = context.WithValue(ctx, "BotNumber", botNumber)
	ctx = context.WithValue(ctx, "BotId", botData.Id)
	ctx = context.WithValue(ctx, "Currency", botData.Configuration.Currency)
	ctx = context.WithValue(ctx, "UserId", botData.UserId)
	return
}
