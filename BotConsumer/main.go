package main

import (
	"local/UbotTrade/API/ConfigurationServiceAPI"
	"local/UbotTrade/BotConsumer/Consumer"
	"local/UbotTrade/ConfigurationService/Service"
	"local/UbotTrade/DataHandlers/RabbitMQ"
	"local/UbotTrade/Global"
	"log"
	"fmt"
	asn1 "encoding/asn1"
	"local/UbotTrade/API/ServiceAPIFactory"
	"reflect"
	"golang.org/x/net/context"
	"local/UbotTrade/Logger"
	"time"
)

func main() { // add dynamic / static

	Handler := RabbitMQ.RabbitHandler{}
	log.Println(time.Now(),"Starting Bot consumer")
	log.Println(time.Now(),"Getting Configuraitions")

	Logger.SetGlobalLogger()
	ctx := context.Background()

	configService := ServiceAPIFactory.GetServiceInstance(ctx, reflect.TypeOf(ConfigurationServiceAPI.ConfigurationServiceAPI{})).(*ConfigurationServiceAPI.ConfigurationServiceAPI)
	ConfigResponse, err := configService.GetConfiguration(ctx, Global.ConfigurationRequest{Key: ConfigurationService.ConnectionStringsConfigurationKey})
	if err != nil {
		log.Println(err.Error())
	}
	connectionStringsConfig := ConfigResponse.Configuration["Data"].(map[string]interface{})
	Handler.Init(RabbitMQ.RabbitConfiguration{
		Topic:            "Bot_Requests",
		ConnectionString: connectionStringsConfig[ConfigurationService.RabbitMQConnectionStringKey].(string),
		Routing:          "trader.*",
		IsConsumer:true,
	})
	Handler.Consume("trader.*", BotConsumer.HandleMessage, HandleErrors, 10)
}


func HandleErrors(message []byte) bool {
	msgObj := RabbitMQ.RabbitMessage{}
	asn1.Unmarshal(message, &msgObj)
	fmt.Println(time.Now(),"Error begin -----------------------------------")
	fmt.Println(time.Now(),msgObj)
	fmt.Println(time.Now(),"Error end -------------------------------------")
	return true
}
