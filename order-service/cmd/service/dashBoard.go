package service

import (
	"fmt"
	"net/http"
	"order-service/cmd/model"
	"order-service/cmd/repository"
	"os"
	"time"
)

type OrderDashBoardService struct {
	DashBoardRepository *repository.DashBoardRepository
	ProductService      *ProductService
}

func NewDashBoardService(dashBoardRepository *repository.DashBoardRepository) *OrderDashBoardService {
	return &OrderDashBoardService{DashBoardRepository: dashBoardRepository,
		ProductService: &ProductService{
			ProductServiceURL: os.Getenv("PRODUCT_SERVICE_URL"),
			client:            &http.Client{Timeout: 5 * time.Second},
		},
	}
}
func (s *OrderDashBoardService) GetDashboardInfo(storeID uint, startDate, endDate time.Time) (*model.DashBoardResponse, error) {
	if startDate.IsZero() {
		startDate = time.Now()
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	// Get all required data
	productsDashBoard, err := s.ProductService.GetStoreProductsDashboard(storeID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get products dashboard: %w", err)
	}

	summary, storeSlug, err := s.DashBoardRepository.GetStoreSummary(storeID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get store summary: %w", err)
	}

	monthlySales, err := s.DashBoardRepository.GetMonthlySales(storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly sales: %w", err)
	}

	latestOrders, err := s.DashBoardRepository.GetLatestOrders(storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest orders: %w", err)
	}

	latestCustomers, err := s.DashBoardRepository.GetLatestCustomers(storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest customers: %w", err)
	}

	// Create response
	summary.TotalProducts = productsDashBoard.TotalProducts
	summary.ProductsChange = productsDashBoard.ProductsChange

	dashboardResponse := model.CreateDashboardResponse(
		storeSlug,
		*summary,
		monthlySales,
		model.GetOrderDashboardResponse(latestOrders),
		latestCustomers,
	)

	return &dashboardResponse, nil
}
