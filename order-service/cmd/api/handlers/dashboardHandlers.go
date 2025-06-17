package handlers

import (
	"errors"
	"log"
	"net/http"
	"order-service/cmd/service"
	"order-service/cmd/utils"
	"time"
)

type DashBoardHandler struct {
	DashBoardService *service.OrderDashBoardService
}

func NewDashBoardHandler(dashBoardService *service.OrderDashBoardService) *DashBoardHandler {
	return &DashBoardHandler{DashBoardService: dashBoardService}
}
func (h *DashBoardHandler) GetDashboardInfo(w http.ResponseWriter, r *http.Request) {
	storeId, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, errors.New("store ID is required"))
		return
	}

	var startDate, endDate time.Time
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	if startDateStr != "" || endDateStr != "" {
		const layout = "2006-01-02" // Correct date format

		if startDateStr != "" {
			startDate, err = time.Parse(layout, startDateStr)
			if err != nil {
				_ = utils.ErrorJSON(w, errors.New("invalid start date format"))
				return
			}
		}

		if endDateStr != "" {
			endDate, err = time.Parse(layout, endDateStr)
			if err != nil {
				_ = utils.ErrorJSON(w, errors.New("invalid end date format"))
				return
			}
		}

		if !startDate.IsZero() && !endDate.IsZero() && startDate.After(endDate) {
			_ = utils.ErrorJSON(w, errors.New("start date must be before end date"))
			return
		}
	}

	dashboardInfo, err := h.DashBoardService.GetDashboardInfo(storeId, startDate, endDate)
	if err != nil {
		log.Println("Error getting dashboard info:", err)
		_ = utils.ErrorJSON(w, errors.New("failed to get dashboard info"))
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, dashboardInfo)
}
