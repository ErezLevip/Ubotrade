package RabbitMQ

import (
	"encoding/asn1"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type RabbitConfiguration struct {
	ConnectionString string
	Routing          string
	Topic            string
	IsConsumer       bool
}

type RabbitHandler struct {
	Channel *amqp.Channel
	Q       *amqp.Queue
	Config  RabbitConfiguration
}

type IRabbitHandler interface {
	Publish(message RabbitMessage)
	Consume(routing string, onMessageHandler OnMessageHandler, failureMessageHandler OnMessageHandler, bulkSize int)
	Init(config RabbitConfiguration)
}

func (self *RabbitHandler) Publish(message RabbitMessage) {
	body, err := asn1.Marshal(message)
	self.handleErrors(err)
	log.Println(time.Now(),"msg topic", message.Topic)
	fmt.Println(time.Now(),message.Routing)
	err = self.Channel.Publish(
		message.Topic,                                  // exchange
		fmt.Sprintf(message.Routing, message.SenderId), // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	self.handleErrors(err)

	log.Printf(" [x] Sent %s", message)
}

type OnMessageHandler func([]byte) bool

func (self *RabbitHandler) Consume(routing string, onMessageHandler OnMessageHandler, failureMessageHandler OnMessageHandler, bulkSize int) {
	self.Channel.Qos(bulkSize, 0, false)
	msgs, err := self.Channel.Consume(
		self.Q.Name, // queue
		"",          // consumer
		true,        // auto ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // args
	)
	self.handleErrors(err)
	forever := make(chan bool)

	go func() {
		var i = 0
		for d := range msgs {
			i++
			if !onMessageHandler(d.Body) {
				self.retryMessageHandler(d.Body, onMessageHandler, failureMessageHandler)
			}
		}
	}()

	<-forever
}
func (self *RabbitHandler) retryMessageHandler(message []byte, messageHandler OnMessageHandler, errorMessageHandler OnMessageHandler) {
	success := false
	for tryCount := 0; tryCount <= 3; tryCount++ {
		success = messageHandler(message)
		if success {
			break
		}
	}
	if !success {
		messageObj := RabbitMessage{}
		_, err := asn1.Unmarshal(message, &messageObj)
		if err != nil {
			log.Panic(err.Error())
		}

		body, err := asn1.Marshal(messageObj)
		self.handleErrors(err)

		self.Channel.Publish(
			self.Config.Topic,                                    // exchange
			fmt.Sprintf(messageObj.Routing, messageObj.SenderId), // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        body,
			})
		self.handleErrors(err)

		log.Printf(" [x] Sent %s", message)
		errorMessageHandler(message)
	}
}

func (self *RabbitHandler) Init(config RabbitConfiguration) {
	self.Config = config //only for consumers
	conn, err := amqp.Dial(config.ConnectionString)
	self.handleErrors(err)

	ch, err := conn.Channel()
	self.handleErrors(err)

	err = ch.ExchangeDeclare(
		config.Topic, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	self.handleErrors(err)
	self.Channel = ch

	if(config.IsConsumer){
		q, err := self.Channel.QueueDeclare(
			"",    // name
			false, // durable
			false, // delete when usused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		self.Q = &q
		self.handleErrors(err)

		err = self.Channel.QueueBind(
			q.Name,         // queue name
			config.Routing, // routing key
			config.Topic,   // exchange
			false,
			nil)
		if err != nil {
			log.Panic(err.Error())
		}
	}
}


func (self *RabbitHandler) handleErrors(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}
