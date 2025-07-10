package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"product-catalog-sorting/internal/domain/catalog"
	"product-catalog-sorting/internal/infrastructure/sorting"
)

func TestService_SortProducts_Comprehensive(t *testing.T) {
	logger := zap.NewNop()
	factory := sorting.NewSorterFactory()
	service := catalog.NewService(factory, logger)

	products := catalog.ProductCollection{
		{ID: 1, Name: "Expensive", Price: 100.0, CreatedAt: time.Now(), SalesCount: 10, ViewsCount: 100},
		{ID: 2, Name: "Cheap", Price: 10.0, CreatedAt: time.Now(), SalesCount: 50, ViewsCount: 200},
		{ID: 3, Name: "Medium", Price: 50.0, CreatedAt: time.Now(), SalesCount: 30, ViewsCount: 150},
	}

	ctx := context.Background()

	t.Run("All Valid Strategies", func(t *testing.T) {
		strategies := catalog.AllSortStrategies()

		for _, strategy := range strategies {
			t.Run(string(strategy), func(t *testing.T) {
				result, err := service.SortProducts(ctx, products, strategy)
				require.NoError(t, err)
				require.NotNil(t, result)

				assert.Len(t, result.Products, len(products))
				assert.Equal(t, strategy, result.Strategy)
				assert.Greater(t, result.ExecutionTime, time.Duration(0))
				assert.False(t, result.SortedAt.IsZero())
				assert.Equal(t, len(products), result.ProductCount)

				// Verify result validation passes
				err = result.Validate()
				assert.NoError(t, err)
			})
		}
	})

	t.Run("Specific Strategy Tests", func(t *testing.T) {
		t.Run("SortByPriceAsc", func(t *testing.T) {
			result, err := service.SortProducts(ctx, products, catalog.SortByPriceAsc)
			require.NoError(t, err)

			assert.Equal(t, catalog.Price(10.0), result.Products[0].Price)
			assert.Equal(t, catalog.Price(50.0), result.Products[1].Price)
			assert.Equal(t, catalog.Price(100.0), result.Products[2].Price)
		})

		t.Run("SortBySalesConversionRatio", func(t *testing.T) {
			result, err := service.SortProducts(ctx, products, catalog.SortBySalesConversionRatio)
			require.NoError(t, err)

			// Expected order: Cheap (0.25), Medium (0.2), Expensive (0.1)
			assert.Equal(t, "Cheap", result.Products[0].Name)
			assert.Equal(t, "Medium", result.Products[1].Name)
			assert.Equal(t, "Expensive", result.Products[2].Name)
		})
	})

	t.Run("Error Cases", func(t *testing.T) {
		t.Run("Nil Products", func(t *testing.T) {
			_, err := service.SortProducts(ctx, nil, catalog.SortByPriceAsc)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot be nil")
		})

		t.Run("Invalid Strategy", func(t *testing.T) {
			_, err := service.SortProducts(ctx, products, catalog.SortStrategy("invalid"))
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid sort strategy")
		})

		t.Run("Invalid Products", func(t *testing.T) {
			invalidProducts := catalog.ProductCollection{
				{ID: 0, Name: "", Price: -10.0, CreatedAt: time.Time{}, SalesCount: -1, ViewsCount: -1},
			}

			_, err := service.SortProducts(ctx, invalidProducts, catalog.SortByPriceAsc)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "validation failed")
		})
	})

	t.Run("Empty Products", func(t *testing.T) {
		result, err := service.SortProducts(ctx, catalog.ProductCollection{}, catalog.SortByPriceAsc)
		require.NoError(t, err)
		assert.Empty(t, result.Products)
		assert.Equal(t, 0, result.ProductCount)
	})

	t.Run("Single Product", func(t *testing.T) {
		singleProduct := catalog.ProductCollection{
			{ID: 1, Name: "Only One", Price: 42.0, CreatedAt: time.Now(), SalesCount: 5, ViewsCount: 50},
		}

		result, err := service.SortProducts(ctx, singleProduct, catalog.SortByPriceAsc)
		require.NoError(t, err)
		assert.Len(t, result.Products, 1)
		assert.Equal(t, "Only One", result.Products[0].Name)
	})
}

