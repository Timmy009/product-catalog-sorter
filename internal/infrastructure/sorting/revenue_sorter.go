package sorting

import (
	"context"
	"sort"

	"product-catalog-sorting/internal/domain/catalog"
)

// RevenueSorter sorts products by revenue generated
type RevenueSorter struct{}

// NewRevenueSorter creates a new revenue sorter
func NewRevenueSorter() catalog.Sorter {
	return &RevenueSorter{}
}

// Sort implements the Sorter interface
func (s *RevenueSorter) Sort(ctx context.Context, products catalog.ProductCollection) (catalog.ProductCollection, error) {
	if len(products) == 0 {
		return catalog.ProductCollection{}, nil
	}

	// Create a copy to avoid mutating the original
	sorted := products.Copy()

	// Sort by revenue generated (descending)
	sort.Slice(sorted, func(i, j int) bool {
		revenueI := sorted[i].RevenueGenerated()
		revenueJ := sorted[j].RevenueGenerated()

		// Primary sort: revenue (higher is better)
		if revenueI != revenueJ {
			return revenueI > revenueJ
		}

		// Secondary sort: sales count (higher is better)
		if sorted[i].SalesCount != sorted[j].SalesCount {
			return sorted[i].SalesCount > sorted[j].SalesCount
		}

		// Tertiary sort: ID for consistent ordering
		return sorted[i].ID < sorted[j].ID
	})

	return sorted, nil
}

// GetStrategy returns the sort strategy
func (s *RevenueSorter) GetStrategy() catalog.SortStrategy {
	return catalog.SortByRevenue
}

// GetDescription returns a human-readable description
func (s *RevenueSorter) GetDescription() string {
	return "Sorts products by revenue generated from highest to lowest"
}
