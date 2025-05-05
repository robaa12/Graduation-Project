package messaging

import (
	"errors"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type MessageHandler func([]byte) error

type Consumer struct {
	client     *RabbitMQClient
	queueName  string
	consumer   string
	autoAck    bool
	exclusive  bool
	noLocal    bool
	noWait     bool
	args       amqp091.Table
	deliveries <-chan amqp091.Delivery
}

func NewConsumer(client *RabbitMQClient, queueName, consumer string) *Consumer {
	return &Consumer{
		client:    client,
		queueName: queueName,
		consumer:  consumer,
		autoAck:   false,
		exclusive: false,
		noLocal:   false,
		noWait:    false,
		args:      nil,
	}
}

func (c *Consumer) Start() error {
	if c.client == nil || !c.client.connected {
		return errors.New("RabbitMQ client not connected")
	}

	var err error
	c.deliveries, err = c.client.channel.Consume(c.queueName, c.consumer, c.autoAck, c.exclusive, c.noLocal, c.noWait, c.args)

	if err != nil {
		return errors.New("failed to start consuming: " + err.Error())
	}

	return nil
}

func (c *Consumer) Listen(handler MessageHandler) {
	if c.deliveries == nil {
		log.Println("Error: Deliveries channel is nil. Did you call Start()?")
		return
	}

	go func() {
		for delivery := range c.deliveries {
			if err := handler(delivery.Body); err != nil {
				log.Printf("Error handling message: %v", err)
				// Negative acknowleadge so the message is requeued
				delivery.Nack(false, true)
			} else {
				// Acknowledge successful processing
				delivery.Ack(false)
			}
		}
		log.Println("Deliveries channel closed")
	}()
}

func (c *Consumer) SetAutoAck(autoAck bool) {
	c.autoAck = autoAck
}

func (c *Consumer) SetExclusive(exclusive bool) {
	c.exclusive = exclusive
}
