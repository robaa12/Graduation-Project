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
	// get order_id from Query parmeter
	orderId, err := utils.GetID(r, "order_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// Read OrderItem Request from json
	var orderItemRequest model.OrderItemRequest
	err = utils.ReadJSON(w, r, &orderItemRequest)
	if err != nil {
		utils.ErrorJSON(w, errors.New("enter valid order item data"))
		return
	}

	// give orderitem respone from service layer
	orderItemResponse, err := orderItemsHandler.OrderItemService.AddOrderItem(orderId, &orderItemRequest)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// write json resonpse
	err = utils.WriteJSON(w, 201, orderItemResponse)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

}

func (orderItemsHandler *OrderItemsHandler) GetAllOrderItems(w http.ResponseWriter, r *http.Request) {
	// get order_id from Query parmeter
	order_id, err := utils.GetID(r, "order_id")
	if err != nil {

		utils.ErrorJSON(w, err)
		return
	}
	// get Orderitem response from service layer
	orderItemsResponse, err := orderItemsHandler.OrderItemService.GetAllOrderItems(utils.ItoS(order_id))
	if err != nil {

		utils.ErrorJSON(w, err)
		return
	}

	// write json resonpse
	err = utils.WriteJSON(w, 200, orderItemsResponse)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

}
func (orderItemsHandler *OrderItemsHandler) GetOrderItem(w http.ResponseWriter, r *http.Request) {
	// get orderitem_id from Query parmeter
	orderItemId, err := utils.GetID(r, "item_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// get Orderitem response from service layer
	orderItemResponse, err := orderItemsHandler.OrderItemService.GetOrderItem(utils.ItoS(orderItemId))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// write json resonpse
	err = utils.WriteJSON(w, 200, orderItemResponse)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

}
func (orderItemsHandler *OrderItemsHandler) UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	// get orderitem_id from Query parmeter
	orderItem_id, err := utils.GetID(r, "item_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	//Read orderitem request from json
	var orderItemRequest *model.OrderItemRequest
	err = utils.ReadJSON(w, r, &orderItemRequest)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// get Orderitem response from service layer
	err = orderItemsHandler.OrderItemService.UpdateOrderItem(utils.ItoS(orderItem_id), orderItemRequest)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// write json response
	err = utils.WriteJSON(w, 200, "Order item updated successfully.")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

}
func (orderItemsHandler *OrderItemsHandler) DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	// get orderitem_id from Query parmeter
	orderitem_id, err := utils.GetID(r, "item_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// get Orderitem response from service layer
	err = orderItemsHandler.OrderItemService.DeleteOrderItem(utils.ItoS(orderitem_id))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// write json resonpse
	err = utils.WriteJSON(w, 200, "Order item deleted successfully.")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

}
