package sorting

import (
	"context"
	"sort"

	"product-catalog-sorting/internal/domain/catalog"
)

// CreatedAtSorter sorts products by creation date
type CreatedAtSorter struct {
	ascending bool
}

// NewCreatedAtSorter creates a new creation date sorter
func NewCreatedAtSorter(ascending bool) catalog.Sorter {
	return &CreatedAtSorter{
		ascending: ascending,
	}
}

// Sort implements the Sorter interface
func (s *CreatedAtSorter) Sort(ctx context.Context, products catalog.ProductCollection) (catalog.ProductCollection, error) {
	if len(products) == 0 {
		return catalog.ProductCollection{}, nil
	}

	// Create a copy to avoid mutating the original
	sorted := products.Copy()

	// Sort by creation date with consistent tie-breaking
	sort.Slice(sorted, func(i, j int) bool {
		timeI := sorted[i].CreatedAt
		timeJ := sorted[j].CreatedAt

		// Primary sort: creation date
		if !timeI.Equal(timeJ) {
			if s.ascending {
				return timeI.Before(timeJ)
			}
			return timeI.After(timeJ)
		}

		// Tie-breaker: ID for consistent ordering
		return sorted[i].ID < sorted[j].ID
	})

	return sorted, nil
}

// GetStrategy returns the sort strategy
func (s *CreatedAtSorter) GetStrategy() catalog.SortStrategy {
	if s.ascending {
		return catalog.SortByCreatedAtAsc
	}
	return catalog.SortByCreatedAtDesc
}

// GetDescription returns a human-readable description
func (s *CreatedAtSorter) GetDescription() string {
	if s.ascending {
		return "Sorts products by creation date from oldest to newest"
	}
	return "Sorts products by creation date from newest to oldest"
}
