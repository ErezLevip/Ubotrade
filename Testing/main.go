package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"local/UbotTrade/API/BotServiceAPI"
	"local/UbotTrade/API/ConfigurationServiceAPI"
	"local/UbotTrade/API/RegistryServiceAPI"
	"local/UbotTrade/ConfigurationService/Service"
	"local/UbotTrade/DataHandlers/MongoDB"
	"local/UbotTrade/DataHandlers/Redis"
	"local/UbotTrade/Global"
	"log"
	"time"
	"local/UbotTrade/API/ServiceAPIFactory"
	"reflect"
	"local/UbotTrade/BotConsumer/Consumer"
	"golang.org/x/net/context"
	"strconv"
)

func main() {

	//ConfigUpdate()
	//ProducerTrigger(Implementations.BitcoinCurrency)
	//time.Sleep(5 * time.Second)
	ProducerTrigger(BotConsumer.EtherCurrency)
	time.Sleep(time.Duration(20 * time.Second))
	ProducerTrigger(BotConsumer.BitcoinCurrency)
	/*TestBOtService()
	TestBOtActivitiesService()
	TestBOtPricesService()*/
	//getservices("Services."+ConfigurationServiceAPI.ConfigurationServiceKey)
	//GetServiceTest(BotServiceAPI.BotServiceKey)
	//TestBOtService()
	/*api :=BaseServiceAPI.BaseServiceAPI{}
	base :=api.Make("UbotTrade/Testing/RegistryService.json","Configuration",&Global.ServiceInformation{
		Url:"http://localhost",
		Port:8081,
		ServiceName:"Configuration",
	})
	res1,err1 :=base.GetRequestMetrics()
	if(err1 != nil){
		log.Println(err1.Error())
	}*/

	//	getservices()
	//GetConfigurations()
}

func ConfigUpdate() {
	/*
		  "mongodbconnectionstring":"localhost",
	  "mongodatabase":"DEV",
	  "Environment":"DEV",
	  "DefaultPort":8083,
	  "RegistryService":"http://localhost:8080",
	  "redisconnectionstring":"localhost:6379",
	  "redisdb":1,
	  "rediscredentials":""
	}


	*/
	handler := MongoDB.MongoHandler{}
	config := make(map[string]string)
	config[MongoDB.ConnectionStringKey] = "localhost"
	config[MongoDB.DatabaseKey] = "DEV"

	handler.Init(config)
	err := handler.Insert(ConfigurationService.ConfigurationsCollection, bson.M{"Key": ConfigurationService.ConnectionStringsConfigurationKey, "Data": ConnectionStrings{
		RabbitMQConnectionString: "amqp://guest:guest@localhost:5672/",
		MongoDbConnectionString:  "localhost",
		RedisConnectionString:    "localhost:6379",
		MongoDatabase:            "DEV",
		RedisDb:                  1,
		RedisCredentials:         "",
	},
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func GetServiceTest(serviceKey string) {
	ctx := context.Background()
	registryService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(RegistryServiceAPI.RegistryServiceAPI{})).(*RegistryServiceAPI.RegistryServiceAPI)
	resp, err := registryService.GetService(ctx, Global.RegistryRequest{ServiceInformation: Global.ServiceInformation{ServiceName: serviceKey}})

	log.Println(time.Now(),"registry response", resp, resp.ServiceInformation)

	if err != nil {
		log.Println(err.Error())
	}
}

/*
func TestBOtService()  {
	botservice :=BotServiceAPI.Make("UbotTrade/Testing/RegistryService.json")
	s :=botservice.Base.ServiceInformation
	log.Println("s",s.Url,s.Port)
	res ,err :=botservice.GetBotInformation(Global.BotInformationRequest{
		Parameters: map[string]interface{}{"IsActive":true,
		},
	})
	if(err != nil){
		log.Println("err",err.Error())
	}
	log.Println("res",res)

}

func TestBOtPricesService()  {
	botservice :=BotServiceAPI.Make("UbotTrade/Testing/RegistryService.json")
	s :=botservice.Base.ServiceInformation
	log.Println("s",s.Url,s.Port)
	res ,err :=botservice.GetPricesData(Global.BotInformationRequest{
		Parameters: map[string]interface{}{},
	})
	if(err != nil){
		log.Println("err",err.Error())
	}
	log.Println("res",res)

}
func TestBOtActivitiesService()  {
	botservice :=BotServiceAPI.Make("UbotTrade/Testing/RegistryService.json")
	s :=botservice.Base.ServiceInformation
	log.Println("s",s.Url,s.Port)
	res ,err :=botservice.GetBotActivities(Global.BotInformationRequest{
		Parameters: map[string]interface{}{},
	})
	if(err != nil){
		log.Println("err",err.Error())
	}
	log.Println("res",res)

}*/

/*
	con :=ConnectionStrings{
		RabbitMQConnectionString:"amqp://guest:guest@localhost:5672/",
		MongoDbConnectionString : "localhost",
		RedisConnectionString:"localhost:6379",
		MongoDatabase:"DEV",
		RedisDb:1,
		RedisCredentials:"",
	}
*/

func TestStaticConfig() {
	/*config := Global.ServiceConfiguration{Url:staticConfig["RegistryService"].(string)}
	fmt.Println("registry config",config)*/
}

func getservices(key string) {
	redis := DataHandlers.RedisHandler{}
	redis.Init(DataHandlers.RedisConfiguration{
		ConnectionString: "localhost:6379",
		Credentials:      "",
		Db:               1,
	})
	res, err := redis.Get(key)
	var response []Global.ServiceInformation
	err = json.Unmarshal([]byte(res), &response)

	for _, val := range response {
		log.Println(val)
	}

	log.Println(time.Now(), "now")
	if err != nil {
		log.Println(err.Error())
	}
}

func GetConfigurations() {
	ctx := context.Background()
	configService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(ConfigurationServiceAPI.ConfigurationServiceAPI{})).(*ConfigurationServiceAPI.ConfigurationServiceAPI)
	//	api := RegistryServiceAPI.RegistryServiceAPI{}
	//api.Make("UbotTrade/Testing/RegistryService.json")
	//api.Register(Global.RegistryRequest{ServiceInformation:Global.ServiceInformation{ServiceName:"test",Port:8081,Url:"localhost"}})
	res, err := configService.GetConfiguration(ctx, Global.ConfigurationRequest{Key: "ConnectionStrings"})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("res", res)
	}
}

