package main

import (
	"log"
	"fmt"
	asn1 "encoding/asn1"
	"reflect"
	"golang.org/x/net/context"
	"time"

	"github.com/erezlevip/Ubotrade/API/ConfigurationServiceAPI"
	"github.com/erezlevip/Ubotrade/BotConsumer/Consumer"
	"github.com/erezlevip/Ubotrade/ConfigurationService/Service"
	"github.com/erezlevip/Ubotrade/DataHandlers/RabbitMQ"
	"github.com/erezlevip/Ubotrade/Global"
	"github.com/erezlevip/Ubotrade/Logger"
	"github.com/erezlevip/Ubotrade/API/ServiceAPIFactory"
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
