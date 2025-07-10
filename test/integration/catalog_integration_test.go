package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"product-catalog-sorting/internal/application"
	"product-catalog-sorting/internal/domain/catalog"
	"product-catalog-sorting/test/testdata"
)

// TestCatalogIntegration tests the complete catalog system end-to-end
func TestCatalogIntegration(t *testing.T) {
	// Initialize application
	logger := zap.NewNop()
	app, err := application.New(application.Config{
		Logger:  logger,
		Context: context.Background(),
	})
	require.NoError(t, err)

	// Use the exact 3 products from the code challenge
	products := testdata.GetTestProducts()
	require.Len(t, products, 3, "Should have exactly 3 products from the challenge")

	ctx := context.Background()

	t.Run("End-to-End Sorting Workflow", func(t *testing.T) {
		// Test all sorting strategies
		strategies := catalog.AllSortStrategies()
		
		for _, strategy := range strategies {
			t.Run(fmt.Sprintf("Strategy_%s", strings.ReplaceAll(string(strategy), "_", "")), func(t *testing.T) {
				result, err := app.SortProducts(ctx, products, strategy)
				require.NoError(t, err)
				require.NotNil(t, result)
				
				// Verify result structure
				assert.Len(t, result.Products, 3)
				assert.Equal(t, strategy, result.Strategy)
				assert.Greater(t, result.ExecutionTime, time.Duration(0))
				assert.Equal(t, 3, result.ProductCount)
				
				// Verify result validation
				err = result.Validate()
				assert.NoError(t, err)
				
				t.Logf("Strategy %s completed in %v", strategy, result.ExecutionTime)
			})
		}
	})

	t.Run("A/B Testing Scenario", func(t *testing.T) {
		// Simulate A/B testing with multiple strategies
		strategies := catalog.NewSortStrategySet(
			catalog.SortByPriceAsc,
			catalog.SortBySalesConversionRatio,
			catalog.SortByPopularity,
			catalog.SortByRevenue,
		)
		
		result, err := app.BatchSort(ctx, products, strategies)
		require.NoError(t, err)
		require.NotNil(t, result)
		
		// Verify batch result structure
		assert.Equal(t, 4, result.StrategyCount)
		assert.Equal(t, 3, result.ProductCount)
		assert.Greater(t, result.TotalTime, time.Duration(0))
		
		// Verify all strategies have results
		for _, strategy := range strategies.ToSlice() {
			strategyResult, exists := result.GetResult(strategy)
			assert.True(t, exists, "Should have result for strategy %s", strategy)
			assert.NotNil(t, strategyResult)
			assert.Len(t, strategyResult.Products, 3)
		}
		
		// Verify batch result validation
		err = result.Validate()
		assert.NoError(t, err)
		
		t.Logf("A/B testing completed with %d strategies in %v", len(strategies), result.TotalTime)
	})

	t.Run("Product Validation Integration", func(t *testing.T) {
		// Test product validation
		err := app.ValidateProducts(ctx, products)
		assert.NoError(t, err)
		
		// Test with invalid products
		invalidProducts := []catalog.Product{
			{ID: 0, Name: "", Price: -10.0}, // Invalid product
		}
		err = app.ValidateProducts(ctx, invalidProducts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("Supported Strategies Integration", func(t *testing.T) {
		strategies := app.GetSupportedStrategies()
		
		assert.NotEmpty(t, strategies)
		assert.Equal(t, len(catalog.AllSortStrategies()), strategies.Len())
		
		// Verify all expected strategies are present
		expectedStrategies := catalog.AllSortStrategies()
		for _, expected := range expectedStrategies {
			assert.True(t, strategies.Contains(expected), "Should contain strategy %s", expected)
		}
	})

	t.Run("Performance Integration Test", func(t *testing.T) {
		// Generate larger dataset for performance testing
		largeProducts := generateLargeDataset(1000)
		
		start := time.Now()
		result, err := app.SortProducts(ctx, largeProducts, catalog.SortBySalesConversionRatio)
		duration := time.Since(start)
		
		require.NoError(t, err)
		assert.Len(t, result.Products, 1000)
		assert.Less(t, duration, 100*time.Millisecond, "Should sort 1000 products quickly")
		
		t.Logf("Sorted 1000 products in %v", duration)
	})
}

// generateLargeDataset creates a large dataset for performance testing
func generateLargeDataset(size int) []catalog.Product {
	products := make([]catalog.Product, size)
	baseTime := time.Now()
	
	for i := 0; i < size; i++ {
		products[i] = catalog.Product{
			ID:         catalog.ProductID(i + 1),
			Name:       fmt.Sprintf("Product %d", i+1),
			Price:      catalog.Price(10 + float64(i%100)),
			CreatedAt:  baseTime.Add(-time.Duration(i%365) * 24 * time.Hour),
			SalesCount: (i%500 + 1),
			ViewsCount: (i%2000 + 100),
		}
	}
	
	return products
}
