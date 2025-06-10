package service

import (
	"errors"
	"fmt"
	"net/http"
	"order-service/cmd/model"
	"order-service/cmd/repository"
	"order-service/cmd/utils"
	"os"
	"time"
)

type OrderService struct {
	OrderRepo      *repository.OrderRepository
	ProductService *ProductService
	UserServiceURL string
	client         *http.Client
}

func NewOrderService(r *repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepo: r,
		ProductService: &ProductService{
			ProductServiceURL: os.Getenv("PRODUCT_SERVICE_URL"),
			client:            &http.Client{Timeout: 5 * time.Second},
		},
		UserServiceURL: os.Getenv("USER_SERVICE_URL"),
		client:         &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *OrderService) AddNewOrder(storeId uint, orderRequest *model.OrderRequestDetails) (*model.OrderResponse, error) {
	err := s.ProductService.VerifyOrderItems(storeId, orderRequest.OrderItems)
	if err != nil {
		return nil, err
	}

	// Payment Logic using payment gateway (To Be Implemented)

	// Manage inventory
	err = s.ProductService.UpdateInventory(orderRequest.OrderItems)
	if err != nil {
		return nil, err
	}

	order, err := s.OrderRepo.AddOrder(storeId, orderRequest)
	if err != nil {
		return nil, err
	}
	orderResponse := order.CreateOrderResponse()
	return orderResponse, nil

}

func (s *OrderService) GetAllOrder(storeId string) (*model.OrdersResponse, error) {
	orders, err := s.OrderRepo.GetAllOrder(storeId)
	if err != nil {
		return nil, err
	}

	/*// Check if no orders were found
	if len(orders) == 0 {
		return nil, errors.New("no orders found")
	}*/

	// mapping order item model into order item response
	ordersResponse := model.GetOrdersReponse(orders)

	return ordersResponse, nil
}

func (s *OrderService) GetOrderDetails(orderId string) (*model.OrderDetailsResponse, error) {
	var order model.Order
	err := s.OrderRepo.GetOrderDetails(&order, orderId)
	if err != nil {
		return nil, err
	}

	// Create basic order details response
	orderDetailsResponse := order.CreateOrderDetailsResponse()

	err = s.ProductService.GetOrderItemDetails(orderDetailsResponse.StoreID, orderDetailsResponse.OrderItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get order item details: %w", err)
	}
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
func (s *OrderService) ChangeOrderStatus(orderId uint, newStatus string) error {
	var order model.Order
	rowAffected, err := s.OrderRepo.IsOrderExist(&order, utils.ItoS(orderId))
	if err != nil {
		return err
	} else if rowAffected == 0 {
		return errors.New("order not found")
	}

	if !model.CanTransition(order.Status, newStatus) {
		return errors.New("invalid status transition from " + order.Status + " to " + newStatus)
	}
	err = s.OrderRepo.ChangeOrderStatus(&order, newStatus)
	if err != nil {
		return err
	}

	return nil
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
