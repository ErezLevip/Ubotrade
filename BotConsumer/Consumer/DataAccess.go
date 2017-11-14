package BotConsumer

import (
	"github.com/erezlevip/Ubotrade/API/ServiceAPIFactory"
	"github.com/erezlevip/Ubotrade/DataHandlers/MongoDB"
	"github.com/erezlevip/Ubotrade/API/UserServiceAPI"
	"github.com/erezlevip/Ubotrade/Global"
	"github.com/erezlevip/Ubotrade/UserService/Service"
	"github.com/nu7hatch/gouuid"

	"fmt"
	"reflect"
	"time"
	"encoding/json"
	"context"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func UpdateBotLastCheck(ctx context.Context) {
	dbHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)
	botId := ctx.Value("BotId").(uuid.UUID)
	err := dbHandler.Update(BotsCollection, bson.M{"Id": botId}, bson.M{"LastHealthCheck": time.Now()})
	if err != nil {
		log.Println(err.Error())
	}
}

func GetActivityById(ctx context.Context, activityId uuid.UUID) (activityData Global.ActivityModel) {
	dbHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)
	res, err := dbHandler.FindFirst(ActivitiesCollection, bson.M{"Id": activityId})
	err = json.Unmarshal([]byte(res["ActivityData"].(string)), &activityData)

	if err != nil {
		log.Println(err.Error())
	}
	return
}

func InsertBot(ctx context.Context, information Global.BotInformation) {
	dbHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)
	botNumber := ctx.Value("BotNumber").(int)
	botId := ctx.Value("BotId").(uuid.UUID)
	userId := ctx.Value("UserId").(string)
	err := dbHandler.Insert(BotsCollection, bson.M{"TimeStamp": time.Now(),
		"Id": botId,
		"BotNumber": botNumber,
		"IsActive": true,
		"LastHealthCheck": time.Now(),
		"BotData": Global.BotInformation{
			Configuration:        information.Configuration,
			Amount:               information.Amount,
			Id:                   botId,
			LastActivityId:       information.LastActivityId,
			Name:                 information.Name,
			UserId:               userId,
			CurrencyAmount:       information.CurrencyAmount,
			LiquidAmountCurrency: information.LiquidAmountCurrency,
			LiquidAmountUSD:      information.LiquidAmountUSD,
			OriginalAmount:       information.OriginalAmount,
		},
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func InsertActivity(ctx context.Context, activity Global.ActivityModel) {
	dbHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)
	botNumber := ctx.Value("BotNumber").(int)
	botId := ctx.Value("BotId").(uuid.UUID)
	activityId, _ := uuid.NewV4()
	err := dbHandler.Insert(ActivitiesCollection, bson.M{"Id": botId, "BotNumber": botNumber,
		"ActivityData": Global.ActivityModel{
			ActivityPrice:        activity.ActivityPrice,
			ActivityType:         activity.ActivityType,
			ActualAmountCurrency: activity.ActualAmountCurrency,
			ActualAmountUSD:      activity.ActualAmountUSD,
			BotId:                botId,
			Id:                   *activityId,
			PriceDifference:      activity.PriceDifference,
			TimeStamp:            time.Now(),
		}})
	if err != nil {
		log.Println(err.Error())
	}
}

func GetLastTickerData(ctx context.Context, priceBlockForSlop int) []float64 {
	dbHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)
	currency := ctx.Value("Currency").(string)
	all, err := dbHandler.FindFMany(TickerDataCollection, bson.M{"Currency": currency})
	pricesCount := len(all)
	if err != nil {
		fmt.Println(err.Error())
	}
	res := make([]float64, 0)
	if pricesCount > 0 {
		bulkSize := priceBlockForSlop

		if pricesCount > bulkSize {
			pricesCount = bulkSize
		}

		index := len(all) - 1
		for i := 0; i < pricesCount; i++ {
			val := all[index]["Data"].(bson.M)
			res = append(res, (val["p"].(float64)))
			index--
		}
	}
	return res
}

func InsertNewNotification(ctx context.Context, message string) {
	userId := ctx.Value("UserId").(string)
	botNumber := ctx.Value("BotNumber").(int)
	userSvc := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(UserServiceAPI.UserServiceAPI{})).(*UserServiceAPI.UserServiceAPI)
	var notification = make(map[string]interface{})
	notification["NotificationMessage"] = message
	notification["BotNumber"] = botNumber
	notification["IsActive"] = true
	notification["UserId"] = userId

	userSvc.SetUser(Global.SetUserRequest{
		DataType:  UserService.NotificationsDataType,
		UserId:    userId,
		Data:      notification,
		Operation: UserService.CreateOperation,
	})
}

func InsertTickerData(ctx context.Context, botNumber int, currency string, tickerData Global.BotTickerDataModel) {
	dbHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)
	err := dbHandler.Insert(TickerDataCollection,
		bson.M{"TimeStamp": time.Now(),
			"BotNumber": botNumber,
			"Currency": currency,
			"Data": tickerData})
	if err != nil {
		log.Println(err.Error())
	}
}

func GetBotNumber(ctx context.Context) int {
	dbHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)
	allBots, err := dbHandler.FindFMany(BotsCollection, nil)
	if err != nil {
		log.Println(err.Error())
	}
	highest := 0
	if len(allBots) > 0 {
		for _, bot := range allBots {
			botVal := bot["BotData"].(bson.M)["configuration"].(bson.M)["botnumber"].(int)
			if botVal > highest {
				highest = botVal
			}
		}
		highest++
	}

	return highest + 1
}
