package handlers

import (
	"net/http"
	"order-service/cmd/model"
	"order-service/cmd/service"
	"order-service/cmd/utils"
)

type CustomerHandler struct {
	CustomerService *service.CustomerService
}

func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{CustomerService: customerService}
}
func (customerHandler *CustomerHandler) CreateNewCustomer(w http.ResponseWriter, r *http.Request) {
	// get store_id  from Query parmeter
	storeId, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	// Read Customer Request
	var customerRequest model.CustomerRequest
	err = utils.ReadJSON(w, r, &customerRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	// Get Customer Response
	customerResponse, err := customerHandler.CustomerService.CreateNewCustomer(&customerRequest, storeId)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusConflict)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, customerResponse)
}
func (customerHandler *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	// get store_id  from Query parmeter
	storeId, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	// Get All Customers Response
	customersResponse, err := customerHandler.CustomerService.GetAllCustomers(storeId)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, customersResponse)

}
func (customerHandler *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	// get store_id , customer_id from Query parmeter
	storeId, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	customerId, err := utils.GetID(r, "customer_id")
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	// Get Customer Response
	customerResponse, err := customerHandler.CustomerService.GetCustomer(storeId, customerId)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, customerResponse)
}

// TODO:
func (customerHandler *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {}
func (customerHandler *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	// get store_id , customer_id from Query parmeter
	storeId, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	customerId, err := utils.GetID(r, "customer_id")
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	// Delete Customer
	err = customerHandler.CustomerService.DeleteCustomer(storeId, customerId)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
		return
	}
	// write json resonpse
	utils.WriteJSON(w, http.StatusNoContent, "Order deleted successfully.")

}
