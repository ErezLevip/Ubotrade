package BotService

import (
	"encoding/json"
	"errors"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
	"local/UbotTrade/API/RegistryServiceAPI"
	"local/UbotTrade/DataHandlers/MongoDB"
	"local/UbotTrade/DataHandlers/RabbitMQ"
	"local/UbotTrade/Global"
	"log"
	"time"
	"local/UbotTrade/BotConsumer/Consumer"
	"local/UbotTrade/API/ServiceAPIFactory"
	"reflect"
	"context"
)

type IBotService interface {
	GetBotInformation(ctx context.Context, Id uuid.UUID, BotNumber int) (result Global.BotInformation, err error)
	GetLastActivities(ctx context.Context, Id uuid.UUID, BotNumber int, MaxResultCount int) ([]Global.ActivityModel, error)
	GetAllActiveBots(ctx context.Context, userId string) ([]Global.GeneralBotInfo, error)
	GetBotTickerData(ctx context.Context, id uuid.UUID, botNumber int, maxResultCount int) ([]Global.BotTickerDataModel, error)
	GetBotProfits(ctx context.Context, id uuid.UUID, botNumber int, days int) ([]Global.ActivityModel, error)
	CreateNewBot(ctx context.Context,config Global.TradingConfig,userId string) error
	Init(serviceInfo Global.ServiceInformation) *BotService
	GetServiceMetrics() (outputMetrics map[string]interface{}, err error)
}
type BotService struct {
}

func (svc BotService) Init(serviceInfo Global.ServiceInformation) (botService *BotService) {
	log.Println("Starting Bot Service")

	botService = &BotService{}

	ctx := context.Background()

	//register the new service
	registryService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(RegistryServiceAPI.RegistryServiceAPI{})).(*RegistryServiceAPI.RegistryServiceAPI)

	_, err := registryService.Register(ctx, Global.RegistryRequest{ServiceInformation: serviceInfo})
	if(err != nil){
		log.Fatal(err.Error())
	}
	return
}

func (svc BotService) GetBotInformation(ctx context.Context, id uuid.UUID, BotNumber int) (Global.BotInformation, error) {
	mongoHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)
	var bot map[string]interface{}
	var err error
	if (id != uuid.UUID{}) {
		bot, err = mongoHandler.FindFirst(BotConsumer.BotsCollection, bson.M{"Id": id})
	} else if BotNumber != 0 {
		bot, err = mongoHandler.FindFirst(BotConsumer.BotsCollection, bson.M{"BotNumber": BotNumber})
	} else {
		err := errors.New("Id or botnumber is required")
		return Global.BotInformation{}, err
	}

	if err != nil {
		log.Println(err.Error())
	}
	if bot != nil && bot["BotData"] != nil {

		var botData Global.BotInformation
		botDataBson, err := bson.Marshal(bot["BotData"].(bson.M))
		err = bson.Unmarshal(botDataBson, &botData)
		return botData, err
	}
	return Global.BotInformation{}, err
}

func (svc BotService) GetAllActiveBots(ctx context.Context, userId string) (result []Global.GeneralBotInfo, err error) {
	mongoHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)

	result = make([]Global.GeneralBotInfo, 0)
	var bots []map[string]interface{}
	bots, err = mongoHandler.FindFMany(BotConsumer.BotsCollection, bson.M{})
	if err != nil {
		return
	}

	for _, bot := range bots {

		healthCheck := bot["LastHealthCheck"].(time.Time)
		if healthCheck.Add(2 * time.Minute).After(time.Now()) && bot != nil && bot["BotData"] != nil {
			var botDataBson []byte
			botDataBson, err = bson.Marshal(bot["BotData"].(bson.M))
			if err != nil {
				return
			}
			var botData Global.BotInformation
			err = bson.Unmarshal(botDataBson, &botData)
			if err != nil {
				return
			} else {

				if (botData.UserId == userId) {
					generalBotModel := Global.GeneralBotInfo{
						Id:        botData.Id,
						BotName:   botData.Name,
						BotNumber: botData.Configuration.BotNumber,
					}
					result = append(result, generalBotModel)
				}
			}
		}
	}
	log.Println(result)
	return
}

