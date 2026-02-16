package main

import "fmt"

// ValidateProduct checks all field constraints, more related to business logic, from the OpenAPI spec.
// Returns an error message string and false if validation fails.
// Returns empty string and true if valid.
func ValidateProduct(p Product) (string, bool) {
	if p.ProductID < 1 {
		return "product_id must be a positive integer (minimum 1)", false
	}

	if len(p.SKU) == 0 {
		return "sku is required and must not be empty", false
	}
	if len(p.SKU) > 100 {
		return fmt.Sprintf("sku must be at most 100 characters, got %d", len(p.SKU)), false
	}

	if len(p.Manufacturer) == 0 {
		return "manufacturer is required and must not be empty", false
	}
	if len(p.Manufacturer) > 200 {
		return fmt.Sprintf("manufacturer must be at most 200 characters, got %d", len(p.Manufacturer)), false
	}

	if p.CategoryID < 1 {
		return "category_id must be a positive integer (minimum 1)", false
	}

	if p.Weight < 0 {
		return "weight must be a non-negative integer (minimum 0)", false
	}

	if p.SomeOtherID < 1 {
		return "some_other_id must be a positive integer (minimum 1)", false
	}

	return "", true
}