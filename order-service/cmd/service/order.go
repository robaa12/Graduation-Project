package service

import (
	"errors"
	"order-service/cmd/model"
	"order-service/cmd/repository"
	"order-service/cmd/utils"
)

type OrderService struct {
	OrderRepo *repository.OrderRepository
}

func NewOrderService(r *repository.OrderRepository) *OrderService {
	return &OrderService{OrderRepo: r}
}
func (s *OrderService) AddNewOrder(orderRequest *model.OrderRequestDetails) (*model.OrderResponse, error) {
	order := orderRequest.CreateOrder()
	err := s.OrderRepo.AddOrder(order)
	if err != nil {
		return nil, err
	}
	orderResponse := order.CreateOrderResponse()
	return orderResponse, nil

}

func (s *OrderService) GetAllOrder(store_id string) ([]model.OrderResponse, error) {

	orders, err := s.OrderRepo.GetAllOrder(store_id)
	if err != nil {
		return nil, err
	}
	// mapping orderitem model into order item response
	var orderResponse []model.OrderResponse
	for _, item := range orders {
		orderResponse = append(orderResponse, *item.CreateOrderResponse())
	}

	return orderResponse, nil
}

func (s *OrderService) GetOrderDetails(order_id string) (*model.OrderDetailsResponse, error) {
	var order model.Order
	err := s.OrderRepo.GetOrderDetails(&order, order_id)
	if err != nil {
		return nil, err
	}

	orderDetailsResponse := order.CreateOrderDetailsResponse()
	return orderDetailsResponse, nil
}

func (s *OrderService) GetOrder(order_id string) (*model.OrderResponse, error) {
	var order model.Order
	err := s.OrderRepo.GetOrder(&order, order_id)
	if err != nil {
		return nil, err
	}

	orderResponse := order.CreateOrderResponse()
	return orderResponse, nil
}

func (s *OrderService) UpdateOrder(order_id uint, orderRequest *model.OrderRequest) error {
	var order model.Order
	rowAffected, err := s.OrderRepo.FindOrder(&order, utils.ItoS(order_id))
	if err != nil {
		return err
	} else if rowAffected == 0 {
		return errors.New("order not found")
	}
	err = s.OrderRepo.UpdateOrder(orderRequest, order_id)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) DeleteOrder(order_id string) error {
	// Delete order  from database by order id
	var order model.Order
	rowAffected, err := s.OrderRepo.FindOrder(&order, order_id)
	if err != nil {
		return err
	} else if rowAffected == 0 {
		return errors.New("order not found")
	}
	err = s.OrderRepo.DeleteOrder(&order)
	if err != nil {
		return err
	}
	return nil

}