func GetFromMongo(collection string, Q bson.M) {
	dbHandler := MongoDB.MongoHandler{}
	dbConfig := make(map[string]string)
	dbConfig[MongoDB.ConnectionStringKey] = "localhost"
	dbConfig[MongoDB.DatabaseKey] = "DEV"
	dbHandler.Init(dbConfig)
	//const configurationsCollection = "Configurations"
	res1, err := dbHandler.FindFMany(collection, Q)
	if err != nil {
		log.Println(err.Error())
	}
	for i, val := range res1 {
		log.Println(i, val)
	}
}

func ProducerTrigger(currency string) {
	original := 0.04
	priceBlockForSlop := 10
	buyOnPositiveSlop := true

	avg := original
	avgHigh := avg + 0.02
	avgLow := avg - 0.02
	//avgMax := avg+0.04
	fmt.Println(time.Now(), "starting point:", "avgLow:", avgLow, "avg:", avg, "avgHigh:", avgHigh, "currency", currency)

	config := Global.TradingConfig{
		BaseCommission:                 0.02,
		Fallback:                       0.05,
		LastStairTimeout:               strconv.FormatInt((time.Duration(10) * time.Minute).Nanoseconds(),10),
		TimeIterations:                 strconv.FormatInt((time.Duration(30) * time.Second).Nanoseconds(),10),
		Stairs:                         []Global.Stair{{Ratio: avgLow}, {Ratio: avg}, {Ratio: avgHigh}},
		Currency:                       currency,
		PriceBlockForSlop:              priceBlockForSlop,
		BuyOnPositiveSlop:              buyOnPositiveSlop,
		PercentagePriceDifferenceOnBuy: 0.02,
		RetryTimeDuration:              strconv.FormatInt(( 10 * time.Second).Nanoseconds(),10),
		BotNumber:                      0,
	}

	ctx := context.Background()
	botsvc := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(BotServiceAPI.BotServiceAPI{})).(*BotServiceAPI.BotServiceAPI)
	res, err := botsvc.CreateNewBot(Global.CreateBotRequest{
		TradingConfiguration: config,
		UserId:               "118103040085940455572",
	})
	log.Println(res)
	if err != nil {
		log.Println(err.Error())
	}
}

type ConnectionStrings struct {
	RedisConnectionString    string `json:"redis_connection_string"`
	RabbitMQConnectionString string `json:"rabbit_mq_connection_string"`
	MongoDbConnectionString  string `json:"mongo_db_connection_string"`
	RedisDb                  int    `json:"redis_db"`
	MongoDatabase            string `json:"mongo_database"`
	RedisCredentials         string `json:"redis_credentials"`
}
