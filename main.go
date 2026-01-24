package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "Apple", Price: 100, Stock: 50},
	{ID: 2, Name: "Banana", Price: 50, Stock: 100},
	{ID: 3, Name: "Tomato", Price: 80, Stock: 75},
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Fruits", Description: "Fresh fruits"},
	{ID: 2, Name: "Vegetables", Description: "Fresh vegetables"},
}

// GET localhost:8080/api/products/{id}
func getProductByID(w http.ResponseWriter, r *http.Request) {
	// GET product ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")

	// Convert ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Find product by ID
	for _, p := range products {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(p); err != nil {
				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// PUT localhost:8080/api/products/{id}
func updateProductByID(w http.ResponseWriter, r *http.Request) {
	// GET product ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")

	// Convert ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Find product by ID
	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Invalid request body"+err.Error(), http.StatusBadRequest)
		return
	}

	// loop through products to find and update the product
	for i := range products {
		if products[i].ID == id {
			updatedProduct.ID = id
			products[i] = updatedProduct

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(products[i]); err != nil {
				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}

// DELETE localhost:8080/api/products/{id}
func deleteProductByID(w http.ResponseWriter, r *http.Request) {
	// GET product ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")

	// Convert ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	//loop through products to find and delete the product
	for i, p := range products {
		if p.ID == id {
			//create new slice with before and after index
			products = append(products[:i], products[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{
				"message": "Product deleted successfully",
			}); err != nil {
				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}

// GET localhost:8080/api/categories/{id}
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(c); err != nil {
				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// PUT localhost:8080/api/categories/{id}
func updateCategoryByID(w http.ResponseWriter, r *http.Request) {
	// GET category ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// Convert ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Find category by ID
	var updatedCategory Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "Invalid request body"+err.Error(), http.StatusBadRequest)
		return
	}

	// loop through categories to find and update the category
	for i := range categories {
		if categories[i].ID == id {
			updatedCategory.ID = id
			categories[i] = updatedCategory

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(categories[i]); err != nil {
				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
}

// DELETE localhost:8080/api/categories/{id}
func deleteCategoryByID(w http.ResponseWriter, r *http.Request) {
	// GET product ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// Convert ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	//loop through categories to find and delete the category
	for i, c := range categories {
		if c.ID == id {
			//create new slice with before and after index
			categories = append(categories[:i], categories[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{
				"message": "Category deleted successfully",
			}); err != nil {
				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
}

func main() {
	//GET & POST localhost:8080/api/categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(categories); err != nil {
				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
				return
			}
		case "POST":
			// read data from request body
			var newCategory Category
			err := json.NewDecoder(r.Body).Decode(&newCategory)
			if err != nil {
				http.Error(w, "Invalid request"+err.Error(), http.StatusBadRequest)
				return
			}

			// add new product to products slice
			newCategory.ID = len(categories) + 1
			categories = append(categories, newCategory)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(newCategory); err != nil {
				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	//GET & PUT localhost:8080/api/products/{id}
	//DELETE localhost:8080/api/products/{id}
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			// Implement getCategoryByID if needed
			getCategoryByID(w, r)
		case "PUT":
			// Implement updateCategoryByID if needed
			updateCategoryByID(w, r)
		case "DELETE":
			// Implement deleteCategoryByID if needed
			deleteCategoryByID(w, r)
		}
	})

	//GET & POST localhost:8080/api/products
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(products); err != nil {
				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
				return
			}
		case "POST":
			// read data from request body
			var newProduct Product
			err := json.NewDecoder(r.Body).Decode(&newProduct)
			if err != nil {
				http.Error(w, "Invalid request"+err.Error(), http.StatusBadRequest)
				return
			}

			// add new product to products slice
			newProduct.ID = len(products) + 1
			products = append(products, newProduct)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(newProduct); err != nil {
				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	//GET & PUT localhost:8080/api/products/{id}
	//DELETE localhost:8080/api/products/{id}
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getProductByID(w, r)
		case "PUT":
			updateProductByID(w, r)
		case "DELETE":
			deleteProductByID(w, r)
		}
	})

	//localhost:8080/api
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"status":  "Success",
			"message": "Welcome to the Cashier API",
			"endpoints": []string{
				"GET /api/categories",
				"POST /api/categories",
				"GET /api/categories/{id}",
				"PUT /api/categories/{id}",
				"DELETE /api/categories/{id}",
				"GET /api/products",
				"POST /api/products",
				"GET /api/products/{id}",
				"PUT /api/products/{id}",
				"DELETE /api/products/{id}",
			},
		}); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server running di localhost:8080 🚀")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
