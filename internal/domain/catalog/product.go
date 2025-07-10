package catalog

import (
	"fmt"
	"time"
)

// ProductID represents a unique product identifier
type ProductID int64

// Price represents a product price with proper precision
type Price float64

// Product represents a product in the catalog domain
// This is the core domain entity following DDD principles
type Product struct {
	ID         ProductID `json:"id"`
	Name       string    `json:"name"`
	Price      Price     `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	SalesCount int       `json:"sales_count"`
	ViewsCount int       `json:"views_count"`
}

// ProductValidationError represents validation errors for products
type ProductValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ProductValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s' with value '%v': %s", 
		e.Field, e.Value, e.Message)
}

// Business logic methods

// SalesConversionRatio calculates the sales-to-views conversion ratio
// This is a key business metric for product performance analysis
func (p Product) SalesConversionRatio() float64 {
	if p.ViewsCount == 0 {
		return 0.0
	}
	return float64(p.SalesCount) / float64(p.ViewsCount)
}

// IsHighPerformer determines if a product is a high performer
// Business rule: conversion ratio > 5% and sales > 50
func (p Product) IsHighPerformer() bool {
	return p.SalesConversionRatio() > 0.05 && p.SalesCount > 50
}

// DaysOnMarket calculates how many days the product has been available
func (p Product) DaysOnMarket() int {
	return int(time.Since(p.CreatedAt).Hours() / 24)
}

// RevenueGenerated calculates total revenue from this product
func (p Product) RevenueGenerated() float64 {
	return float64(p.Price) * float64(p.SalesCount)
}

// Validation methods

// Validate performs comprehensive validation of the product
// Returns detailed validation errors for better debugging
func (p Product) Validate() error {
	var validationErrors []error

	// Validate ID
	if p.ID <= 0 {
		validationErrors = append(validationErrors, ProductValidationError{
			Field:   "ID",
			Value:   p.ID,
			Message: "must be positive",
		})
	}

	// Validate Name
	if p.Name == "" {
		validationErrors = append(validationErrors, ProductValidationError{
			Field:   "Name",
			Value:   p.Name,
			Message: "cannot be empty",
		})
	}
	if len(p.Name) > 255 {
		validationErrors = append(validationErrors, ProductValidationError{
			Field:   "Name",
			Value:   len(p.Name),
			Message: "cannot exceed 255 characters",
		})
	}

	// Validate Price
	if p.Price < 0 {
		validationErrors = append(validationErrors, ProductValidationError{
			Field:   "Price",
			Value:   p.Price,
			Message: "cannot be negative",
		})
	}
	if p.Price > 999999.99 {
		validationErrors = append(validationErrors, ProductValidationError{
			Field:   "Price",
			Value:   p.Price,
			Message: "exceeds maximum allowed value (999999.99)",
		})
	}

	// Validate CreatedAt
	if p.CreatedAt.IsZero() {
		validationErrors = append(validationErrors, ProductValidationError{
			Field:   "CreatedAt",
			Value:   p.CreatedAt,
			Message: "must be set",
		})
	}
	if p.CreatedAt.After(time.Now()) {
		validationErrors = append(validationErrors, ProductValidationError{
			Field:   "CreatedAt",
			Value:   p.CreatedAt,
			Message: "cannot be in the future",
		})
	}

	// Validate SalesCount
	if p.SalesCount < 0 {
		validationErrors = append(validationErrors, ProductValidationError{
			Field:   "SalesCount",
			Value:   p.SalesCount,
			Message: "cannot be negative",
		})
	}

	// Validate ViewsCount
	if p.ViewsCount < 0 {
		validationErrors = append(validationErrors, ProductValidationError{
			Field:   "ViewsCount",
			Value:   p.ViewsCount,
			Message: "cannot be negative",
		})
	}

	// Business rule validation: Sales cannot exceed views
	if p.SalesCount > p.ViewsCount {
		validationErrors = append(validationErrors, ProductValidationError{
			Field:   "SalesCount",
			Value:   fmt.Sprintf("sales: %d, views: %d", p.SalesCount, p.ViewsCount),
			Message: "sales count cannot exceed views count",
		})
	}

	// Return combined validation errors
	if len(validationErrors) > 0 {
		return fmt.Errorf("product validation failed with %d errors: %v", 
			len(validationErrors), validationErrors)
	}

	return nil
}

// IsValid performs quick validation check
func (p Product) IsValid() bool {
	return p.Validate() == nil
}

// String provides a detailed string representation for debugging
func (p Product) String() string {
	return fmt.Sprintf("Product{ID: %d, Name: %q, Price: %.2f, Sales: %d, Views: %d, Ratio: %.4f, Created: %s}",
		p.ID, p.Name, p.Price, p.SalesCount, p.ViewsCount, 
		p.SalesConversionRatio(), p.CreatedAt.Format("2006-01-02"))
}

// ProductCollection represents a collection of products with utility methods
type ProductCollection []Product

// Validate validates all products in the collection
func (pc ProductCollection) Validate() error {
	var validationErrors []error

	for i, product := range pc {
		if err := product.Validate(); err != nil {
			validationErrors = append(validationErrors, 
				fmt.Errorf("product at index %d (ID: %d): %w", i, product.ID, err))
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("collection validation failed with %d errors: %v", 
			len(validationErrors), validationErrors)
	}

	return nil
}

// Copy creates a deep copy of the product collection
// Ensures immutability in sorting operations
func (pc ProductCollection) Copy() ProductCollection {
	if pc == nil {
		return nil
	}
	
	copied := make(ProductCollection, len(pc))
	copy(copied, pc)
	return copied
}

// FilterHighPerformers returns only high-performing products
func (pc ProductCollection) FilterHighPerformers() ProductCollection {
	var highPerformers ProductCollection
	for _, product := range pc {
		if product.IsHighPerformer() {
			highPerformers = append(highPerformers, product)
		}
	}
	return highPerformers
}

// TotalRevenue calculates total revenue for all products in collection
func (pc ProductCollection) TotalRevenue() float64 {
	var total float64
	for _, product := range pc {
		total += product.RevenueGenerated()
	}
	return total
}

// AverageConversionRatio calculates average conversion ratio
func (pc ProductCollection) AverageConversionRatio() float64 {
	if len(pc) == 0 {
		return 0.0
	}

	var total float64
	for _, product := range pc {
		total += product.SalesConversionRatio()
	}
	return total / float64(len(pc))
}

// Len returns the length of the collection (for sort.Interface)
func (pc ProductCollection) Len() int {
	return len(pc)
}

// Less compares products by ID (default comparison for sort.Interface)
func (pc ProductCollection) Less(i, j int) bool {
	return pc[i].ID < pc[j].ID
}

// Swap swaps two products in the collection (for sort.Interface)
func (pc ProductCollection) Swap(i, j int) {
	pc[i], pc[j] = pc[j], pc[i]
}

// Value objects for type safety

// String returns the string representation of ProductID
func (id ProductID) String() string {
	return fmt.Sprintf("ProductID(%d)", int64(id))
}

// IsValid checks if ProductID is valid
func (id ProductID) IsValid() bool {
	return id > 0
}

// String returns the string representation of Price
func (p Price) String() string {
	return fmt.Sprintf("$%.2f", float64(p))
}

// IsValid checks if Price is valid
func (p Price) IsValid() bool {
	return p >= 0 && p <= 999999.99
}

// ToFloat64 converts Price to float64
func (p Price) ToFloat64() float64 {
	return float64(p)
}
