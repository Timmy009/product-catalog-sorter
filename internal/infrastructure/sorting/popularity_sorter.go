package sorting

import (
	"context"
	"sort"

	"product-catalog-sorting/internal/domain/catalog"
)

// PopularitySorter sorts products by view count (popularity)
type PopularitySorter struct{}

// NewPopularitySorter creates a new popularity sorter
func NewPopularitySorter() catalog.Sorter {
	return &PopularitySorter{}
}

// Sort implements the Sorter interface
func (s *PopularitySorter) Sort(ctx context.Context, products catalog.ProductCollection) (catalog.ProductCollection, error) {
	if len(products) == 0 {
		return catalog.ProductCollection{}, nil
	}

	// Create a copy to avoid mutating the original
	sorted := products.Copy()

	// Sort by popularity (views) with tie-breaking logic
	sort.Slice(sorted, func(i, j int) bool {
		// Primary sort: view count (higher is better)
		if sorted[i].ViewsCount != sorted[j].ViewsCount {
			return sorted[i].ViewsCount > sorted[j].ViewsCount
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
func (s *PopularitySorter) GetStrategy() catalog.SortStrategy {
	return catalog.SortByPopularity
}

// GetDescription returns a human-readable description
func (s *PopularitySorter) GetDescription() string {
	return "Sorts products by popularity (view count) from highest to lowest"
}
