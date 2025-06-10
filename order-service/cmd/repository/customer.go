package repository

import (
	"errors"
	"fmt"
	"order-service/cmd/model"

	"gorm.io/gorm"
)

// CustomerRepository handles database operations for customers
type CustomerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new instance of CustomerRepository
func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// AddCustomer inserts a new customer into the database
func AddCustomer(customer *model.Customer, tx *gorm.DB) error {
	return tx.FirstOrCreate(customer, model.Customer{Email: customer.Email}).Error
}

// AddStoreCustomer creates a new relationship between a store and a customer
func AddStoreCustomer(storeCustomer *model.StoreCustomer, tx *gorm.DB) (bool, error) {
	// Use FirstOrCreate to find or create the relationship
	result := tx.FirstOrCreate(storeCustomer, model.StoreCustomer{StoreID: storeCustomer.StoreID, CustomerID: storeCustomer.CustomerID})

	// Handle errors
	if result.Error != nil {
		return false, fmt.Errorf("failed to add store-customer relationship: %w", result.Error)
	}

	// Log whether a new record was created or an existing one was found
	if result.RowsAffected == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

// CreateCustomer inserts a new customer into the database
func (r *CustomerRepository) CreateCustomer(customer *model.Customer, storeID uint) error {
	// start transaction
	tx := r.db.Begin()
	// Defer rollback if transaction fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := AddCustomer(customer, tx); err != nil {
		tx.Rollback()
		return err
	}
	storeCustomer := model.CreateStoreCustmer(storeID, customer.ID)
	isNew, err := AddStoreCustomer(storeCustomer, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	if !isNew {
		return errors.New("customer is already existing in the store")
	}
	return nil
}

// GetStoreCustomerWithOrders retrieves a customer by ID along with their orders
func (r *CustomerRepository) GetStoreCustomerWithOrders(storeCustomer *model.StoreCustomer) (*model.Customer, error) {
	var customer *model.Customer
	err := r.db.Preload("Orders", "store_id = ?", storeCustomer.StoreID).First(&customer, storeCustomer.CustomerID).
		Error
	if err != nil {
		return nil, err
	}

	return customer, nil
}

// GetStoreCustomerByEmail retrieves a customer using their email with their orders
func (r *CustomerRepository) GetStoreCustomerByEmail(customer *model.Customer, storeID uint) error {

	return r.db.Preload("Orders", "store_id = ?", storeID).Where("email = ?", customer.Email).First(&customer).Error

}

// GetStoreCustomers retrieves all customers linked to a specific store
func (r *CustomerRepository) GetStoreCustomers(storeID uint) ([]model.StoreCustomerItem, error) {
	storeCustomerItems := []model.StoreCustomerItem{}

	// Query to retrieve customers and their order statistics for a specific store
	query := r.db.Model(&model.Customer{}).
		Select(`
            customers.id as customer_id, 
            customers.email as customer_email,  
            COUNT(orders.id) as number_of_orders, 
            COALESCE(SUM(orders.total_price), 0) as total_spent
        `).
		Joins("JOIN store_customers ON store_customers.customer_id = customers.id").
		Joins("LEFT JOIN orders ON orders.customer_id = customers.id AND orders.store_id = ?", storeID).
		Where("store_customers.store_id = ?", storeID).
		Group("customers.id")

	// Execute the query and scan the results
	if err := query.Scan(&storeCustomerItems).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve store customers: %w", err)
	}

	return storeCustomerItems, nil
}

// FindStoreCustomer Find Store Customer retrieves rows Affected
func (r *CustomerRepository) FindStoreCustomer(storeCustomer *model.StoreCustomer) (int64, error) {

	result := r.db.Where("store_id = ? AND customer_id = ?", storeCustomer.StoreID, storeCustomer.CustomerID).Find(storeCustomer)
	return result.RowsAffected, result.Error
}

// DeleteStoreCustomer removes a relationship between a store and a customer
func (r *CustomerRepository) DeleteStoreCustomer(storeCustomer *model.StoreCustomer) error {
	// Start a transaction
	tx := r.db.Begin()

	// Defer rollback if transaction fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//  Delete the customer's orders for the specified store
	err := tx.Unscoped().Delete(&model.Order{}, model.Order{StoreID: storeCustomer.StoreID, CustomerID: storeCustomer.CustomerID}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	//Delete the customer's relationship with the store
	err = tx.Unscoped().Delete(storeCustomer, storeCustomer).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// FindCustomer Find Customer By ID retrieves rows Affected
func (r *CustomerRepository) FindCustomer(customer *model.Customer) (int64, error) {

	result := r.db.Find(customer, customer.ID)
	return result.RowsAffected, result.Error
}

// DeleteCustomer removes a customer by ID
func (r *CustomerRepository) DeleteCustomer(customer *model.Customer) error {
	return r.db.Unscoped().Delete(customer, customer.ID).Error
}
