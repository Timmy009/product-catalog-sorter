package sorting

import (
	"fmt"

	"product-catalog-sorting/internal/domain/catalog"
)

// DefaultSorterFactory implements the SorterFactory interface
type DefaultSorterFactory struct{}

// NewSorterFactory creates a new default sorter factory
func NewSorterFactory() catalog.SorterFactory {
	return &DefaultSorterFactory{}
}

// CreateSorter creates a sorter for the given strategy
func (f *DefaultSorterFactory) CreateSorter(strategy catalog.SortStrategy) (catalog.Sorter, error) {
	switch strategy {
	case catalog.SortByPriceAsc:
		return NewPriceSorter(true), nil
	case catalog.SortByPriceDesc:
		return NewPriceSorter(false), nil
	case catalog.SortBySalesConversionRatio:
		return NewSalesConversionRatioSorter(), nil
	case catalog.SortByCreatedAtDesc:
		return NewCreatedAtSorter(false), nil
	case catalog.SortByCreatedAtAsc:
		return NewCreatedAtSorter(true), nil
	case catalog.SortByPopularity:
		return NewPopularitySorter(), nil
	case catalog.SortByRevenue:
		return NewRevenueSorter(), nil
	case catalog.SortByName:
		return NewNameSorter(), nil
	default:
		return nil, fmt.Errorf("unsupported sort strategy: %s", strategy)
	}
}

// GetSupportedStrategies returns all supported strategies
func (f *DefaultSorterFactory) GetSupportedStrategies() catalog.SortStrategySet {
	return catalog.NewSortStrategySet(catalog.AllSortStrategies()...)
}

// IsSupported checks if a strategy is supported
func (f *DefaultSorterFactory) IsSupported(strategy catalog.SortStrategy) bool {
	_, err := f.CreateSorter(strategy)
	return err == nil
}
