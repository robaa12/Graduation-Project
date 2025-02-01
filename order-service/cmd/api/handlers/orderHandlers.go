package handlers

import (
	"errors"
	"net/http"
	"order-service/cmd/model"
	"order-service/cmd/service"
	"order-service/cmd/utils"
)

type OrderHandler struct {
	OrderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{OrderService: orderService}
}
func (orderHandler *OrderHandler) AddNewOrder(w http.ResponseWriter, r *http.Request) {
	// Read Order Request from json
	var orderRequest model.OrderRequestDetails
	err := utils.ReadJSON(w, r, &orderRequest)
	if err != nil {
		utils.ErrorJSON(w, errors.New("enter valid order item data"))
		return
	}

	// give orderitem respone from service layer
	orderResponse, err := orderHandler.OrderService.AddNewOrder(&orderRequest)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// write json resonpse
	err = utils.WriteJSON(w, 201, orderResponse)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

}

func (orderHandler *OrderHandler) GetAllOrder(w http.ResponseWriter, r *http.Request) {
	// get store_id from Query parmeter
	store_id, err := utils.GetID(r, "store_id")
	if err != nil {

		utils.ErrorJSON(w, err)
		return
	}
	// get Order response from service layer
	orderResponse, err := orderHandler.OrderService.GetAllOrder(utils.ItoS(store_id))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// write json resonpse
	err = utils.WriteJSON(w, 200, orderResponse)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

}
func (orderHandler *OrderHandler) GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	// get order_id from Query parmeter
	order_id, err := utils.GetID(r, "order_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// get Orderitem response from service layer
	orderDetailsResponse, err := orderHandler.OrderService.GetOrderDetails(utils.ItoS(order_id))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// write json resonpse
	err = utils.WriteJSON(w, 200, orderDetailsResponse)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

}
func (orderHandler *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// get order_id from Query parmeter
	order_id, err := utils.GetID(r, "order_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// get Orderitem response from service layer
	orderResponse, err := orderHandler.OrderService.GetOrder(utils.ItoS(order_id))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// write json resonpse
	err = utils.WriteJSON(w, 200, orderResponse)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

}
func (orderHandler *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	// get order_id from Query parmeter
	order_id, err := utils.GetID(r, "order_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	//Read order request from json
	var orderRequest *model.OrderRequest
	err = utils.ReadJSON(w, r, &orderRequest)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// get Order response from service layer
	err = orderHandler.OrderService.UpdateOrder(order_id, orderRequest)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// write json response
	err = utils.WriteJSON(w, 200, "Order updated successfully.")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

}
func (orderHandler *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// get order_id from Query parmeter
	order_id, err := utils.GetID(r, "order_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// get Orderitem response from service layer
	err = orderHandler.OrderService.DeleteOrder(utils.ItoS(order_id))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// write json resonpse
	err = utils.WriteJSON(w, 200, "Order deleted successfully.")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

}
