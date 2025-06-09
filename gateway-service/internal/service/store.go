package service

import (
	"fmt"
	"log"

	apperrors "github.com/robaa12/gatway-service/internal/errors"

	httpcient "github.com/robaa12/gatway-service/internal/http-cient"
	"github.com/robaa12/gatway-service/internal/model"
)

type StoreService struct {
	Client *httpcient.Client
}

func NewStoreService(client *httpcient.Client) *StoreService {
	return &StoreService{
		Client: client,
	}
}

func (s *StoreService) CreateStore(storeRequest *model.StoreRequest) (*model.Store, error) {

	// Step 1: Create store in user service
	storeUserResponse, err := s.Client.CreateStoreInUserService(storeRequest)

	if err != nil {

		return nil, apperrors.NewInternalServerError(fmt.Sprintf("failed to create store in user service: %v", err))
	}

	// Parse the user service response
	store := storeUserResponse.GetStore()

	// Extract the store ID from the response
	storeServicesRequest := store.ToServiceCreateStoreRequest()

	// Step 2: Create store in product and order services concurrently
	servicesResponse, successfulServices := s.Client.CreateStoreInServices(&storeServicesRequest)

	// Step 3: Check if all services succeeded
	allSucceeded := true
	for _, result := range servicesResponse {
		if !result.Success {
			allSucceeded = false
			break
		}
	}

	// Step 4: Handle compensating transactions if necessary
	if !allSucceeded {
		log.Printf("Store creation failed in some services. Initiating compensating transactions.")
		s.Client.CompensateStoreCreation(successfulServices, store.ID)

		// Respond with error
		return nil, apperrors.NewInternalServerError(fmt.Sprintf("store creation failed in some services: %v", servicesResponse))
	}

	return &store, nil

}
