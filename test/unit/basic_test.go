package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"product-catalog-sorting/internal/application"
	"product-catalog-sorting/internal/domain/catalog"
	"product-catalog-sorting/test/testdata"
)

func TestBasicFunctionality(t *testing.T) {
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

	t.Run("Sort by Price Ascending", func(t *testing.T) {
		result, err := app.SortProducts(ctx, products, catalog.SortByPriceAsc)
		require.NoError(t, err)
		require.Len(t, result.Products, 3)

		// Should be sorted by price ascending: Coffee Table ($10.00), Alabaster Table ($12.99), Zebra Table ($44.49)
		assert.Equal(t, "Coffee Table", result.Products[0].Name)
		assert.Equal(t, catalog.Price(10.00), result.Products[0].Price)
		
		assert.Equal(t, "Alabaster Table", result.Products[1].Name)
		assert.Equal(t, catalog.Price(12.99), result.Products[1].Price)
		
		assert.Equal(t, "Zebra Table", result.Products[2].Name)
		assert.Equal(t, catalog.Price(44.49), result.Products[2].Price)
	})

	t.Run("Sort by Price Descending", func(t *testing.T) {
		result, err := app.SortProducts(ctx, products, catalog.SortByPriceDesc)
		require.NoError(t, err)
		require.Len(t, result.Products, 3)

		// Should be sorted by price descending: Zebra Table ($44.49), Alabaster Table ($12.99), Coffee Table ($10.00)
		assert.Equal(t, "Zebra Table", result.Products[0].Name)
		assert.Equal(t, catalog.Price(44.49), result.Products[0].Price)
	})

	t.Run("Sort by Sales Conversion Ratio", func(t *testing.T) {
		result, err := app.SortProducts(ctx, products, catalog.SortBySalesConversionRatio)
		require.NoError(t, err)
		require.Len(t, result.Products, 3)

		// Expected order by conversion ratio:
		// Zebra Table: 301/3279 = 9.18%
		// Coffee Table: 1048/20123 = 5.21%
		// Alabaster Table: 32/730 = 4.38%
		assert.Equal(t, "Zebra Table", result.Products[0].Name)
		assert.Equal(t, "Coffee Table", result.Products[1].Name)
		assert.Equal(t, "Alabaster Table", result.Products[2].Name)
	})

	t.Run("Sort by Revenue", func(t *testing.T) {
		result, err := app.SortProducts(ctx, products, catalog.SortByRevenue)
		require.NoError(t, err)
		require.Len(t, result.Products, 3)

		// Expected order by revenue:
		// Coffee Table: 1048 * $10.00 = $10,480
		// Zebra Table: 301 * $44.49 = $13,391.49
		// Alabaster Table: 32 * $12.99 = $415.68
		assert.Equal(t, "Zebra Table", result.Products[0].Name) // Highest revenue
		assert.Equal(t, "Coffee Table", result.Products[1].Name)
		assert.Equal(t, "Alabaster Table", result.Products[2].Name)
	})

	t.Run("Batch Sort", func(t *testing.T) {
		strategies := catalog.NewSortStrategySet(
			catalog.SortByPriceAsc,
			catalog.SortBySalesConversionRatio,
		)

		results, err := app.BatchSort(ctx, products, strategies)
		require.NoError(t, err)
		require.Len(t, results.Results, 2)

		// Check price ascending result
		priceResult, exists := results.GetResult(catalog.SortByPriceAsc)
		assert.True(t, exists)
		assert.Equal(t, "Coffee Table", priceResult.Products[0].Name)

		// Check sales ratio result
		ratioResult, exists := results.GetResult(catalog.SortBySalesConversionRatio)
		assert.True(t, exists)
		assert.Equal(t, "Zebra Table", ratioResult.Products[0].Name)
	})

	t.Run("Get Supported Strategies", func(t *testing.T) {
		strategies := app.GetSupportedStrategies()
		assert.NotEmpty(t, strategies)
		assert.True(t, strategies.Contains(catalog.SortByPriceAsc))
		assert.True(t, strategies.Contains(catalog.SortBySalesConversionRatio))
	})

	t.Run("Validate Products", func(t *testing.T) {
		err := app.ValidateProducts(ctx, products)
		assert.NoError(t, err)

		// Test with invalid product
		invalidProducts := []catalog.Product{
			{ID: 0, Name: "", Price: -10.0}, // Invalid product
		}
		err = app.ValidateProducts(ctx, invalidProducts)
		assert.Error(t, err)
	})
}

func TestProductValidation(t *testing.T) {
	t.Run("Valid Product", func(t *testing.T) {
		product := catalog.Product{
			ID:         1,
			Name:       "Valid Product",
			Price:      10.0,
			CreatedAt:  time.Now(),
			SalesCount: 5,
			ViewsCount: 50,
		}

		err := product.Validate()
		assert.NoError(t, err)
		assert.True(t, product.IsValid())
	})

	t.Run("Invalid Product", func(t *testing.T) {
		product := catalog.Product{
			ID:         0,
			Name:       "",
			Price:      -10.0,
			CreatedAt:  time.Time{},
			SalesCount: -1,
			ViewsCount: -1,
		}

		err := product.Validate()
		assert.Error(t, err)
		assert.False(t, product.IsValid())
	})
}

func TestSortStrategies(t *testing.T) {
	t.Run("All Strategies Are Valid", func(t *testing.T) {
		strategies := catalog.AllSortStrategies()
		assert.NotEmpty(t, strategies)

		for _, strategy := range strategies {
			assert.True(t, strategy.IsValid(), "Strategy %s should be valid", strategy)
			assert.NotEmpty(t, strategy.Description())
			assert.Greater(t, strategy.Priority(), 0)
		}
	})

	t.Run("Strategy Set Operations", func(t *testing.T) {
		strategies := catalog.NewSortStrategySet(
			catalog.SortByPriceAsc,
			catalog.SortBySalesConversionRatio,
		)

		assert.Equal(t, 2, strategies.Len())
		assert.True(t, strategies.Contains(catalog.SortByPriceAsc))
		assert.False(t, strategies.Contains(catalog.SortByPopularity))

		err := strategies.Validate()
		assert.NoError(t, err)
	})
}
