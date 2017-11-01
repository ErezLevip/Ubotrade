package RabbitMQ

import (
	"time"
)

type RabbitMessage struct {
	Topic    string
	Payload  string
	SentTime time.Time
	Routing  string
	SenderId string
	Monitor bool
}
