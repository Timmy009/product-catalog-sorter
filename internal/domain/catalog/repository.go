package catalog

import (
	"context"
	"time"
)

// Repository defines the contract for product data access
// Following Repository pattern for clean architecture
type Repository interface {
	// GetProducts retrieves products with optional filtering
	GetProducts(ctx context.Context, filter ProductFilter) (ProductCollection, error)

	// GetProductByID retrieves a single product by ID
	GetProductByID(ctx context.Context, id ProductID) (*Product, error)

	// SaveProduct saves or updates a product
	SaveProduct(ctx context.Context, product *Product) error

	// DeleteProduct removes a product
	DeleteProduct(ctx context.Context, id ProductID) error

	// GetProductCount returns the total number of products
	GetProductCount(ctx context.Context, filter ProductFilter) (int, error)
}

// ProductFilter represents filtering criteria for product queries
type ProductFilter struct {
	IDs           []ProductID `json:"ids,omitempty"`
	NameContains  string      `json:"name_contains,omitempty"`
	MinPrice      *Price      `json:"min_price,omitempty"`
	MaxPrice      *Price      `json:"max_price,omitempty"`
	MinSales      *int        `json:"min_sales,omitempty"`
	MaxSales      *int        `json:"max_sales,omitempty"`
	MinViews      *int        `json:"min_views,omitempty"`
	MaxViews      *int        `json:"max_views,omitempty"`
	CreatedAfter  *time.Time  `json:"created_after,omitempty"`
	CreatedBefore *time.Time  `json:"created_before,omitempty"`
	Limit         int         `json:"limit,omitempty"`
	Offset        int         `json:"offset,omitempty"`
}

// IsEmpty returns true if no filters are applied
func (f ProductFilter) IsEmpty() bool {
	return len(f.IDs) == 0 &&
		f.NameContains == "" &&
		f.MinPrice == nil &&
		f.MaxPrice == nil &&
		f.MinSales == nil &&
		f.MaxSales == nil &&
		f.MinViews == nil &&
		f.MaxViews == nil &&
		f.CreatedAfter == nil &&
		f.CreatedBefore == nil
}
