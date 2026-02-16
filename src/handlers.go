package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProduct handles GET /products/:productId
// Logic flow:
//   1. Parse productId from URL path
//   2. Validate it's a positive integer
//   3. Look up in store
//   4. Return 200 with product, or 404 if not found
func GetProduct(store *ProductStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Extract path parameter
		idStr := c.Param("productId")

		// Step 2: Parse to integer
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIError{
				Error:   "INVALID_INPUT",
				Message: "Product ID must be a valid integer",
				Details: "Received: " + idStr,
			})
			return
		}

		// Step 3: Validate range (spec says minimum: 1)
		if id < 1 {
			c.JSON(http.StatusBadRequest, APIError{
				Error:   "INVALID_INPUT",
				Message: "Product ID must be a positive integer",
				Details: "product_id must be >= 1",
			})
			return
		}

		// Step 4: Look up in store
		product, found := store.Get(id)
		if !found {
			c.JSON(http.StatusNotFound, APIError{
				Error:   "NOT_FOUND",
				Message: "Product not found",
				Details: "No product with ID: " + idStr,
			})
			return
		}

		// Step 5: Return product
		c.JSON(http.StatusOK, product)
	}
}

// AddProductDetails handles POST /products/:productId/details
// Logic flow:
//   1. Parse and validate productId from URL path
//   2. Decode and validate JSON body
//   3. Check product exists in store
//   4. Update product
//   5. Return 204 (no body)
func AddProductDetails(store *ProductStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Extract and parse path parameter
		idStr := c.Param("productId")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIError{
				Error:   "INVALID_INPUT",
				Message: "Product ID must be a valid integer",
				Details: "Received: " + idStr,
			})
			return
		}

		if id < 1 {
			c.JSON(http.StatusBadRequest, APIError{
				Error:   "INVALID_INPUT",
				Message: "Product ID must be a positive integer",
				Details: "product_id must be >= 1",
			})
			return
		}

		// Step 2: Decode JSON body
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, APIError{
				Error:   "INVALID_INPUT",
				Message: "Invalid JSON in request body",
				Details: err.Error(),
			})
			return
		}

		// Step 3: Validate all fields against spec constraints
		if detail, valid := ValidateProduct(product); !valid {
			c.JSON(http.StatusBadRequest, APIError{
				Error:   "INVALID_INPUT",
				Message: "The provided input data is invalid",
				Details: detail,
			})
			return
		}

		// Step 4: Update in store (checks existence internally)
		if !store.Update(id, product) {
			c.JSON(http.StatusNotFound, APIError{
				Error:   "NOT_FOUND",
				Message: "Product not found",
				Details: "No product with ID: " + idStr,
			})
			return
		}

		// Step 5: Return 204 No Content — MUST NOT have a response body
		// Use c.Status(), NOT c.JSON() — c.JSON(204, ...) would add a body
		c.Status(http.StatusNoContent)
	}
}