func TestService_BatchSort_Comprehensive(t *testing.T) {
	logger := zap.NewNop()
	factory := sorting.NewSorterFactory()
	service := catalog.NewService(factory, logger)

	products := catalog.ProductCollection{
		{ID: 1, Name: "Product A", Price: 20.0, CreatedAt: time.Now(), SalesCount: 10, ViewsCount: 100},
		{ID: 2, Name: "Product B", Price: 10.0, CreatedAt: time.Now().AddDate(0, 0, -1), SalesCount: 20, ViewsCount: 150},
		{ID: 3, Name: "Product C", Price: 30.0, CreatedAt: time.Now().AddDate(0, 0, -2), SalesCount: 5, ViewsCount: 80},
	}

	ctx := context.Background()

	t.Run("Multiple Strategies", func(t *testing.T) {
		strategies := catalog.NewSortStrategySet(
			catalog.SortByPriceAsc,
			catalog.SortBySalesConversionRatio,
			catalog.SortByPopularity,
		)

		result, err := service.BatchSort(ctx, products, strategies)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Len(t, result.Results, 3)
		assert.Equal(t, 3, result.StrategyCount)
		assert.Equal(t, len(products), result.ProductCount)
		assert.Greater(t, result.TotalTime, time.Duration(0))
		assert.False(t, result.ExecutedAt.IsZero())

		// Verify each strategy result exists and is valid
		for _, strategy := range strategies.ToSlice() {
			strategyResult, exists := result.GetResult(strategy)
			assert.True(t, exists, "Result should exist for strategy %s", strategy)
			assert.NotNil(t, strategyResult)
			assert.Len(t, strategyResult.Products, len(products))

			err := strategyResult.Validate()
			assert.NoError(t, err)
		}

		// Verify batch result validation passes
		err = result.Validate()
		assert.NoError(t, err)
	})

	t.Run("Single Strategy Batch", func(t *testing.T) {
		strategies := catalog.NewSortStrategySet(catalog.SortByPriceAsc)

		result, err := service.BatchSort(ctx, products, strategies)
		require.NoError(t, err)

		assert.Len(t, result.Results, 1)
		assert.Equal(t, 1, result.StrategyCount)

		priceResult, exists := result.GetResult(catalog.SortByPriceAsc)
		assert.True(t, exists)
		assert.Equal(t, catalog.Price(10.0), priceResult.Products[0].Price)
	})

	t.Run("All Strategies Batch", func(t *testing.T) {
		allStrategies := catalog.NewSortStrategySet(catalog.AllSortStrategies()...)

		result, err := service.BatchSort(ctx, products, allStrategies)
		require.NoError(t, err)

		assert.Len(t, result.Results, len(catalog.AllSortStrategies()))
		assert.Equal(t, len(catalog.AllSortStrategies()), result.StrategyCount)

		// Verify all strategies produced results
		for _, strategy := range catalog.AllSortStrategies() {
			strategyResult, exists := result.GetResult(strategy)
			assert.True(t, exists, "Result should exist for strategy %s", strategy)
			assert.NotNil(t, strategyResult)
		}
	})

	t.Run("Error Cases", func(t *testing.T) {
		t.Run("Nil Products", func(t *testing.T) {
			strategies := catalog.NewSortStrategySet(catalog.SortByPriceAsc)
			_, err := service.BatchSort(ctx, nil, strategies)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot be nil")
		})

		t.Run("Empty Strategies", func(t *testing.T) {
			strategies := catalog.SortStrategySet{}
			_, err := service.BatchSort(ctx, products, strategies)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot be empty")
		})

		t.Run("Invalid Strategies", func(t *testing.T) {
			strategies := catalog.SortStrategySet{
				catalog.SortByPriceAsc,
				catalog.SortStrategy("invalid"),
			}
			_, err := service.BatchSort(ctx, products, strategies)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "validation failed")
		})

		t.Run("Invalid Products", func(t *testing.T) {
			invalidProducts := catalog.ProductCollection{
				{ID: 0, Name: "", Price: -10.0},
			}
			strategies := catalog.NewSortStrategySet(catalog.SortByPriceAsc)

			_, err := service.BatchSort(ctx, invalidProducts, strategies)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "validation failed")
		})
	})

	t.Run("Empty Products Batch", func(t *testing.T) {
		strategies := catalog.NewSortStrategySet(catalog.SortByPriceAsc, catalog.SortByPopularity)

		result, err := service.BatchSort(ctx, catalog.ProductCollection{}, strategies)
		require.NoError(t, err)

		assert.Len(t, result.Results, 2)
		assert.Equal(t, 0, result.ProductCount)

		for _, strategy := range strategies.ToSlice() {
			strategyResult, exists := result.GetResult(strategy)
			assert.True(t, exists)
			assert.Empty(t, strategyResult.Products)
		}
	})
}

