package handlers

import (
	"errors"
	"net/http"
	"order-service/cmd/model"
	"order-service/cmd/service"
	"order-service/cmd/utils"
)

type OrderItemsHandler struct {
	OrderItemService *service.OrderItemService
}

func NewOrderItemsHandler(orderItemService *service.OrderItemService) *OrderItemsHandler {
	return &OrderItemsHandler{OrderItemService: orderItemService}
}
func (orderItemsHandler *OrderItemsHandler) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	// get order_id from Query parameter
	orderId, err := utils.GetID(r, "order_id")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// Read OrderItem Request from json
	var orderItemRequest model.OrderItemRequest
	err = utils.ReadJSON(w, r, &orderItemRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, errors.New("enter valid order item data"))
		return
	}

	// give order item response from service layer
	orderItemResponse, err := orderItemsHandler.OrderItemService.AddOrderItem(orderId, &orderItemRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// write json response
	err = utils.WriteJSON(w, 201, orderItemResponse)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}

func (orderItemsHandler *OrderItemsHandler) GetAllOrderItems(w http.ResponseWriter, r *http.Request) {
	// get order_id from Query parameter
	orderId, err := utils.GetID(r, "order_id")
	if err != nil {

		_ = utils.ErrorJSON(w, err)
		return
	}
	// get Order item response from service layer
	orderItemsResponse, err := orderItemsHandler.OrderItemService.GetAllOrderItems(utils.ItoS(orderId))
	if err != nil {

		_ = utils.ErrorJSON(w, err)
		return
	}

	// write json response
	err = utils.WriteJSON(w, 200, orderItemsResponse)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}
func (orderItemsHandler *OrderItemsHandler) GetOrderItem(w http.ResponseWriter, r *http.Request) {
	// get order item_id from Query parameter
	orderItemId, err := utils.GetID(r, "item_id")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// get Order item response from service layer
	orderItemResponse, err := orderItemsHandler.OrderItemService.GetOrderItem(utils.ItoS(orderItemId))
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// write json response
	err = utils.WriteJSON(w, 200, orderItemResponse)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}
func (orderItemsHandler *OrderItemsHandler) UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	// get order item ID from Query parameter
	orderItemId, err := utils.GetID(r, "item_id")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	//Read order item request from json
	var orderItemRequest *model.OrderItemRequest
	err = utils.ReadJSON(w, r, &orderItemRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// get Order item response from service layer
	err = orderItemsHandler.OrderItemService.UpdateOrderItem(utils.ItoS(orderItemId), orderItemRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// write json response
	err = utils.WriteJSON(w, 200, "Order item updated successfully.")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}
func (orderItemsHandler *OrderItemsHandler) DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	// get order item_id from Query parameter
	orderItemId, err := utils.GetID(r, "item_id")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// get Order item response from service layer
	err = orderItemsHandler.OrderItemService.DeleteOrderItem(utils.ItoS(orderItemId))
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// write json response
	err = utils.WriteJSON(w, 200, "Order item deleted successfully.")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

}
