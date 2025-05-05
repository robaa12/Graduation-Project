package messaging

import (
	"errors"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn              *amqp091.Connection
	channel           *amqp091.Channel
	uri               string
	connectionRetries int
	retryDelay        time.Duration
	connected         bool
}

func NewRabbitClient(uri string) *RabbitMQClient {
	return &RabbitMQClient{
		uri:               uri,
		connectionRetries: 5,
		retryDelay:        5 * time.Second,
		connected:         false,
	}
}

func (c *RabbitMQClient) Connect() error {
	var err error

	for i := 0; i < c.connectionRetries; i++ {
		c.conn, err = amqp091.Dial(c.uri)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ: %v. Retrying in %v... (%d/%d) ", err, c.retryDelay, i+1, c.connectionRetries)
		time.Sleep(c.retryDelay)
	}

	if err != nil {
		return errors.New("failed to connect to RabbitMQ after multiple attempts")
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		_ = c.conn.Close()
		return errors.New("failed to open a channel: " + err.Error())
	}

	c.connected = true
	log.Println("Successfully connected to RabbitMQ")

	go c.handleReconnection()

	return nil
}

func (c *RabbitMQClient) handleReconnection() {
	connErr := make(chan *amqp091.Error)
	c.conn.NotifyClose(connErr)

	// Wait for a connection error
	err := <-connErr
	c.connected = false
	log.Printf("RabbitMQ connection closed: %v. Attempting to reconnect...", err)

	// Attempt to handleReconnection
	for {
		if err := c.Connect(); err == nil {
			log.Println("Successfully reconnected to RabbitMQ")
			return
		}

		log.Println("Failed to reconnect to RabbitMQ. Retrying...")
		time.Sleep(c.retryDelay)
	}
}

// Close closes the connection to RabbitMQ
func (c *RabbitMQClient) Close() error {
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *RabbitMQClient) DeclareExchange(name string, kind string, durable bool) error {
	if !c.connected {
		return errors.New("not connected to RabbitMQ")
	}

	return c.channel.ExchangeDeclare(name, kind, durable, false, false, false, nil)
}

func (c *RabbitMQClient) DeclareQueue(name string, durable bool) (amqp091.Queue, error) {
	if !c.connected {
		return amqp091.Queue{}, errors.New("not connected to RabbitMQ")
	}

	return c.channel.QueueDeclare(name, durable, false, false, false, nil)
}

func (c *RabbitMQClient) BindQueue(queueName, key, exchangeName string) error {
	if !c.connected {
		return errors.New("not connected to RabbitMQ")
	}

	return c.channel.QueueBind(queueName, key, exchangeName, false, nil)
}

func (c *RabbitMQClient) GetChannel() *amqp091.Channel {
	return c.channel
}
