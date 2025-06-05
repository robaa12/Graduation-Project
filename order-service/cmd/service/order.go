package service

import (
	"encoding/json"
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

type UserServiceResponse struct {
	Status  bool         `json:"status"`
	Message string       `json:"message"`
	Data    StoreDetails `json:"data"`
}

type StoreDetails struct {
	ID            uint    `json:"id"`
	StoreName     string  `json:"store_name"`
	Href          *string `json:"href"`
	Slug          *string `json:"slug"`
	Description   string  `json:"description"`
	BusinessPhone string  `json:"business_phone"`
	CategoryID    uint    `json:"category_id"`
	PlanID        uint    `json:"plan_id"`
	StoreCurrency string  `json:"store_currency"`
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

// GetStoreDetails fetches store details from the user service
func (s *OrderService) GetStoreDetails(storeID uint) (*StoreDetails, error) {
	// Create request to user service
	url := fmt.Sprintf("%s/store/%d", s.UserServiceURL, storeID)
	
	// Make GET request
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get store details, status code: %d", resp.StatusCode)
	}
	
	// Parse response
	var response UserServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse store details: %w", err)
	}
	
	// Debug log to check response content
	fmt.Printf("Store details response: %+v\n", response)
	
	return &response.Data, nil
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
	
	// Check if no orders were found
	if len(orders) == 0 {
		return nil, errors.New("no orders found")
	}
	
	// Get store details
	storeIdUint := utils.StoUint(storeId)
	storeDetails, err := s.GetStoreDetails(storeIdUint)
	if err != nil {
		// Log error but continue without store name
		fmt.Printf("Warning: Failed to get store details: %v\n", err)
	} else {
		fmt.Printf("Successfully retrieved store details. Store name: %s\n", storeDetails.StoreName)
	}
	
	// mapping order item model into order item response
	var orderResponse []model.OrderResponse
	for _, item := range orders {
		response := item.CreateOrderResponse()
		
		// Add store name if available
		if storeDetails != nil {
			response.OrderResponseInfo.StoreName = storeDetails.StoreName
		}
		
		orderResponse = append(orderResponse, *response)
	}

	return orderResponse, nil
}

func (s *OrderService) GetOrderDetails(orderId string) (*model.OrderDetailsResponse, error) {
	var order model.Order
	err := s.OrderRepo.GetOrderDetails(&order, orderId)
	if err != nil {
		return nil, err
	}

	// Create basic order details response
	orderDetailsResponse := order.CreateOrderDetailsResponse()
	
	// Get store details
	storeDetails, err := s.GetStoreDetails(order.StoreID)
	if err != nil {
		// Log error but continue without store name
		fmt.Printf("Warning: Failed to get store details: %v\n", err)
	} else if storeDetails != nil {
		// Add store name to response
		orderDetailsResponse.OrderResponseInfo.StoreName = storeDetails.StoreName
	}
	
	// Collect all SKU IDs from order items
	var skuIDs []uint
	for _, item := range orderDetailsResponse.OrderItems {
		skuIDs = append(skuIDs, item.SkuID)
	}
	
	// Only fetch SKU details if there are items
	if len(skuIDs) > 0 {
		// Get SKU details from product service
		skuDetails, err := s.ProductService.GetSkuDetails(order.StoreID, skuIDs)
		if err != nil {
			// Log error but continue without SKU details
			fmt.Printf("Warning: Failed to get SKU details: %v\n", err)
		} else {
			// Enrich order items with SKU details
			for i, item := range orderDetailsResponse.OrderItems {
				if details, found := skuDetails[item.SkuID]; found {
					orderDetailsResponse.OrderItems[i].SkuName = details.SkuName
					orderDetailsResponse.OrderItems[i].ProductID = details.ProductID
					orderDetailsResponse.OrderItems[i].ProductName = details.ProductName
					orderDetailsResponse.OrderItems[i].ImageURL = details.ImageURL
				}
			}
		}
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
	
	// Get store details
	storeDetails, err := s.GetStoreDetails(order.StoreID)
	if err != nil {
		// Log error but continue without store name
		fmt.Printf("Warning: Failed to get store details: %v\n", err)
	} else if storeDetails != nil {
		// Add store name to response
		orderResponse.OrderResponseInfo.StoreName = storeDetails.StoreName
	}
	
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
