package model

type Service struct {
	Name string `json:"name" validate:"required"`
	URL  string `json:"url" validate:"required,url"`
}

type ServiceCreateStoreRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name" validate:"required"`
}

// ServiceResult tracks the result of operations on individual services
type ServiceResponse struct {
	Success bool                      `json:"success"`
	Error   string                    `json:"error,omitempty"`
	Data    ServiceCreateStoreRequest `json:"store,omitempty"`
}

func ToServiceResponse(success bool, errorMsg string, data ServiceCreateStoreRequest) ServiceResponse {
	return ServiceResponse{
		Success: success,
		Error:   errorMsg,
		Data:    data,
	}
}
