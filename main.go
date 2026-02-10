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
	// Category
	categoryRepo := repositories.NewCategoryRepository(pool)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	// Transaction
	transactionRepo := repositories.NewTransactionRepository(pool)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	// Report
	reportRepo := repositories.NewReportRepository(pool)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	// Setup routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)
	//post /api/checkout
	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)
	http.HandleFunc("/api/report/today", reportHandler.HandleReportToday)

	//localhost:8080/api
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"status":  "Success",
			"message": "Welcome to the Cashier API",
			"endpoints": []string{
				"GET /api/products",
				"POST /api/products",
				"GET /api/product/{id}",
				"PUT /api/product/{id}",
				"DELETE /api/product/{id}",
				"GET /api/products?name={name}",

				"GET /api/categories",
				"POST /api/categories",
				"GET /api/category/{id}",
				"PUT /api/category/{id}",
				"DELETE /api/category/{id}",

				"POST /api/checkout",
				"GET /api/report/today",
				"Comming Soon GET /api/report?date={date}",
			},
		}); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Bind ke semua interface (IPv4/IPv6)
	addr := ":" + config.Port
	fmt.Println("Server running in", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
