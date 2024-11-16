package utils

import (
	"encoding/json"
	"log"

	"github.com/wagslane/go-rabbitmq"
)

type AmqpConfig struct {
	Conn         *rabbitmq.Conn
	Queue        string
	RoutingKey   string
	ExchangeName string
}

type amqpPublishData struct {
	Types   string `json:"types"`
	Payload any    `json:"payload"`
}

func routeMessage(body []byte) error {
	var msg = new(amqpPublishData)
	if err := json.Unmarshal(body, &msg); err != nil {
		return err
	}

	switch msg.Types {
	case "ffmpeg":
		payloadMap, ok := msg.Payload.(map[string]interface{})
		if ok {
			FFmpegExexute(payloadMap["uuid"].(string), payloadMap["location"].(string))
		}
	case "hellow":
		log.Println("hello worlds")
	default:
		log.Printf("Unknown message type: %s", msg.Types)
	}

	return nil
}

func NewAmqpConn() *AmqpConfig {
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &AmqpConfig{
		Conn:         conn,
		Queue:        "gostream",
		RoutingKey:   "gostream_routekey",
		ExchangeName: "event",
	}
}

func (a *AmqpConfig) NewConsumer() {
	consumer, err := rabbitmq.NewConsumer(
		a.Conn,
		a.Queue,
		rabbitmq.WithConsumerOptionsRoutingKey(a.RoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(a.ExchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	err = consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		// log.Printf("consumed: %v", string(d.Body))
		if err := routeMessage(d.Body); err != nil {
			log.Printf("Error processing message: %v", err)
			return rabbitmq.NackRequeue
		}
		// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
		return rabbitmq.Ack
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (a *AmqpConfig) NewPublisher(types string, payload any) error {
	message := amqpPublishData{
		Types:   types,
		Payload: payload,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	publisher, err := rabbitmq.NewPublisher(
		a.Conn,
		rabbitmq.WithPublisherOptionsExchangeName(a.ExchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	return publisher.Publish(
		data,
		[]string{a.RoutingKey},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(a.ExchangeName),
	)
}
