package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MessagingService provides high-level messaging functionality
type MessagingService struct {
	client      *RabbitMQClient
	publishers  map[string]*Publisher
	consumers   map[string]*Consumer
	responses   map[string]chan []byte
	responseMu  sync.RWMutex
	initialized bool
}

func NewMessagingService(rabbitMQURI string) (*MessagingService, error) {
	client := NewRabbitClient(rabbitMQURI)
	if err := client.Connect(); err != nil {
		return nil, err
	}

	service := &MessagingService{
		client:     client,
		publishers: make(map[string]*Publisher),
		consumers:  make(map[string]*Consumer),
		responses:  make(map[string]chan []byte),
	}

	exchanges := map[string]string{
		OrderExchange:     "topic",
		ProductExchange:   "topic",
		InventoryExchange: "topic",
	}

	for name, kind := range exchanges {
		if err := client.DeclareExchange(name, kind, true); err != nil {
			client.Close()
			return nil, errors.New("failed to declare exchange " + name + ":" + err.Error())
		}
	}
	return service, nil
}

func (s *MessagingService) Initialize() error {
	if s.initialized {
		return nil
	}

	// Set up response queue
	responseQueueName := "response-" + uuid.NewString()
	responseQueue, err := s.client.DeclareQueue(responseQueueName, false)
	if err != nil {
		return err
	}

	// Bind response queue to all response routing keys
	if err := s.client.BindQueue(responseQueue.Name, "response.*", InventoryExchange); err != nil {
		return err
	}

	// Set up response consumer
	responseConsumer := NewConsumer(s.client, responseQueue.Name, "")
	responseConsumer.SetAutoAck(true)
	if err := responseConsumer.Start(); err != nil {
		return err
	}

	// Process responses
	responseConsumer.Listen(s.handleResponse)

	// Set up common publishers
	publishers := map[string]struct {
		exchange   string
		routingKey string
	}{
		"inventory_verification": {InventoryExchange, "inventory.verification"},
		"inventory_update":       {InventoryExchange, "inventory.update"},
		"order_created":          {OrderExchange, "order.created"},
	}

	for name, config := range publishers {
		s.publishers[name] = NewPublisher(s.client, config.exchange, config.routingKey)
	}

	s.initialized = true
	return nil
}

func (s *MessagingService) handleResponse(body []byte) error {
	var base BaseMessage
	if err := json.Unmarshal(body, &base); err != nil {
		return err
	}

	if base.CorrelID == "" {
		return nil // Ignore messages without correlation IDs
	}

	s.responseMu.RLock()
	ch, exists := s.responses[base.CorrelID]
	s.responseMu.RUnlock()

	if exists {
		ch <- body
	}

	return nil
}

func (s *MessagingService) VerifyInventory(ctx context.Context, storeID uint, items []InventoryItem) (*InventoryVerificationResponse, error) {
	if !s.initialized {
		if err := s.Initialize(); err != nil {
			return nil, err
		}
	}

	correlID := uuid.NewString()

	responseCh := make(chan []byte, 1)
	s.responseMu.Lock()
	s.responses[correlID] = responseCh
	s.responseMu.Unlock()

	// Clean up when done
	defer func() {
		s.responseMu.Lock()
		delete(s.responses, correlID)
		s.responseMu.Unlock()
		close(responseCh)
	}()

	// Create request
	req := InventoryVerificationRequest{
		BaseMessage: BaseMessage{
			ID:        uuid.NewString(),
			Type:      InventoryUpdateRequestType,
			Timestamp: time.Now(),
			CorrelID:  correlID,
		},
		StoreID: storeID,
		Items:   items,
	}

	// Get publishers
	pub, exists := s.publishers["inventory_verification"]
	if !exists {
		return nil, errors.New("inventory verification publisher not configured")
	}

	// Send request
	if err := pub.PublishJSON(req); err != nil {
		return nil, err
	}

	select {
	case responseBody := <-responseCh:
		var response InventoryVerificationResponse
		if err := json.Unmarshal(responseBody, &response); err != nil {
			return nil, err
		}
		return &response, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(10 * time.Second):
		return nil, errors.New("timeout waiting for inventory verification response")
	}
}

func (s *MessagingService) UpdateInventory(storeID uint, orderID uint, items []InventoryItem) error {
	if !s.initialized {
		if err := s.Initialize(); err != nil {
			return err
		}
	}

	req := InventoryUpdateRequest{
		BaseMessage: BaseMessage{
			ID:        uuid.NewString(),
			Type:      InventoryUpdateRequestType,
			Timestamp: time.Now(),
		},
		OrderID: orderID,
		StoreID: storeID,
		Items:   items,
	}

	pub, exists := s.publishers["inventory_update"]
	if !exists {
		return errors.New("inventory update publisher not configured")
	}

	// Send request (fire and forget)
	return pub.PublishJSON(req)
}

func (s *MessagingService) PublishOrderCreated(orderID, storeID uint, customerInfo CustomerInfo, items []OrderItemInfo, total float64) error {
	if !s.initialized {
		if err := s.Initialize(); err != nil {
			return err
		}
	}

	// Create event
	event := OrderCreatedEvent{
		BaseMessage: BaseMessage{
			ID:        uuid.NewString(),
			Type:      OrderCreatedEventType,
			Timestamp: time.Now(),
		},
		OrderID:  orderID,
		StoreID:  storeID,
		Customer: customerInfo,
		Items:    items,
		Total:    total,
	}

	// Get publisher
	pub, exists := s.publishers["order_created"]
	if !exists {
		return errors.New("order created publisher not configured")
	}

	return pub.PublishJSON(event)
}

func (s *MessagingService) Close() error {
	return s.client.Close()
}
