package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// HandleProducts - GET /api/products|POST /api/products
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Println("Filtering products by name:", name)
	products, err := h.service.GetAll(name)
	if err != nil {
		http.Error(w, "Failed to get all products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode products: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&newProduct); err != nil {
		http.Error(w, "Failed to create product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newProduct)
}

// HandleProductByID - GET|PUT|DEL /api/product/{id}
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// perbaiki prefix: gunakan /api/product/
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(strings.ToLower(err.Error()), "not found") {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get product: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Baca body sekali
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Struct untuk field lain (pointer → partial update)
	type UpdateReq struct {
		Name  *string `json:"name"`
		Price *int    `json:"price"`
		Stock *int    `json:"stock"`
		// CategoryID akan diproses manual untuk bedakan "missing" vs "null"
	}

	var req UpdateReq
	if len(body) > 0 {
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Deteksi presence key "category_id" dan nil vs value
	var raw map[string]*json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	categoryPresent := false
	var categoryIDDecoded *int // jika present & bukan null → pointer ke int, jika present & null → tetap nil
	if raw != nil {
		if rm, ok := raw["category_id"]; ok {
			categoryPresent = true
			if rm != nil {
				var v int
				if err := json.Unmarshal(*rm, &v); err != nil {
					http.Error(w, "Invalid category_id (must be number or null)", http.StatusBadRequest)
					return
				}
				categoryIDDecoded = &v
			} else {
				// rm == nil → "category_id": null
				categoryIDDecoded = nil
			}
		}
	}

	old, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(strings.ToLower(err.Error()), "not found") {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get product: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if req.Name != nil {
		old.Name = *req.Name
	}
	if req.Price != nil {
		old.Price = *req.Price
	}
	if req.Stock != nil {
		old.Stock = *req.Stock
	}

	// category_id: hanya ubah kalau key hadir
	if categoryPresent {
		// present + null → unset (set NULL di DB)
		// present + number → set new value
		old.CategoryID = categoryIDDecoded
	}

	if err := h.service.Update(old); err != nil {
		http.Error(w, "Failed to update product: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Re-fetch setelah update agar category_name hasil JOIN ikut terbarui
	updated, err := h.service.GetByID(id)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(updated)
		return
	}

	// Fallback: jika re-fetch gagal, set category_name sesuai perubahan category_id
	if categoryPresent && categoryIDDecoded == nil {
		old.CategoryName = ""
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(old)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Product deleted successfully",
	})
}
