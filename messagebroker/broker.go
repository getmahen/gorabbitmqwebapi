package messagebroker

import (
	"encoding/json"
	"fmt"
	"log"
	"rabbitmqwebapp/models"
	"time"

	"github.com/streadway/amqp"
)

const queueName = "Port"

type IMessageBroker interface {
	EstablishConnection() (conn *amqp.Connection, channel *amqp.Channel)
	DispatchMessage(requestId string, result models.PortResponse, messageListner IMessageListner) error
}

type IMessageListner interface {
	MessageConsumed(result models.PortResponse)
}

type messageBrokerImplementer struct {
	url            string
	queueName      string
	connection     *amqp.Connection
	channel        *amqp.Channel
	messageListner IMessageListner
}

func NewMessageBroker(url string) IMessageBroker {
	ret := &messageBrokerImplementer{
		url: url,
	}
	ret.EstablishConnection()

	return ret
}

func (p *messageBrokerImplementer) EstablishConnection() (conn *amqp.Connection, channel *amqp.Channel) {
	connection, err := amqp.Dial(p.url)
	p.failOnError(err, "Failed to connect to RabbitMQ")
	println(fmt.Sprintf("Successfully established RabbitMq connection"))
	p.connection = connection

	ch, err := connection.Channel()
	p.failOnError(err, "Failed to open a channel")
	println(fmt.Sprintf("Successfully created RabbitMq channel"))
	p.channel = ch

	p.queueName = queueName

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	p.failOnError(err, "Failed to declare the queue")
	println(fmt.Sprintf("Successfully created Queue named %s in RabbitMq", q.Name))

	//Start listening to the queue to consume messages - This is a gorountine
	p.beginConsumption()

	return connection, ch
}

func (p *messageBrokerImplementer) DispatchMessage(requestId string, result models.PortResponse, messageListner IMessageListner) error {
	if p.connection != nil {
		//messageToSend := fmt.Sprintf("Hey there! Received message for requestId: %s", requestId)
		messageToSend, marshalErr := json.Marshal(result)
		p.failOnError(marshalErr, "Failed to JSON Marshal")

		//Put the message on the queue
		msg := amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(messageToSend),
		}

		err := p.channel.Publish("", p.queueName, false, false, msg)
		p.failOnError(err, "Error sending message to Queue in RabbitMQ")
		println(fmt.Sprintf("Successfully sent message %s to Queue: %s", messageToSend, p.queueName))

		p.messageListner = messageListner
	}
	return nil
}

func (p *messageBrokerImplementer) beginConsumption() {
	println("Started beginConsumption as a goroutine........")

	go func() {
		msgs, err := p.channel.Consume(p.queueName, "", true, false, false, false, nil)
		p.failOnError(err, "Error in consuming message from the queue")

		for msg := range msgs {
			time.Sleep(1 * time.Second)

			target := models.PortResponse{}
			err := json.Unmarshal(msg.Body, &target)
			p.failOnError(err, "Failed to Unmarshal Json response from the queue")

			println(fmt.Sprintf("Message received from the queueName: %s and message: %s", queueName, target))
			p.messageListner.MessageConsumed(target)
		}
	}()
}

func (p *messageBrokerImplementer) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
