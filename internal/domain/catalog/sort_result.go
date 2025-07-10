package catalog

import (
	"fmt"
	"time"
)

// SortResult represents the result of a sorting operation
// Contains both the sorted products and metadata about the operation
type SortResult struct {
	Products       ProductCollection `json:"products"`
	Strategy       SortStrategy      `json:"strategy"`
	ExecutionTime  time.Duration     `json:"execution_time"`
	ProductCount   int               `json:"product_count"`
	SortedAt       time.Time         `json:"sorted_at"`
}

// NewSortResult creates a new sort result with the given parameters
func NewSortResult(products ProductCollection, strategy SortStrategy, executionTime time.Duration) *SortResult {
	return &SortResult{
		Products:      products,
		Strategy:      strategy,
		ExecutionTime: executionTime,
		ProductCount:  len(products),
		SortedAt:      time.Now(),
	}
}

// Validate ensures the sort result is valid and consistent
func (sr *SortResult) Validate() error {
	if sr == nil {
		return fmt.Errorf("sort result cannot be nil")
	}

	if !sr.Strategy.IsValid() {
		return fmt.Errorf("invalid sort strategy: %s", sr.Strategy)
	}

	if sr.ProductCount != len(sr.Products) {
		return fmt.Errorf("product count mismatch: expected %d, got %d", 
			sr.ProductCount, len(sr.Products))
	}

	if sr.ExecutionTime < 0 {
		return fmt.Errorf("execution time cannot be negative")
	}

	if sr.SortedAt.IsZero() {
		return fmt.Errorf("sorted timestamp must be set")
	}

	return nil
}

// GetTopProducts returns the top N products from the sorted result
func (sr *SortResult) GetTopProducts(n int) ProductCollection {
	if n <= 0 || len(sr.Products) == 0 {
		return ProductCollection{}
	}

	if n >= len(sr.Products) {
		return sr.Products.Copy()
	}

	return sr.Products[:n].Copy()
}

// String provides a detailed string representation of the sort result
func (sr *SortResult) String() string {
	return fmt.Sprintf("SortResult{Strategy: %s, Products: %d, ExecutionTime: %v, SortedAt: %s}",
		sr.Strategy, sr.ProductCount, sr.ExecutionTime, sr.SortedAt.Format(time.RFC3339))
}

// BatchSortResult represents the result of a batch sorting operation
type BatchSortResult struct {
	Results       map[SortStrategy]*SortResult `json:"results"`
	TotalTime     time.Duration                `json:"total_time"`
	StrategyCount int                          `json:"strategy_count"`
	ProductCount  int                          `json:"product_count"`
	ExecutedAt    time.Time                    `json:"executed_at"`
}

// NewBatchSortResult creates a new batch sort result
func NewBatchSortResult(results map[SortStrategy]*SortResult, totalTime time.Duration) *BatchSortResult {
	productCount := 0
	if len(results) > 0 {
		// Get product count from first result (all should be the same)
		for _, result := range results {
			productCount = result.ProductCount
			break
		}
	}

	return &BatchSortResult{
		Results:       results,
		TotalTime:     totalTime,
		StrategyCount: len(results),
		ProductCount:  productCount,
		ExecutedAt:    time.Now(),
	}
}

// GetResult returns the sort result for a specific strategy
func (bsr *BatchSortResult) GetResult(strategy SortStrategy) (*SortResult, bool) {
	result, exists := bsr.Results[strategy]
	return result, exists
}

// Validate ensures the batch sort result is valid
func (bsr *BatchSortResult) Validate() error {
	if bsr == nil {
		return fmt.Errorf("batch sort result cannot be nil")
	}

	if len(bsr.Results) == 0 {
		return fmt.Errorf("batch sort result must contain at least one result")
	}

	if bsr.StrategyCount != len(bsr.Results) {
		return fmt.Errorf("strategy count mismatch: expected %d, got %d",
			bsr.StrategyCount, len(bsr.Results))
	}

	// Validate each individual result
	for strategy, result := range bsr.Results {
		if result == nil {
			return fmt.Errorf("result for strategy %s is nil", strategy)
		}
		if err := result.Validate(); err != nil {
			return fmt.Errorf("invalid result for strategy %s: %w", strategy, err)
		}
	}

	return nil
}

// String provides a string representation of the batch sort result
func (bsr *BatchSortResult) String() string {
	return fmt.Sprintf("BatchSortResult{Strategies: %d, Products: %d, TotalTime: %v}",
		bsr.StrategyCount, bsr.ProductCount, bsr.TotalTime)
}
