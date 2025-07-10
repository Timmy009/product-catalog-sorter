package sorting

import (
	"context"
	"sort"

	"product-catalog-sorting/internal/domain/catalog"
)

// SalesConversionRatioSorter sorts products by sales conversion ratio
type SalesConversionRatioSorter struct{}

// NewSalesConversionRatioSorter creates a new sales conversion ratio sorter
func NewSalesConversionRatioSorter() catalog.Sorter {
	return &SalesConversionRatioSorter{}
}

// Sort implements the Sorter interface
func (s *SalesConversionRatioSorter) Sort(ctx context.Context, products catalog.ProductCollection) (catalog.ProductCollection, error) {
	if len(products) == 0 {
		return catalog.ProductCollection{}, nil
	}

	// Create a copy to avoid mutating the original
	sorted := products.Copy()

	// Sort by conversion ratio (descending), then by sales count (descending)
	sort.Slice(sorted, func(i, j int) bool {
		ratioI := sorted[i].SalesConversionRatio()
		ratioJ := sorted[j].SalesConversionRatio()

		// Primary sort: conversion ratio (higher is better)
		if ratioI != ratioJ {
			return ratioI > ratioJ
		}

		// Secondary sort: sales count (higher is better) for tie-breaking
		if sorted[i].SalesCount != sorted[j].SalesCount {
			return sorted[i].SalesCount > sorted[j].SalesCount
		}

		// Tertiary sort: ID for consistent ordering
		return sorted[i].ID < sorted[j].ID
	})

	return sorted, nil
}

// GetStrategy returns the sort strategy
func (s *SalesConversionRatioSorter) GetStrategy() catalog.SortStrategy {
	return catalog.SortBySalesConversionRatio
}

// GetDescription returns a human-readable description
func (s *SalesConversionRatioSorter) GetDescription() string {
	return "Sorts products by sales conversion ratio (sales/views) from highest to lowest"
}
