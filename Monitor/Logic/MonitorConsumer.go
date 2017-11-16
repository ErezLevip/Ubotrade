package Logic
/*
import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"local/UbotTrade/Bot/Logic"
	"local/UbotTrade/BotConsumer/Logic"
	"local/UbotTrade/DataHandlers/MongoDB"
	"local/UbotTrade/DataHandlers/RabbitMQ"
	"local/UbotTrade/Global"
	"log"
	"time"
)

type Monitor interface {
	Make(interval time.Duration)
	StartMonitor()
}

type ActiveBotMonitor struct {
	Interval      time.Duration
	DbHandler     MongoDB.MongoHandler
	RabbitHandler RabbitMQ.RabbitHandler
}

func (monitor ActiveBotMonitor) Make(interval time.Duration) {
	interval = interval /// initialize rabbit and mongo handlers
}

func (monitor ActiveBotMonitor) StartMonitor() {
	log.Println("Monitor has started ", time.Now())
	for {
		bots := monitor.getAllActiveBots()
		for i := 0; i < len(bots); i++ {
			var botData Global.BotInformation
			err := json.Unmarshal([]byte(bots[i]["BotData"].(string)), &botData)
			if time.Now().After(bots[i]["LastHealthCheck"].(time.Time).Add(monitor.Interval)) {
				//var payloadJson string
				var payloadBytes, err = json.Marshal(botData)
				if err != nil {
					log.Println(err.Error())
				}
				msg := RabbitMQ.RabbitMessage{
					Payload:  string(payloadBytes),
					SenderId: botData.UserId,
					SentTime: time.Now(),
					Topic:    "Bot_Requests",
					Routing:  "trader.",
					Monitor:true,
				}
				monitor.RabbitHandler.Publish(msg)

			}
			if err != nil {
				log.Println(err.Error())
			}
		}

		time.Sleep(monitor.Interval)
	}
}

func (monitor ActiveBotMonitor) getAllActiveBots() []map[string]interface{} {

	res, err := monitor.DbHandler.FindFMany(Implementations.BotsCollection, bson.M{"IsActive": true})
	if err != nil {
		log.Println(err.Error())
	}
	return res
}*/
