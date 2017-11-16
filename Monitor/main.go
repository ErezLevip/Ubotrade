package main
/*
import (
	"local/UbotTrade/DataHandlers/MongoDB"
	"local/UbotTrade/DataHandlers/RabbitMQ"
	"local/UbotTrade/Monitor/Logic"
	"time"
)

func main() {

	var monitor Logic.Monitor
	dbHandler := MongoDB.MongoHandler{}
	dbConfig := make(map[string]string)
	dbConfig[MongoDB.ConnectionStringKey] = "localhost"
	dbConfig[MongoDB.DatabaseKey] = "DEV"
	dbHandler.Init(dbConfig)

	rabbitHandler := RabbitMQ.RabbitHandler{}
	rabbitHandler.Init("amqp://guest:guest@localhost:5672/", "Bot_Requests")

	monitor = Logic.ActiveBotMonitor{
		Interval:      time.Second * time.Duration(2),
		DbHandler:     dbHandler,
		RabbitHandler: rabbitHandler,
	}

	monitor.StartMonitor()

}
*/