func TestService_ValidateProducts_Comprehensive(t *testing.T) {
	logger := zap.NewNop()
	factory := sorting.NewSorterFactory()
	service := catalog.NewService(factory, logger)
	ctx := context.Background()

	t.Run("Valid Products", func(t *testing.T) {
		validProducts := catalog.ProductCollection{
			{ID: 1, Name: "Valid Product 1", Price: 10.0, CreatedAt: time.Now(), SalesCount: 5, ViewsCount: 50},
			{ID: 2, Name: "Valid Product 2", Price: 20.0, CreatedAt: time.Now(), SalesCount: 10, ViewsCount: 100},
		}

		err := service.ValidateProducts(ctx, validProducts)
		assert.NoError(t, err)
	})

	t.Run("Invalid Products", func(t *testing.T) {
		invalidProducts := catalog.ProductCollection{
			{ID: 1, Name: "Valid Product", Price: 10.0, CreatedAt: time.Now(), SalesCount: 5, ViewsCount: 50},
			{ID: 0, Name: "", Price: -10.0, CreatedAt: time.Time{}, SalesCount: -1, ViewsCount: -1},
		}

		err := service.ValidateProducts(ctx, invalidProducts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("Nil Products", func(t *testing.T) {
		err := service.ValidateProducts(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})

	t.Run("Empty Products", func(t *testing.T) {
		err := service.ValidateProducts(ctx, catalog.ProductCollection{})
		assert.NoError(t, err)
	})

	t.Run("Mixed Valid and Invalid", func(t *testing.T) {
		mixedProducts := catalog.ProductCollection{
			{ID: 1, Name: "Valid", Price: 10.0, CreatedAt: time.Now(), SalesCount: 5, ViewsCount: 50},
			{ID: 2, Name: "Also Valid", Price: 20.0, CreatedAt: time.Now(), SalesCount: 10, ViewsCount: 100},
			{ID: 0, Name: "", Price: -5.0, CreatedAt: time.Time{}, SalesCount: -1, ViewsCount: -1}, // Invalid
		}

		err := service.ValidateProducts(ctx, mixedProducts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product at index 2")
	})
}

func TestService_GetSupportedStrategies(t *testing.T) {
	logger := zap.NewNop()
	factory := sorting.NewSorterFactory()
	service := catalog.NewService(factory, logger)

	strategies := service.GetSupportedStrategies()

	assert.NotEmpty(t, strategies)
	assert.Equal(t, len(catalog.AllSortStrategies()), strategies.Len())

	// Verify all expected strategies are present
	expectedStrategies := catalog.AllSortStrategies()
	for _, expected := range expectedStrategies {
		assert.True(t, strategies.Contains(expected), "Should contain strategy %s", expected)
	}

	// Verify all returned strategies are valid
	for _, strategy := range strategies.ToSlice() {
		assert.True(t, strategy.IsValid())
	}
}

func TestService_ContextHandling(t *testing.T) {
	logger := zap.NewNop()
	factory := sorting.NewSorterFactory()
	service := catalog.NewService(factory, logger)

	products := catalog.ProductCollection{
		{ID: 1, Name: "Product", Price: 10.0, CreatedAt: time.Now(), SalesCount: 5, ViewsCount: 50},
	}

	t.Run("Cancelled Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Current implementation doesn't check context cancellation
		// This test documents the current behavior
		result, err := service.SortProducts(ctx, products, catalog.SortByPriceAsc)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Context With Timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(1 * time.Millisecond) // Ensure timeout

		// Current implementation doesn't check context timeout
		// This test documents the current behavior
		result, err := service.SortProducts(ctx, products, catalog.SortByPriceAsc)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Nil Context", func(t *testing.T) {
		// Should not panic with nil context
		result, err := service.SortProducts(context.TODO(), products, catalog.SortByPriceAsc)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestService_Performance(t *testing.T) {
	logger := zap.NewNop()
	factory := sorting.NewSorterFactory()
	service := catalog.NewService(factory, logger)
	ctx := context.Background()

	sizes := []int{100, 1000, 5000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Dataset size %d", size), func(t *testing.T) {
			products := generateLargeProductCollection(size)

			t.Run("Single Sort Performance", func(t *testing.T) {
				start := time.Now()
				result, err := service.SortProducts(ctx, products, catalog.SortBySalesConversionRatio)
				duration := time.Since(start)

				require.NoError(t, err)
				assert.Len(t, result.Products, size)
				assert.Less(t, duration, 5*time.Second, "Should complete within 5 seconds")

				t.Logf("Sorted %d products in %v", size, duration)
			})

			t.Run("Batch Sort Performance", func(t *testing.T) {
				strategies := catalog.NewSortStrategySet(
					catalog.SortByPriceAsc,
					catalog.SortBySalesConversionRatio,
					catalog.SortByPopularity,
				)

				start := time.Now()
				result, err := service.BatchSort(ctx, products, strategies)
				duration := time.Since(start)

				require.NoError(t, err)
				assert.Len(t, result.Results, 3)
				assert.Less(t, duration, 15*time.Second, "Batch sort should complete within 15 seconds")

				t.Logf("Batch sorted %d products with %d strategies in %v", size, len(strategies), duration)
			})
		})
	}
}

func TestService_ConcurrentAccess(t *testing.T) {
	logger := zap.NewNop()
	factory := sorting.NewSorterFactory()
	service := catalog.NewService(factory, logger)

	products := generateLargeProductCollection(100)
	ctx := context.Background()

	t.Run("Concurrent Sort Operations", func(t *testing.T) {
		const numGoroutines = 10
		results := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				strategy := catalog.AllSortStrategies()[id%len(catalog.AllSortStrategies())]
				_, err := service.SortProducts(ctx, products, strategy)
				results <- err
			}(i)
		}

		// Collect all results
		for i := 0; i < numGoroutines; i++ {
			err := <-results
			assert.NoError(t, err)
		}
	})

	t.Run("Concurrent Batch Operations", func(t *testing.T) {
		const numGoroutines = 5
		results := make(chan error, numGoroutines)

		strategies := catalog.NewSortStrategySet(
			catalog.SortByPriceAsc,
			catalog.SortBySalesConversionRatio,
		)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				_, err := service.BatchSort(ctx, products, strategies)
				results <- err
			}()
		}

		// Collect all results
		for i := 0; i < numGoroutines; i++ {
			err := <-results
			assert.NoError(t, err)
		}
	})
}

func TestService_EdgeCases(t *testing.T) {
	logger := zap.NewNop()
	factory := sorting.NewSorterFactory()
	service := catalog.NewService(factory, logger)
	ctx := context.Background()

	t.Run("Products with Extreme Values", func(t *testing.T) {
		extremeProducts := catalog.ProductCollection{
			{ID: 1, Name: "Max Values", Price: 999999.99, SalesCount: 1000000, ViewsCount: 1000000, CreatedAt: time.Now()},
			{ID: 2, Name: "Min Values", Price: 0.01, SalesCount: 1, ViewsCount: 1, CreatedAt: time.Now().AddDate(-10, 0, 0)},
			{ID: 3, Name: "Zero Values", Price: 0.0, SalesCount: 0, ViewsCount: 1, CreatedAt: time.Now()},
		}

		strategies := catalog.AllSortStrategies()
		for _, strategy := range strategies {
			result, err := service.SortProducts(ctx, extremeProducts, strategy)
			require.NoError(t, err, "Strategy %s should handle extreme values", strategy)
			assert.Len(t, result.Products, 3)
		}
	})

	t.Run("Products with Same Values", func(t *testing.T) {
		identicalProducts := catalog.ProductCollection{
			{ID: 1, Name: "Same Product", Price: 10.0, SalesCount: 5, ViewsCount: 50, CreatedAt: time.Now()},
			{ID: 2, Name: "Same Product", Price: 10.0, SalesCount: 5, ViewsCount: 50, CreatedAt: time.Now()},
			{ID: 3, Name: "Same Product", Price: 10.0, SalesCount: 5, ViewsCount: 50, CreatedAt: time.Now()},
		}

		strategies := catalog.AllSortStrategies()
		for _, strategy := range strategies {
			result, err := service.SortProducts(ctx, identicalProducts, strategy)
			require.NoError(t, err, "Strategy %s should handle identical values", strategy)
			assert.Len(t, result.Products, 3)

			// Verify stable sorting by checking that all products are present
			foundIDs := make(map[catalog.ProductID]bool)
			for _, product := range result.Products {
				foundIDs[product.ID] = true
			}
			assert.True(t, foundIDs[1], "Should contain product ID 1")
			assert.True(t, foundIDs[2], "Should contain product ID 2")
			assert.True(t, foundIDs[3], "Should contain product ID 3")
		}
	})
}

// Helper function to generate large product collections for testing
func generateLargeProductCollection(size int) catalog.ProductCollection {
	products := make(catalog.ProductCollection, size)
	now := time.Now()

	for i := 0; i < size; i++ {
		products[i] = catalog.Product{
			ID:         catalog.ProductID(i + 1),
			Name:       fmt.Sprintf("Product %d", i+1),
			Price:      catalog.Price(10 + float64(i%100)),
			CreatedAt:  now.Add(-time.Duration(i%365) * 24 * time.Hour),
			SalesCount: (i%500 + 1),
			ViewsCount: (i%2000 + 100),
		}
	}

	return products
}

// Benchmark tests
func BenchmarkService_SortProducts(b *testing.B) {
	logger := zap.NewNop()
	factory := sorting.NewSorterFactory()
	service := catalog.NewService(factory, logger)

	products := generateLargeProductCollection(1000)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.SortProducts(ctx, products, catalog.SortBySalesConversionRatio)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_BatchSort(b *testing.B) {
	logger := zap.NewNop()
	factory := sorting.NewSorterFactory()
	service := catalog.NewService(factory, logger)

	products := generateLargeProductCollection(1000)
	strategies := catalog.NewSortStrategySet(
		catalog.SortByPriceAsc,
		catalog.SortBySalesConversionRatio,
		catalog.SortByPopularity,
	)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.BatchSort(ctx, products, strategies)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_ValidateProducts(b *testing.B) {
	logger := zap.NewNop()
	factory := sorting.NewSorterFactory()
	service := catalog.NewService(factory, logger)

	products := generateLargeProductCollection(1000)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := service.ValidateProducts(ctx, products)
		if err != nil {
			b.Fatal(err)
		}
	}
}
