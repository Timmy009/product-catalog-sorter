package catalog

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Service defines the core business operations for the catalog domain
type Service interface {
	// SortProducts sorts a collection of products using the specified strategy
	SortProducts(ctx context.Context, products ProductCollection, strategy SortStrategy) (*SortResult, error)
	
	// BatchSort sorts products using multiple strategies simultaneously
	BatchSort(ctx context.Context, products ProductCollection, strategies SortStrategySet) (*BatchSortResult, error)
	
	// GetSupportedStrategies returns all supported sorting strategies
	GetSupportedStrategies() SortStrategySet
	
	// ValidateProducts validates a collection of products
	ValidateProducts(ctx context.Context, products ProductCollection) error
}

// DefaultService implements the Service interface
type DefaultService struct {
	sorterFactory SorterFactory
	logger        *zap.Logger
}

// NewService creates a new catalog service with dependencies
func NewService(factory SorterFactory, logger *zap.Logger) Service {
	return &DefaultService{
		sorterFactory: factory,
		logger:        logger,
	}
}

// SortProducts implements the core sorting business logic
func (s *DefaultService) SortProducts(ctx context.Context, products ProductCollection, strategy SortStrategy) (*SortResult, error) {
	// Validate inputs
	if err := s.validateSortRequest(products, strategy); err != nil {
		return nil, fmt.Errorf("sort request validation failed: %w", err)
	}

	// Record start time
	start := time.Now()

	s.logger.Debug("Starting sort operation",
		zap.String("strategy", string(strategy)),
		zap.Int("product_count", len(products)),
	)

	// Create sorter
	sorter, err := s.sorterFactory.CreateSorter(strategy)
	if err != nil {
		return nil, fmt.Errorf("failed to create sorter for strategy %s: %w", strategy, err)
	}

	// Execute sorting
	sortedProducts, err := sorter.Sort(ctx, products)
	if err != nil {
		return nil, fmt.Errorf("sorting failed for strategy %s: %w", strategy, err)
	}

	// Calculate execution time
	executionTime := time.Since(start)

	// Create result
	result := NewSortResult(sortedProducts, strategy, executionTime)

	s.logger.Debug("Sort operation completed",
		zap.String("strategy", string(strategy)),
		zap.Int("product_count", len(sortedProducts)),
		zap.Duration("execution_time", executionTime),
	)

	return result, nil
}

// BatchSort sorts products using multiple strategies
func (s *DefaultService) BatchSort(ctx context.Context, products ProductCollection, strategies SortStrategySet) (*BatchSortResult, error) {
	// Validate inputs
	if err := s.validateBatchSortRequest(products, strategies); err != nil {
		return nil, fmt.Errorf("batch sort request validation failed: %w", err)
	}

	start := time.Now()
	results := make(map[SortStrategy]*SortResult)

	s.logger.Debug("Starting batch sort operation",
		zap.Int("strategy_count", len(strategies)),
		zap.Int("product_count", len(products)),
	)

	// Execute each sorting strategy
	for _, strategy := range strategies {
		result, err := s.SortProducts(ctx, products, strategy)
		if err != nil {
			return nil, fmt.Errorf("batch sort failed for strategy %s: %w", strategy, err)
		}
		results[strategy] = result
	}

	totalTime := time.Since(start)
	batchResult := NewBatchSortResult(results, totalTime)

	s.logger.Debug("Batch sort operation completed",
		zap.Int("strategy_count", len(strategies)),
		zap.Duration("total_time", totalTime),
	)

	return batchResult, nil
}

// GetSupportedStrategies returns all supported sorting strategies
func (s *DefaultService) GetSupportedStrategies() SortStrategySet {
	return s.sorterFactory.GetSupportedStrategies()
}

// ValidateProducts validates a collection of products
func (s *DefaultService) ValidateProducts(ctx context.Context, products ProductCollection) error {
	if products == nil {
		return fmt.Errorf("products collection cannot be nil")
	}

	return products.Validate()
}

// validateSortRequest validates the sort request parameters
func (s *DefaultService) validateSortRequest(products ProductCollection, strategy SortStrategy) error {
	if products == nil {
		return fmt.Errorf("products collection cannot be nil")
	}

	if !strategy.IsValid() {
		return fmt.Errorf("invalid sort strategy: %s", strategy)
	}

	if err := products.Validate(); err != nil {
		return fmt.Errorf("product validation failed: %w", err)
	}

	return nil
}

// validateBatchSortRequest validates the batch sort request parameters
func (s *DefaultService) validateBatchSortRequest(products ProductCollection, strategies SortStrategySet) error {
	if products == nil {
		return fmt.Errorf("products collection cannot be nil")
	}

	if len(strategies) == 0 {
		return fmt.Errorf("strategies set cannot be empty")
	}

	if err := strategies.Validate(); err != nil {
		return fmt.Errorf("strategies validation failed: %w", err)
	}

	if err := products.Validate(); err != nil {
		return fmt.Errorf("product validation failed: %w", err)
	}

	return nil
}