func (svc BotService) GetLastActivities(ctx context.Context, id uuid.UUID, botNumber int, maxResultCount int) ([]Global.ActivityModel, error) {
	mongoHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)

	var result = make([]Global.ActivityModel, 0)
	var activities = make([]map[string]interface{}, 0)
	var err error

	if (id != uuid.UUID{}) {
		activities, err = mongoHandler.FindFMany(BotConsumer.ActivitiesCollection, bson.M{"Id": id})
	} else if botNumber != 0 {
		activities, err = mongoHandler.FindFMany(BotConsumer.ActivitiesCollection, bson.M{"BotNumber": botNumber})
	} else {
		err = errors.New("Id cant be nill")
	}

	if err != nil {
		return result, err
	}
	if len(activities) > 0 {
		maxIndex := len(activities) - 1
		if maxResultCount == 0 {
			maxResultCount = 20
		} else if maxResultCount > maxIndex {
			maxResultCount = maxIndex
		}

		result = make([]Global.ActivityModel, maxResultCount)
		for i := maxIndex; i >= (maxIndex - maxResultCount); i-- {
			activity := activities[i]

			var activityData Global.ActivityModel
			activityBson, err := bson.Marshal(activity["ActivityData"].(bson.M))
			err = bson.Unmarshal(activityBson, &activityData)
			if err != nil {
				return result, err
			}
			result = append(result, activityData)
		}
		return result, err
	}
	return result, err
}

func (svc BotService) GetBotProfits(ctx context.Context, id uuid.UUID, botNumber int, days int) ([]Global.ActivityModel, error) {
	var result = make([]Global.ActivityModel, 0)
	mongoHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)

	var activities = make([]map[string]interface{}, 0)
	var err error

	if (id != uuid.UUID{}) {
		activities, err = mongoHandler.FindFMany(BotConsumer.ActivitiesCollection, bson.M{"Id": id})
	} else if botNumber != 0 {
		activities, err = mongoHandler.FindFMany(BotConsumer.ActivitiesCollection, bson.M{"BotNumber": botNumber})
	} else {
		err = errors.New("Id cant be nill")
	}

	if err != nil {
		return result, err
	}
	activitiesCount := len(activities)
	if activitiesCount > 0 {
		activitiesPerDay := map[int][] Global.ActivityModel{}
		currentDay := 1
		for i := activitiesCount - 1; i >= 0; i-- {
			activity := activities[i]

			var activityData Global.ActivityModel
			activityBson, err := bson.Marshal(activity["ActivityData"].(bson.M))
			err = bson.Unmarshal(activityBson, &activityData)
			if err != nil {
				return result, err
			}
			if activityData.TimeStamp.Day() == time.Now().Day() - currentDay && activityData.ActivityType == "Sell" {
				activitiesPerDay[currentDay] = append(activitiesPerDay[currentDay], activityData)
			} else if (activityData.TimeStamp.Day() < time.Now().Day() - currentDay) {
				currentDay--
			}
		}

		for _, v := range activitiesPerDay {
			result = append(result, v[0])
		}
	}
	return result, err
}
func (svc BotService) GetBotTickerData(ctx context.Context, id uuid.UUID, botNumber int, maxResultCount int) (tickerDataResult []Global.BotTickerDataModel, err error) {
	if maxResultCount == 0 {
		maxResultCount = 20
	}

	mongoHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)

	var tickerDatas []map[string]interface{}
	if (id != uuid.UUID{}) {
		tickerDatas, err = mongoHandler.FindFMany(BotConsumer.TickerDataCollection, bson.M{"Id": id})
	} else if botNumber != 0 {
		tickerDatas, err = mongoHandler.FindFMany(BotConsumer.TickerDataCollection, bson.M{"BotNumber": botNumber})
	}
	if err != nil {
		log.Println(err.Error())
	}
	if len(tickerDatas) > 0 {
		maxIndex := len(tickerDatas) - 1
		if maxResultCount == 0 {
			maxResultCount = 20
		} else if maxResultCount > maxIndex {
			maxResultCount = maxIndex
		}
		tickerDataResult = make([]Global.BotTickerDataModel, 0)

		for i := maxIndex; i >= (maxIndex - maxResultCount); i-- {
			var tickerData = tickerDatas[i]
			if tickerData["Data"] != nil {
				var tickerDataModel Global.BotTickerDataModel
				var modelString []byte
				modelString, err = json.Marshal(tickerData["Data"])
				json.Unmarshal(modelString, &tickerDataModel)
				tickerDataResult = append(tickerDataResult, tickerDataModel)
			}
		}
	}
	return
}

func (svc BotService) CreateNewBot(ctx context.Context, tradingConfiguration Global.TradingConfig, userId string) (err error) {
	rabbitHandler := ctx.Value("RabbitHandler").(*RabbitMQ.RabbitHandler)
	mongoHandler := ctx.Value("MongoHandler").(*MongoDB.MongoHandler)

	producer := BotProducer{
		DbHandler:     mongoHandler,
		RabbitHandler: rabbitHandler,
	}

	err = producer.Produce(tradingConfiguration, userId)
	if err != nil {
		log.Println(err.Error())
	}

	return
}

func (svc BotService) GetServiceMetrics() (outputMetrics map[string]interface{}, err error) {
	return make(map[string]interface{}), nil
}
