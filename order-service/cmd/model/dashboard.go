package model

type MonthlySales struct {
	Month string  `json:"month"`
	Sales float64 `json:"sales"`
}
type ProductsDashboardResponse struct {
	TotalProducts  int64   `json:"totalProducts"`
	ProductsChange float64 `json:"productsChange"`
}
type Summary struct {
	TotalRevenue   float64 `json:"totalRevenue"`
	TotalProducts  int64   `json:"totalProducts"`
	TotalOrders    int64   `json:"totalOrders"`
	RevenueChange  float64 `json:"revenueChange"`
	ProductsChange float64 `json:"productsChange"`
	OrdersChange   float64 `json:"ordersChange"`
}

type DashBoardResponse struct {
	Store_Slug      string                      `json:"storeSlug"`
	Summary         Summary                     `json:"summary"`
	MonthlySales    []MonthlySales              `json:"monthlySales"`
	LatestOrders    []OrderDashboardResponse    `json:"latestOrders"`
	LatestCustomers []CustomerDashboardResponse `json:"latestCustomers"`
}

func (p *ProductsDashboardResponse) CreateSummary(summary Summary) Summary {
	return Summary{
		TotalRevenue:   summary.TotalRevenue,
		TotalProducts:  p.TotalProducts,
		TotalOrders:    summary.TotalOrders,
		RevenueChange:  summary.RevenueChange,
		ProductsChange: p.ProductsChange,
		OrdersChange:   summary.OrdersChange,
	}
}
func CreateDashboardResponse(storeSlug string, summary Summary, monthlySales []MonthlySales, latestOrders []OrderDashboardResponse, latestCustomers []CustomerDashboardResponse) DashBoardResponse {
	return DashBoardResponse{
		Store_Slug:      storeSlug,
		Summary:         summary,
		MonthlySales:    monthlySales,
		LatestOrders:    latestOrders,
		LatestCustomers: latestCustomers}
}
