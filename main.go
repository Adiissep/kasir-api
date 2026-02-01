package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// // GET localhost:8080/api/categories/{id}
// func getCategoryByID(w http.ResponseWriter, r *http.Request) {
// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid category ID", http.StatusBadRequest)
// 		return
// 	}

// 	for _, c := range categories {
// 		if c.ID == id {
// 			w.Header().Set("Content-Type", "application/json")
// 			if err := json.NewEncoder(w).Encode(c); err != nil {
// 				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			return
// 		}
// 	}

// 	http.Error(w, "Category not found", http.StatusNotFound)
// }

// // PUT localhost:8080/api/categories/{id}
// func updateCategoryByID(w http.ResponseWriter, r *http.Request) {
// 	// GET category ID from URL
// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

// 	// Convert ID to integer
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid category ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Find category by ID
// 	var updatedCategory Category
// 	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
// 	if err != nil {
// 		http.Error(w, "Invalid request body"+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// loop through categories to find and update the category
// 	for i := range categories {
// 		if categories[i].ID == id {
// 			updatedCategory.ID = id
// 			categories[i] = updatedCategory

// 			w.Header().Set("Content-Type", "application/json")
// 			if err := json.NewEncoder(w).Encode(categories[i]); err != nil {
// 				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			return
// 		}
// 	}
// 	http.Error(w, "Category not found", http.StatusNotFound)
// }

// DELETE localhost:8080/api/categories/{id}
// func deleteCategoryByID(w http.ResponseWriter, r *http.Request) {
// 	// GET product ID from URL
// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

// 	// Convert ID to integer
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid category ID", http.StatusBadRequest)
// 		return
// 	}

// 	//loop through categories to find and delete the category
// 	for i, c := range categories {
// 		if c.ID == id {
// 			//create new slice with before and after index
// 			categories = append(categories[:i], categories[i+1:]...)
// 			w.Header().Set("Content-Type", "application/json")
// 			if err := json.NewEncoder(w).Encode(map[string]string{
// 				"message": "Category deleted successfully",
// 			}); err != nil {
// 				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			return
// 		}
// 	}
// 	http.Error(w, "Category not found", http.StatusNotFound)
// }

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}
	if config.DBConn == "" {
		log.Fatal("DB_CONN is empty. Ensure .env has DB_CONN=<connection string>")
	}
	fmt.Println("Using port:", config.Port)
	//fmt.Println("DB Connection:", config.DBConn)

	// Setup DB (pgxpool)
	pool, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer pool.Close()

	// Injeksi pgxpool ke repository (pastikan constructor repo menerima *pgxpool.Pool)
	productRepo := repositories.NewProductRepository(pool)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)

	//localhost:8080/api
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"status":  "Success",
			"message": "Welcome to the Cashier API",
			"endpoints": []string{
				"GET /api/products",
				"POST /api/products",
				"GET /api/products/{id}",
				"PUT /api/products/{id}",
				"DELETE /api/products/{id}",

				"GET /api/categories",
				"POST /api/categories",
				"GET /api/categories/{id}",
				"PUT /api/categories/{id}",
				"DELETE /api/categories/{id}",
			},
		}); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Bind ke semua interface (IPv4/IPv6)
	addr := ":" + config.Port
	fmt.Println("Server running in:", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
