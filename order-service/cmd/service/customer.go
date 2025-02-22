package service

import (
	"errors"
	"order-service/cmd/model"
	"order-service/cmd/repository"
)

type CustomerService struct {
	CustomerRepo *repository.CustomerRepository
}

func NewCustomerService(r *repository.CustomerRepository) *CustomerService {
	return &CustomerService{CustomerRepo: r}
}
func (customerService *CustomerService) CreateNewCustomer(customerRequest *model.CustomerRequest, storeID uint) (*model.CustomerResponseInfo, error) {
	//TODO: Verfiy Store ID Before Create New Customer From My Store Table If Not Existing Go To User Services

	customer := customerRequest.CreateCustomer()
	err := customerService.CustomerRepo.CreateCustomer(customer, storeID)
	if err != nil {
		return nil, err
	}
	customerResponse := customer.CreateCustomerResponseInfo()
	return customerResponse, nil
}
func (customerService *CustomerService) GetAllCustomers(storeID uint) ([]model.StoreCustomerItem, error) {

	customers, err := customerService.CustomerRepo.GetStoreCustomers(storeID)
	if err != nil {
		return nil, err
	}
	return customers, nil
}
func (customerService *CustomerService) GetCustomer(storeID, customerID uint) (*model.CustomerResponseDetails, error) {
	//TODO: Verfiy Store ID Before Create New Customer From My Store Table If not Existing ? Ask Robaa
	storeCustomer := model.CreateStoreCustmer(storeID, customerID)
	//  Validate the customer's relationship with the store
	rowAffected, err := customerService.CustomerRepo.FindStoreCustomer(storeCustomer)
	if err != nil {
		return nil, err
	} else if rowAffected == 0 {
		return nil, errors.New("customer not found")
	}

	//  Fetch the customer and their orders for the specified store
	customer, err := customerService.CustomerRepo.GetStoreCustomerWithOrders(storeCustomer)
	if err != nil {
		return nil, err
	}
	customerResponse := customer.CreateCustomerResponseDetails()
	return customerResponse, nil
}
func (customerService *CustomerService) UpdateCustomer() {
	//TODO: Update Customer Info
}
func (customerService *CustomerService) DeleteCustomer(storeID, customerID uint) error {
	storeCustomer := model.CreateStoreCustmer(storeID, customerID)
	//  Validate the customer's relationship with the store
	rowAffected, err := customerService.CustomerRepo.FindStoreCustomer(storeCustomer)
	if err != nil {
		return err
	} else if rowAffected == 0 {
		return errors.New("customer not found")
	}
	//  Delete the customer and their orders for the specified store
	err = customerService.CustomerRepo.DeleteStoreCustomer(storeCustomer)
	if err != nil {
		return err
	}
	return nil

}
