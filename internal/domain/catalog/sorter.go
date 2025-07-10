package catalog

import (
	"context"
)

// Sorter defines the interface for product sorting implementations
type Sorter interface {
	// Sort applies the sorting strategy to a collection of products
	Sort(ctx context.Context, products ProductCollection) (ProductCollection, error)
	
	// GetStrategy returns the sort strategy this sorter implements
	GetStrategy() SortStrategy
	
	// GetDescription returns a human-readable description
	GetDescription() string
}

// SorterFactory creates sorters for different strategies
type SorterFactory interface {
	// CreateSorter creates a sorter for the given strategy
	CreateSorter(strategy SortStrategy) (Sorter, error)
	
	// GetSupportedStrategies returns all supported strategies
	GetSupportedStrategies() SortStrategySet
	
	// IsSupported checks if a strategy is supported
	IsSupported(strategy SortStrategy) bool
}
