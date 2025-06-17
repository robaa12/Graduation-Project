package main

import (
	"net/http"
	"order-service/cmd/api/handlers"
	"order-service/cmd/repository"
	"order-service/cmd/service"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

var db *gorm.DB

func (app *Config) routes() http.Handler {
	db = app.db
	mux := chi.NewRouter()
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Route("/orders/{order_id}/items", orderItems)
	mux.Route("/stores/{store_id}/orders", order)
	mux.Route("/stores/{store_id}/customers", customer)
	mux.Route("/stores/{store_id}/dashboard", dashboard)
	mux.Route("/stores", store)

	return mux
}
func store(r chi.Router) {
	storeRepo := repository.NewStoreRepository(db)
	storeService := service.NewStoreService(storeRepo)
	storeHandler := handlers.NewStoreHandler(storeService)

	r.Post("/", storeHandler.CreateStore)
	r.Route("/{store_id}", func(r chi.Router) {
		r.Delete("/", storeHandler.DeleteStore)
	})
}
func order(r chi.Router) {
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	r.Post("/", orderHandler.AddNewOrder)
	r.Get("/", orderHandler.GetAllOrder)
	r.Route("/{order_id}", func(r chi.Router) {
		r.Get("/", orderHandler.GetOrder)
		r.Get("/details", orderHandler.GetOrderDetails)
		r.Put("/", orderHandler.UpdateOrder)
		r.Put("/status/{status}", orderHandler.UpdateOrderStatus)
		r.Delete("/", orderHandler.DeleteOrder)
	})

}
func orderItems(r chi.Router) {
	orderItemsRepo := repository.NewOrderItemRepository(db)
	orderItemService := service.NewOrderItemService(orderItemsRepo)
	orderItemsHandler := handlers.NewOrderItemsHandler(orderItemService)

	r.Post("/", orderItemsHandler.AddOrderItem)
	r.Get("/", orderItemsHandler.GetAllOrderItems)

	r.Route("/{item_id}", func(r chi.Router) {
		r.Get("/", orderItemsHandler.GetOrderItem)
		r.Put("/", orderItemsHandler.UpdateOrderItem)
		r.Delete("/", orderItemsHandler.DeleteOrderItem)

	})

}
func customer(r chi.Router) {
	customerrRepo := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerrRepo)
	customerHandler := handlers.NewCustomerHandler(customerService)

	r.Post("/", customerHandler.CreateNewCustomer)
	r.Get("/", customerHandler.GetAllCustomers)
	r.Route("/{customer_id}", func(r chi.Router) {
		r.Get("/", customerHandler.GetCustomer)
		//	r.Put("/", customerHandler.UpdateCustomer)
		r.Delete("/", customerHandler.DeleteCustomer)
	})

}
func dashboard(r chi.Router) {
	dashboardRepo := repository.NewDashBoardRepository(db)
	dashboardService := service.NewDashBoardService(dashboardRepo)
	dashboardHandler := handlers.NewDashBoardHandler(dashboardService)

	r.Get("/", dashboardHandler.GetDashboardInfo)

}
