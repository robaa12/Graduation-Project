package messaging

import (
	"encoding/json"
	"errors"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	client     *RabbitMQClient
	exchange   string
	routingKey string
	mandatory  bool
	immediate  bool
}

func NewPublisher(client *RabbitMQClient, exchange, routingKey string) *Publisher {
	return &Publisher{
		client:     client,
		exchange:   exchange,
		routingKey: routingKey,
		mandatory:  false,
		immediate:  false,
	}
}

func (p *Publisher) PublishJSON(message interface{}) error {
	if p.client == nil || !p.client.connected {
		return errors.New("RabbitMQ client not connected")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return errors.New("failed to marshal message to JSON: " + err.Error())
	}

	return p.Publish(data, "application/json")
}

func (p *Publisher) Publish(body []byte, contentType string) error {
	if p.client == nil || !p.client.connected {
		return errors.New("RabbitMQ client not connected")
	}

	return p.client.channel.Publish(p.exchange, p.routingKey, p.mandatory, p.immediate, amqp091.Publishing{
		ContentType:  contentType,
		Body:         body,
		DeliveryMode: amqp091.Persistent,
	})
}

func (p *Publisher) SetMandatory(mandatory bool) {
	p.mandatory = mandatory
}

func (p *Publisher) SetImmediate(immediate bool) {
	p.immediate = immediate
}
