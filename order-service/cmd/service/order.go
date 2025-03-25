package service

import (
	"errors"
	"order-service/cmd/model"
	"order-service/cmd/repository"
	"order-service/cmd/utils"
	"os"
)

type OrderService struct {
	OrderRepo      *repository.OrderRepository
	ProductService *ProductService
}

func NewOrderService(r *repository.OrderRepository) *OrderService {
	return &OrderService{OrderRepo: r,
		ProductService: &ProductService{ProductServiceURL: os.Getenv("PRODUCT_SERVICE_URL")},
	}
}
func (s *OrderService) AddNewOrder(orderRequest *model.OrderRequestDetails) (*model.OrderResponse, error) {
	err := s.ProductService.VerifyOrderItems(orderRequest.StoreID, orderRequest.OrderItems)
	if err != nil {
		return nil, err
	}

	// Payment Logic using payment gateway (To Be Implemented)

	// Manage inventory
	err = s.ProductService.UpdateInventory(orderRequest.OrderItems)
	if err != nil {
		return nil, err
	}

	order, err := s.OrderRepo.AddOrder(orderRequest)
	if err != nil {
		return nil, err
	}
	orderResponse := order.CreateOrderResponse()
	return orderResponse, nil

}

func (s *OrderService) GetAllOrder(storeId string) ([]model.OrderResponse, error) {

	orders, err := s.OrderRepo.GetAllOrder(storeId)
	if err != nil {
		return nil, err
	}
	// mapping order item model into order item response
	var orderResponse []model.OrderResponse
	for _, item := range orders {
		orderResponse = append(orderResponse, *item.CreateOrderResponse())
	}

	return orderResponse, nil
}

func (s *OrderService) GetOrderDetails(orderId string) (*model.OrderDetailsResponse, error) {
	var order model.Order
	err := s.OrderRepo.GetOrderDetails(&order, orderId)
	if err != nil {
		return nil, err
	}

	orderDetailsResponse := order.CreateOrderDetailsResponse()
	return orderDetailsResponse, nil
}

func (s *OrderService) GetOrder(orderId string) (*model.OrderResponse, error) {
	var order model.Order
	err := s.OrderRepo.GetOrder(&order, orderId)
	if err != nil {
		return nil, err
	}

	orderResponse := order.CreateOrderResponse()
	return orderResponse, nil
}

func (s *OrderService) UpdateOrder(orderId uint, orderRequest *model.OrderRequest) error {
	var order model.Order
	rowAffected, err := s.OrderRepo.FindOrder(&order, utils.ItoS(orderId))
	if err != nil {
		return err
	} else if rowAffected == 0 {
		return errors.New("order not found")
	}
	err = s.OrderRepo.UpdateOrder(orderRequest, orderId)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) DeleteOrder(orderId string) error {
	// Delete order  from database by order id
	var order model.Order
	rowAffected, err := s.OrderRepo.FindOrder(&order, orderId)
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
