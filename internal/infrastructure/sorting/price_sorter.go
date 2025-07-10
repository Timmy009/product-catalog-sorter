package sorting

import (
	"context"
	"sort"

	"product-catalog-sorting/internal/domain/catalog"
)

// PriceSorter sorts products by price
type PriceSorter struct {
	ascending bool
}

// NewPriceSorter creates a new price sorter
func NewPriceSorter(ascending bool) catalog.Sorter {
	return &PriceSorter{
		ascending: ascending,
	}
}

// Sort implements the Sorter interface
func (s *PriceSorter) Sort(ctx context.Context, products catalog.ProductCollection) (catalog.ProductCollection, error) {
	if len(products) == 0 {
		return catalog.ProductCollection{}, nil
	}

	// Create a copy to avoid mutating the original
	sorted := products.Copy()

	// Sort using Go's built-in sort package
	sort.Slice(sorted, func(i, j int) bool {
		if s.ascending {
			return sorted[i].Price < sorted[j].Price
		}
		return sorted[i].Price > sorted[j].Price
	})

	return sorted, nil
}

// GetStrategy returns the sort strategy
func (s *PriceSorter) GetStrategy() catalog.SortStrategy {
	if s.ascending {
		return catalog.SortByPriceAsc
	}
	return catalog.SortByPriceDesc
}

// GetDescription returns a human-readable description
func (s *PriceSorter) GetDescription() string {
	if s.ascending {
		return "Sorts products by price from lowest to highest"
	}
	return "Sorts products by price from highest to lowest"
}
