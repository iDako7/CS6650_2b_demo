package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create the in-memory store (with seed data)
	store := NewProductStore()

	// Create the Gin router
	router := gin.Default()

	// Register routes — these must match the OpenAPI spec paths
	// Gin uses :param syntax for path parameters (not {param})
	router.GET("/products/:productId", GetProduct(store))
	router.POST("/products/:productId/details", AddProductDetails(store))

	// Start server on port 8080
	log.Println("Starting Product API server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}