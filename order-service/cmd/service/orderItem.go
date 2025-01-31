package service

import (
	"errors"
	"log"
	"order-service/cmd/model"
	"order-service/cmd/repository"
)

type OrderItemService struct {
	OrderItemRepo *repository.OrderItemRepository
}

func NewOrderItemService(orderItemRepo *repository.OrderItemRepository) *OrderItemService {
	return &OrderItemService{OrderItemRepo: orderItemRepo}
}

func (s *OrderItemService) AddOrderItem(orderId uint, orderItemRequest *model.OrderItemRequest) (*model.OrderItemResponse, error) {

	// mapping order item request into order Item Model and add to database
	orderItem := orderItemRequest.CreateOrderItem(orderId)
	err := s.OrderItemRepo.AddOrderItem(orderItem)
	if err != nil {
		log.Println("Error, Creating orderItem in database")
		return nil, err
	}

	// mapping orderitem model into order item response
	orderItemResponse := orderItem.CreateOrderItemResponse()
	return orderItemResponse, nil
}

func (s *OrderItemService) GetAllOrderItems(order_id string) ([]model.OrderItemResponse, error) {

	// Query to Get all order item from database
	orderItems, err := s.OrderItemRepo.GetAllOrderItems(order_id)
	if err != nil {

		return nil, err
	}

	// mapping orderitem model into order item response
	var orderItemsResponse []model.OrderItemResponse
	for _, item := range orderItems {
		orderItemsResponse = append(orderItemsResponse, *item.CreateOrderItemResponse())
	}

	return orderItemsResponse, nil
}
func (s *OrderItemService) GetOrderItem(orderItemId string) (*model.OrderItemResponse, error) {
	// Query to Get order item from database
	var orderItem model.OrderItem
	err := s.OrderItemRepo.GetOrderItem(&orderItem, orderItemId)
	if err != nil {
		return nil, err
	}
	// mapping orderitem model into order item response
	orderItemResponse := orderItem.CreateOrderItemResponse()
	return orderItemResponse, nil
}
func (s *OrderItemService) UpdateOrderItem(orderItem_id string, orderItemRequest *model.OrderItemRequest) error {

	var orderItem model.OrderItem
	rowAffected, err := s.OrderItemRepo.FindOrderItem(&orderItem, orderItem_id)
	if err != nil {
		return err
	} else if rowAffected == 0 {
		return errors.New("order item not found")
	}
	// verfiy sku id comming from request as same in database
	if orderItem.SkuID != orderItemRequest.SkuID && orderItemRequest.SkuID != 0 {
		return errors.New("can't change sku for item")
	}

	err = s.OrderItemRepo.UpdateOrderItem(orderItemRequest, orderItem_id)
	if err != nil {
		return err
	}

	return nil
}
func (s *OrderItemService) DeleteOrderItem(orderItem_id string) error {
	// Delete orderItem  from database by order id
	var orderItem model.OrderItem
	rowAffected, err := s.OrderItemRepo.FindOrderItem(&orderItem, orderItem_id)
	if err != nil {
		return err
	} else if rowAffected == 0 {
		return errors.New("order item not found")
	}
	err = s.OrderItemRepo.DeleteOrderItem(&orderItem)
	if err != nil {
		return err
	}
	return nil

}
