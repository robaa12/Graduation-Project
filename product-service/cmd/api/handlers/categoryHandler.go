package handlers

import (
	"net/http"

	apperrors "github.com/robaa12/product-service/cmd/errors"
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/service"
	"github.com/robaa12/product-service/cmd/utils"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (ch *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	var categoryRequest model.CategoryRequest
	if err := utils.ReadJSON(w, r, &categoryRequest); err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid request payload"))
		return
	}

	categoryResponse, err := ch.service.CreateCategory(storeID, &categoryRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusCreated, categoryResponse)
}

// GetCategories - GET /stores/{store_id}/categories/
func (ch *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {

	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	categories, err := ch.service.GetCategories(storeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, categories)
}

// GetCategoryByID - GET /stores/{store_id}/categories/{category_id}
func (ch *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {

	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	categoryID, err := utils.GetID(r, "category_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid category_id"))
		return
	}
	category, err := ch.service.GetCategoryByID(storeID, categoryID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, category)
}

// GetCategoryBySlug - GET /stores/{store_id}/categories/{category_slug}
func (ch *CategoryHandler) GetCategoryBySlug(w http.ResponseWriter, r *http.Request) {

	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	categorySlug, err := utils.GetString(r, "category_slug")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid category_slug"))
		return
	}
	category, err := ch.service.GetCategoryBySlug(storeID, categorySlug)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, category)
}

// UpdateCategory Update Category - POST /stores/{store_id}/categories/{category_id}
func (ch *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	categoryID, err := utils.GetID(r, "category_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid category_id"))
		return
	}

	var categoryRequest model.CategoryRequest
	if err := utils.ReadJSON(w, r, &categoryRequest); err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid request payload"))
		return
	}
	err = ch.service.UpdateCategory(storeID, categoryID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Category Updated successfully",
	})

}

// DeleteCategory Delete category  - DELETE /stores/{store_id}/categories/{category_id}
func (ch *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	categoryID, err := utils.GetID(r, "category_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid category_id"))
		return
	}

	err = ch.service.DeleteCategory(storeID, categoryID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Category Deleted successfully",
	})
}
