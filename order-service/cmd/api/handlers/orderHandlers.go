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
	// get store_id from Query parameter
	storeId, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// Read Order Request from json
	var orderRequest model.OrderRequestDetails
	err = utils.ReadJSON(w, r, &orderRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, errors.New("enter valid order item data"))
		return
	}

	// give order item response from service layer
	orderResponse, err := orderHandler.OrderService.AddNewOrder(storeId, &orderRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// write json response
	err = utils.WriteJSON(w, 201, orderResponse)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}

func (orderHandler *OrderHandler) GetAllOrder(w http.ResponseWriter, r *http.Request) {
	// get store_id from Query parameter
	storeId, err := utils.GetID(r, "store_id")
	if err != nil {

		_ = utils.ErrorJSON(w, err)
		return
	}
	// get Order response from service layer
	orderResponse, err := orderHandler.OrderService.GetAllOrder(utils.ItoS(storeId))
	if err != nil {
		if err.Error() == "no orders found" {
			_ = utils.ErrorJSON(w, errors.New("no orders found"), http.StatusNotFound)
			return
		}
		_ = utils.ErrorJSON(w, err)
		return
	}

	// write json response
	err = utils.WriteJSON(w, 200, orderResponse)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}
func (orderHandler *OrderHandler) GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	// get order_id from Query parameter
	orderId, err := utils.GetID(r, "order_id")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// get Order item response from service layer
	orderDetailsResponse, err := orderHandler.OrderService.GetOrderDetails(utils.ItoS(orderId))
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// write json response
	err = utils.WriteJSON(w, 200, orderDetailsResponse)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}
func (orderHandler *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// get order_id from Query parameter
	orderId, err := utils.GetID(r, "order_id")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// get Order item response from service layer
	orderResponse, err := orderHandler.OrderService.GetOrder(utils.ItoS(orderId))
	if err != nil {
		_ = utils.ErrorJSON(w, errors.New("order not found"), 404)
		return
	}
	// write json response
	err = utils.WriteJSON(w, 200, orderResponse)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}
func (orderHandler *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {

	// get order_id from Query parameter
	orderId, err := utils.GetID(r, "order_id")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// get status from Query parameter
	status, err := utils.GetString(r, "status")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// validate status
	if !model.IsValidStatus(status) {
		_ = utils.ErrorJSON(w, errors.New("invalid status: "+status), http.StatusBadRequest)
		return
	}
	err = orderHandler.OrderService.ChangeOrderStatus(orderId, status)
	if err != nil {

		_ = utils.ErrorJSON(w, err, http.StatusNotModified)
		return
	}

	// write json response
	err = utils.WriteJSON(w, http.StatusOK, "Order status updated successfully.")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}
func (orderHandler *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	// get order_id from Query parameter
	orderId, err := utils.GetID(r, "order_id")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	//Read order request from json
	var orderRequest *model.OrderRequest
	err = utils.ReadJSON(w, r, &orderRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// get Order response from service layer
	err = orderHandler.OrderService.UpdateOrder(orderId, orderRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// write json response
	err = utils.WriteJSON(w, 201, "Order updated successfully.")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}
func (orderHandler *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// get order_id from Query parameter
	orderId, err := utils.GetID(r, "order_id")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// get Order item response from service layer
	err = orderHandler.OrderService.DeleteOrder(utils.ItoS(orderId))
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// write json response
	err = utils.WriteJSON(w, 204, "Order deleted successfully.")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}
