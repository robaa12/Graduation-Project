package repository

import (
	"order-service/cmd/model"
	"time"

	"gorm.io/gorm"
)

type DashBoardRepository struct {
	db *gorm.DB
}

func NewDashBoardRepository(db *gorm.DB) *DashBoardRepository {
	return &DashBoardRepository{db: db}
}
func (d *DashBoardRepository) GetMonthlySales(storeID uint) ([]model.MonthlySales, error) {
	monthlySales := []model.MonthlySales{}

	createdAt := time.Now().AddDate(0, -12, 0)

	result := d.db.Table("orders").
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, COALESCE(SUM(total_price), 0) as sales").
		Where("store_id = ? AND created_at >= ? AND status != ?", storeID, createdAt, "cancelled").
		Group("month").
		Order("month").
		Scan(&monthlySales)

	if result.Error != nil {
		return nil, result.Error
	}

	return monthlySales, nil
}
func (d *DashBoardRepository) GetLatestOrders(storeID uint) ([]model.Order, error) {
	var latestOrders []model.Order
	result := d.db.Table("orders").
		Where("store_id = ? AND status != 'cancelled'", storeID).
		Order("created_at DESC").
		Limit(7).
		Find(&latestOrders)
	if result.Error != nil {
		return nil, result.Error
	}
	return latestOrders, nil
}
func (d *DashBoardRepository) GetLatestCustomers(storeID uint) ([]model.CustomerDashboardResponse, error) {
	latestCustomers := []model.CustomerDashboardResponse{}

	result := d.db.Table("customers").
		Select(`
            customers.id AS customer_id,
            customers.email AS customer_name,
            COUNT(orders.id) AS number_of_orders,
            COALESCE(SUM(orders.total_price), 0) AS total_spent,
            TO_CHAR(customers.created_at, 'YYYY-MM') AS join_date
        `).
		Joins("JOIN store_customers ON store_customers.customer_id = customers.id").
		Joins("LEFT JOIN orders ON orders.customer_id = customers.id AND orders.store_id = ?", storeID).
		Where("store_customers.store_id = ?", storeID).
		Group("customers.id, customers.email, customers.created_at").
		Order("customers.created_at DESC").
		Limit(7).
		Scan(&latestCustomers)
	if result.Error != nil {
		return nil, result.Error
	}
	return latestCustomers, nil
}

func (r *DashBoardRepository) GetStoreSummary(storeID uint, startDate, endDate time.Time) (*model.Summary, string, error) {
	var summary model.Summary

	// Get store slug
	var storeSlug string
	if err := r.db.Model(&model.Store{}).
		Select("slug").
		Where("id = ?", storeID).
		Scan(&storeSlug).Error; err != nil {
		return nil, "", err
	}

	// Current period totals - consider only non-cancelled orders
	if err := r.db.Model(&model.Order{}).
		Where("store_id = ? AND status != ?", storeID, "cancelled").
		Count(&summary.TotalOrders).
		Select("COALESCE(SUM(total_price), 0)").
		Row().
		Scan(&summary.TotalRevenue); err != nil {
		return nil, "", err
	}

	// Get previous period metrics for comparison
	periodDuration := endDate.Sub(startDate)
	prevStart := startDate.Add(-periodDuration)
	prevEnd := startDate

	var prevRevenue float64
	var prevOrders int64
	if err := r.db.Model(&model.Order{}).
		Where("store_id = ? AND created_at BETWEEN ? AND ? AND status != ?",
			storeID, prevStart, prevEnd, "cancelled").
		Count(&prevOrders).
		Select("COALESCE(SUM(total_price), 0)").
		Row().
		Scan(&prevRevenue); err != nil {
		return nil, "", err
	}

	// Calculate percentage changes
	if prevRevenue > 0 {
		summary.RevenueChange = ((summary.TotalRevenue - prevRevenue) / prevRevenue) * 100.0
	}
	if prevOrders > 0 {
		summary.OrdersChange = (float64(summary.TotalOrders-prevOrders) / float64(prevOrders)) * 100.0
	}

	return &summary, storeSlug, nil
}
