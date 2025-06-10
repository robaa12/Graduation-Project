package httpcient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/robaa12/gatway-service/internal/model"
)

type Client struct {
	userServiceURL    string
	productServiceURL string
	orderServiceURL   string
	client            *http.Client
}

// NewClient creates a new HTTP client for the gateway service
func NewClient(userServiceURL, productServiceURL, orderServiceURL string) *Client {
	return &Client{
		userServiceURL:    userServiceURL,
		productServiceURL: productServiceURL,
		orderServiceURL:   orderServiceURL,
		client:            &http.Client{},
	}
}

// Helper method to create store in user service
func (c *Client) CreateStoreInUserService(storeRequest *model.StoreRequest) (*model.StoreUserResponse, error) {
	requestBody, err := json.Marshal(storeRequest)
	if err != nil {
		return nil, fmt.Errorf("marshaling request body failed: %w", err)
	}
	resp, respBody, err := c.sendRequest(http.MethodPost, c.userServiceURL+"/store", requestBody)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	// Parse the store response to extract the store IDAdd commentMore actions
	var storeResponse model.StoreUserResponse
	if err := json.Unmarshal(respBody, &storeResponse); err != nil {
		log.Printf("Error parsing store response: %v", err)
		// Continue anyway to return the original response
	}

	if resp.StatusCode != http.StatusCreated {
		// Forward the error response from the user service

		return nil, fmt.Errorf("user service returned status code %d", resp.StatusCode)
	}
	return &storeResponse, nil
}
func (c *Client) CreateStoreInServices(storeServicesRequest *model.ServiceCreateStoreRequest) (map[string]model.ServiceResponse, []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	serviceResponses := make(map[string]model.ServiceResponse)
	successfulServices := make([]string, 0)
	services := []model.Service{
		{Name: "product_service", URL: c.productServiceURL},
		{Name: "order_service", URL: c.orderServiceURL},
		// Add more here
	}
	for _, svc := range services {
		c.createStoreInService(&wg, &mu, svc.Name, svc.URL, storeServicesRequest, serviceResponses, &successfulServices)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return serviceResponses, successfulServices

}

// Helper function to create store in a service
func (c *Client) createStoreInService(
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	serviceName string,
	url string,
	serviceRequest *model.ServiceCreateStoreRequest,
	serviceResponses map[string]model.ServiceResponse,
	successfulServices *[]string,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		requestBody, err := json.Marshal(serviceRequest)
		var seviceResponse model.ServiceCreateStoreRequest
		success := false
		errorMsg := ""
		if err != nil {
			errorMsg = "marshaling request body failed: " + err.Error()
		} else {
			resp, respBody, reqErr := c.sendRequest(http.MethodPost, url+"/stores", requestBody)
			if reqErr != nil {
				errorMsg = reqErr.Error()
			} else if resp != nil && resp.StatusCode != http.StatusCreated {
				errorMsg = fmt.Sprintf("%s returned status code %d", serviceName, resp.StatusCode)
			} else {
				success = true
				if respBody != nil {
					_ = json.Unmarshal(respBody, &seviceResponse)
				}
			}
		}
		mu.Lock()
		defer mu.Unlock()
		if success {
			*successfulServices = append(*successfulServices, serviceName)
		}
		serviceResponses[serviceName] = model.ToServiceResponse(success, errorMsg, seviceResponse)
	}()
}

// Compensating transaction: Delete store from all successful services
func (c *Client) CompensateStoreCreation(successfulServices []string, storeID uint) {
	var wg sync.WaitGroup

	for _, serviceName := range successfulServices {
		wg.Add(1)
		go func(service string) {
			defer wg.Done()
			var err error
			switch service {
			case "product_service":
				err = c.deleteStoreFromService(c.productServiceURL+"/stores", storeID)
			case "order_service":
				err = c.deleteStoreFromService(c.orderServiceURL+"/stores", storeID)
			}
			if err != nil {
				log.Printf("Compensation transaction failed for %s: %v", service, err)
			} else {
				log.Printf("Successfully performed compensating transaction for %s", service)
			}
		}(serviceName)
	}

	// Always attempt to delete from user service as well
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := c.deleteStoreFromService(c.userServiceURL+"/store", storeID)
		if err != nil {
			log.Printf("Failed to delete store from user service: %v", err)
		} else {
			log.Printf("Successfully deleted store from user service")
		}
	}()

	wg.Wait()
}

// Compensating transaction: Delete store from any service
func (c *Client) deleteStoreFromService(serviceURL string, storeID uint) error {
	url := fmt.Sprintf("%s/%d", serviceURL, storeID)

	resp, _, err := c.sendRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("delete store from service failed with status code %d", resp.StatusCode)
	}

	return nil
}

// Helper method to send HTTP requests
func (c *Client) sendRequest(method, url string, body []byte) (*http.Response, []byte, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("creating request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return resp, nil, fmt.Errorf("reading response body failed: %w", err)
	}
	resp.Body.Close()

	return resp, respBody, nil
}
