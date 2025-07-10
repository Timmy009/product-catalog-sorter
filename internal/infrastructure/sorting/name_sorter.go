package sorting

import (
	"context"
	"sort"
	"strings"

	"product-catalog-sorting/internal/domain/catalog"
)

// NameSorter sorts products alphabetically by name
type NameSorter struct{}

// NewNameSorter creates a new name sorter
func NewNameSorter() catalog.Sorter {
	return &NameSorter{}
}

// Sort implements the Sorter interface
func (s *NameSorter) Sort(ctx context.Context, products catalog.ProductCollection) (catalog.ProductCollection, error) {
	if len(products) == 0 {
		return catalog.ProductCollection{}, nil
	}

	// Create a copy to avoid mutating the original
	sorted := products.Copy()

	// Sort alphabetically (case-insensitive)
	sort.Slice(sorted, func(i, j int) bool {
		nameI := strings.ToLower(sorted[i].Name)
		nameJ := strings.ToLower(sorted[j].Name)

		// Primary sort: name (alphabetical)
		if nameI != nameJ {
			return nameI < nameJ
		}

		// Tie-breaker: ID for consistent ordering
		return sorted[i].ID < sorted[j].ID
	})

	return sorted, nil
}

// GetStrategy returns the sort strategy
func (s *NameSorter) GetStrategy() catalog.SortStrategy {
	return catalog.SortByName
}

// GetDescription returns a human-readable description
func (s *NameSorter) GetDescription() string {
	return "Sorts products alphabetically by name (case-insensitive)"
}
