package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"product-catalog-sorting/internal/application"
	"product-catalog-sorting/internal/domain/catalog"
	"product-catalog-sorting/test/testdata"
)

func TestSimpleWorkflow(t *testing.T) {
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

	t.Run("Basic Price Sort", func(t *testing.T) {
		result, err := app.SortProducts(ctx, products, catalog.SortByPriceAsc)
		require.NoError(t, err)
		require.Len(t, result.Products, 3)
		
		// Should be sorted by price ascending - Coffee Table is cheapest at $10.00
		assert.Equal(t, catalog.Price(10.00), result.Products[0].Price)
		assert.Equal(t, "Coffee Table", result.Products[0].Name)
	})

	t.Run("Basic Batch Sort", func(t *testing.T) {
		strategies := catalog.NewSortStrategySet(
			catalog.SortByPriceAsc,
			catalog.SortBySalesConversionRatio,
		)
		
		results, err := app.BatchSort(ctx, products, strategies)
		require.NoError(t, err)
		require.Len(t, results.Results, 2)
		
		// Verify both results exist
		priceResult, exists := results.GetResult(catalog.SortByPriceAsc)
		assert.True(t, exists)
		assert.Len(t, priceResult.Products, 3)
		
		ratioResult, exists := results.GetResult(catalog.SortBySalesConversionRatio)
		assert.True(t, exists)
		assert.Len(t, ratioResult.Products, 3)
	})
}
