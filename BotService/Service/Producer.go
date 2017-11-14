package BotService

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nu7hatch/gouuid"
	"github.com/erezlevip/Ubotrade/DataHandlers/MongoDB"
	"github.com/erezlevip/Ubotrade/DataHandlers/RabbitMQ"
	"github.com/erezlevip/Ubotrade/Global"
)

type BotProducer struct {
	DbHandler     MongoDB.IDbHandler
	RabbitHandler RabbitMQ.IRabbitHandler
}

func (producer *BotProducer) Produce(config Global.TradingConfig, userId string) (err error) {

	botId, _ := uuid.NewV4()
	log.Println(time.Now(),"New Bot Id", botId)

	botData := Global.BotInformation{
		Id:            *botId,
		Configuration: config,
		Name:          config.Currency + " Trader",
	}

	configJson, err := json.Marshal(botData)
	producer.RabbitHandler.Publish(RabbitMQ.RabbitMessage{SenderId: userId,
		SentTime: time.Now(),
		Topic:    "Bot_Requests",
		Payload:  string(configJson),
		Routing:  "trader.",
		Monitor:false,
	})

	return
}
