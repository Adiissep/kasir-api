package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// HandleCategories handles requests for GET|POST /api/categories
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllCategories(w, r)
	case http.MethodPost:
		h.createCategory(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) getAllCategories(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAllCategories()
	if err != nil {
		http.Error(w, "Failed to get all categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

func (h *CategoryHandler) createCategory(w http.ResponseWriter, r *http.Request) {
	var c models.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateCategory(&c); err != nil {
		http.Error(w, "Failed to create category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(c)
}

// HandleCategoryByID handles requests for GET|PUT|DELETE /api/categories/{id}
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getCategoryByID(w, r)
	case http.MethodPut:
		h.updateCategory(w, r)
	case http.MethodDelete:
		h.deleteCategory(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		http.Error(w, "Category not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	type UpdateReq struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	var req UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	old, err := h.service.GetCategoryByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if req.Name != nil {
		old.Name = *req.Name
	}
	if req.Description != nil {
		old.Description = *req.Description
	}

	if err := h.service.UpdateCategory(old); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(old)

	var updatedCategory models.Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedCategory.ID = id
	if err := h.service.UpdateCategory(&updatedCategory); err != nil {
		http.Error(w, "Failed to update category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updatedCategory)
}

func (h *CategoryHandler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCategory(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "Category deleted successfully",
	})

}
