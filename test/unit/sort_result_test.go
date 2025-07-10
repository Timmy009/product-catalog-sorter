package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"product-catalog-sorting/internal/domain/catalog"
)

func TestNewSortResult(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now()},
		{ID: 2, Name: "Product 2", Price: 20.0, CreatedAt: time.Now()},
	}

	strategy := catalog.SortByPriceAsc
	executionTime := 50 * time.Millisecond

	result := catalog.NewSortResult(products, strategy, executionTime)

	assert.NotNil(t, result)
	assert.Equal(t, products, result.Products)
	assert.Equal(t, strategy, result.Strategy)
	assert.Equal(t, executionTime, result.ExecutionTime)
	assert.Equal(t, len(products), result.ProductCount)
	assert.False(t, result.SortedAt.IsZero())
}

func TestSortResult_Validate(t *testing.T) {
	validProducts := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now()},
		{ID: 2, Name: "Product 2", Price: 20.0, CreatedAt: time.Now()},
	}

	tests := []struct {
		name          string
		result        *catalog.SortResult
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid result",
			result: &catalog.SortResult{
				Products:      validProducts,
				Strategy:      catalog.SortByPriceAsc,
				ExecutionTime: 50 * time.Millisecond,
				ProductCount:  len(validProducts),
				SortedAt:      time.Now(),
			},
			expectError: false,
		},
		{
			name:          "Nil result",
			result:        nil,
			expectError:   true,
			errorContains: "cannot be nil",
		},
		{
			name: "Invalid strategy",
			result: &catalog.SortResult{
				Products:      validProducts,
				Strategy:      catalog.SortStrategy("invalid"),
				ExecutionTime: 50 * time.Millisecond,
				ProductCount:  len(validProducts),
				SortedAt:      time.Now(),
			},
			expectError:   true,
			errorContains: "invalid sort strategy",
		},
		{
			name: "Product count mismatch",
			result: &catalog.SortResult{
				Products:      validProducts,
				Strategy:      catalog.SortByPriceAsc,
				ExecutionTime: 50 * time.Millisecond,
				ProductCount:  5, // wrong count
				SortedAt:      time.Now(),
			},
			expectError:   true,
			errorContains: "product count mismatch",
		},
		{
			name: "Negative execution time",
			result: &catalog.SortResult{
				Products:      validProducts,
				Strategy:      catalog.SortByPriceAsc,
				ExecutionTime: -50 * time.Millisecond,
				ProductCount:  len(validProducts),
				SortedAt:      time.Now(),
			},
			expectError:   true,
			errorContains: "execution time cannot be negative",
		},
		{
			name: "Zero sorted timestamp",
			result: &catalog.SortResult{
				Products:      validProducts,
				Strategy:      catalog.SortByPriceAsc,
				ExecutionTime: 50 * time.Millisecond,
				ProductCount:  len(validProducts),
				SortedAt:      time.Time{},
			},
			expectError:   true,
			errorContains: "sorted timestamp must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.result.Validate()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSortResult_GetTopProducts(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "First", Price: 10.0},
		{ID: 2, Name: "Second", Price: 20.0},
		{ID: 3, Name: "Third", Price: 30.0},
		{ID: 4, Name: "Fourth", Price: 40.0},
		{ID: 5, Name: "Fifth", Price: 50.0},
	}

	result := catalog.NewSortResult(products, catalog.SortByPriceAsc, time.Millisecond)

	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{"Get top 3", 3, 3},
		{"Get top 10 (more than available)", 10, 5},
		{"Get top 0", 0, 0},
		{"Get top -1", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			top := result.GetTopProducts(tt.n)
			assert.Len(t, top, tt.expected)

			if tt.expected > 0 {
				// Verify it's a copy, not the original
				top[0].Name = "Modified"
				assert.NotEqual(t, "Modified", result.Products[0].Name)
			}
		})
	}

	// Test with empty products
	emptyResult := catalog.NewSortResult(catalog.ProductCollection{}, catalog.SortByPriceAsc, time.Millisecond)
	top := emptyResult.GetTopProducts(3)
	assert.Empty(t, top)
}

func TestSortResult_String(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now()},
	}

	result := catalog.NewSortResult(products, catalog.SortByPriceAsc, 50*time.Millisecond)

	str := result.String()
	assert.Contains(t, str, "SortResult{")
	assert.Contains(t, str, "price_asc")
	assert.Contains(t, str, "Products: 1")
	assert.Contains(t, str, "50ms")
}

func TestNewBatchSortResult(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now()},
		{ID: 2, Name: "Product 2", Price: 20.0, CreatedAt: time.Now()},
	}

	results := map[catalog.SortStrategy]*catalog.SortResult{
		catalog.SortByPriceAsc:  catalog.NewSortResult(products, catalog.SortByPriceAsc, 30*time.Millisecond),
		catalog.SortByPriceDesc: catalog.NewSortResult(products, catalog.SortByPriceDesc, 25*time.Millisecond),
	}

	totalTime := 100 * time.Millisecond

	batchResult := catalog.NewBatchSortResult(results, totalTime)

	assert.NotNil(t, batchResult)
	assert.Equal(t, results, batchResult.Results)
	assert.Equal(t, totalTime, batchResult.TotalTime)
	assert.Equal(t, len(results), batchResult.StrategyCount)
	assert.Equal(t, len(products), batchResult.ProductCount)
	assert.False(t, batchResult.ExecutedAt.IsZero())
}

func TestBatchSortResult_GetResult(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now()},
	}

	result1 := catalog.NewSortResult(products, catalog.SortByPriceAsc, 30*time.Millisecond)
	result2 := catalog.NewSortResult(products, catalog.SortByPriceDesc, 25*time.Millisecond)

	results := map[catalog.SortStrategy]*catalog.SortResult{
		catalog.SortByPriceAsc:  result1,
		catalog.SortByPriceDesc: result2,
	}

	batchResult := catalog.NewBatchSortResult(results, 100*time.Millisecond)

	// Test existing strategy
	foundResult, exists := batchResult.GetResult(catalog.SortByPriceAsc)
	assert.True(t, exists)
	assert.Equal(t, result1, foundResult)

	// Test non-existing strategy
	notFoundResult, exists := batchResult.GetResult(catalog.SortByPopularity)
	assert.False(t, exists)
	assert.Nil(t, notFoundResult)
}

func TestBatchSortResult_Validate(t *testing.T) {
	validProducts := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now()},
	}

	validResult := catalog.NewSortResult(validProducts, catalog.SortByPriceAsc, 30*time.Millisecond)

	tests := []struct {
		name          string
		batchResult   *catalog.BatchSortResult
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid batch result",
			batchResult: &catalog.BatchSortResult{
				Results: map[catalog.SortStrategy]*catalog.SortResult{
					catalog.SortByPriceAsc: validResult,
				},
				TotalTime:     100 * time.Millisecond,
				StrategyCount: 1,
				ProductCount:  1,
				ExecutedAt:    time.Now(),
			},
			expectError: false,
		},
		{
			name:          "Nil batch result",
			batchResult:   nil,
			expectError:   true,
			errorContains: "cannot be nil",
		},
		{
			name: "Empty results",
			batchResult: &catalog.BatchSortResult{
				Results:       map[catalog.SortStrategy]*catalog.SortResult{},
				TotalTime:     100 * time.Millisecond,
				StrategyCount: 0,
				ProductCount:  0,
				ExecutedAt:    time.Now(),
			},
			expectError:   true,
			errorContains: "must contain at least one result",
		},
		{
			name: "Strategy count mismatch",
			batchResult: &catalog.BatchSortResult{
				Results: map[catalog.SortStrategy]*catalog.SortResult{
					catalog.SortByPriceAsc: validResult,
				},
				TotalTime:     100 * time.Millisecond,
				StrategyCount: 2, // wrong count
				ProductCount:  1,
				ExecutedAt:    time.Now(),
			},
			expectError:   true,
			errorContains: "strategy count mismatch",
		},
		{
			name: "Nil result in map",
			batchResult: &catalog.BatchSortResult{
				Results: map[catalog.SortStrategy]*catalog.SortResult{
					catalog.SortByPriceAsc: nil,
				},
				TotalTime:     100 * time.Millisecond,
				StrategyCount: 1,
				ProductCount:  1,
				ExecutedAt:    time.Now(),
			},
			expectError:   true,
			errorContains: "is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.batchResult.Validate()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBatchSortResult_String(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now()},
	}

	results := map[catalog.SortStrategy]*catalog.SortResult{
		catalog.SortByPriceAsc: catalog.NewSortResult(products, catalog.SortByPriceAsc, 30*time.Millisecond),
	}

	batchResult := catalog.NewBatchSortResult(results, 100*time.Millisecond)

	str := batchResult.String()
	assert.Contains(t, str, "BatchSortResult{")
	assert.Contains(t, str, "Strategies: 1")
	assert.Contains(t, str, "Products: 1")
	assert.Contains(t, str, "100ms")
}

// Edge cases and performance tests
func TestSortResult_EdgeCases(t *testing.T) {
	t.Run("Empty products collection", func(t *testing.T) {
		emptyProducts := catalog.ProductCollection{}
		result := catalog.NewSortResult(emptyProducts, catalog.SortByPriceAsc, time.Millisecond)

		assert.Equal(t, 0, result.ProductCount)
		assert.Empty(t, result.Products)

		top := result.GetTopProducts(5)
		assert.Empty(t, top)
	})

	t.Run("Zero execution time", func(t *testing.T) {
		products := catalog.ProductCollection{
			{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now()},
		}

		result := catalog.NewSortResult(products, catalog.SortByPriceAsc, 0)
		assert.Equal(t, time.Duration(0), result.ExecutionTime)

		err := result.Validate()
		assert.NoError(t, err) // Zero execution time is valid
	})
}

func TestBatchSortResult_EdgeCases(t *testing.T) {
	t.Run("Empty products in all results", func(t *testing.T) {
		emptyProducts := catalog.ProductCollection{}

		results := map[catalog.SortStrategy]*catalog.SortResult{
			catalog.SortByPriceAsc: catalog.NewSortResult(emptyProducts, catalog.SortByPriceAsc, time.Millisecond),
		}

		batchResult := catalog.NewBatchSortResult(results, time.Millisecond)
		assert.Equal(t, 0, batchResult.ProductCount)
	})

	t.Run("Single strategy batch", func(t *testing.T) {
		products := catalog.ProductCollection{
			{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now()},
		}

		results := map[catalog.SortStrategy]*catalog.SortResult{
			catalog.SortByPriceAsc: catalog.NewSortResult(products, catalog.SortByPriceAsc, time.Millisecond),
		}

		batchResult := catalog.NewBatchSortResult(results, time.Millisecond)
		assert.Equal(t, 1, batchResult.StrategyCount)
		assert.Equal(t, 1, batchResult.ProductCount)
	})
}

// Benchmark tests
func BenchmarkSortResult_Validate(b *testing.B) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now()},
		{ID: 2, Name: "Product 2", Price: 20.0, CreatedAt: time.Now()},
	}

	result := catalog.NewSortResult(products, catalog.SortByPriceAsc, 50*time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = result.Validate()
	}
}

func BenchmarkSortResult_GetTopProducts(b *testing.B) {
	products := make(catalog.ProductCollection, 1000)
	for i := 0; i < 1000; i++ {
		products[i] = catalog.Product{
			ID:        catalog.ProductID(i + 1),
			Name:      "Product",
			Price:     catalog.Price(float64(i + 1)),
			CreatedAt: time.Now(),
		}
	}

	result := catalog.NewSortResult(products, catalog.SortByPriceAsc, time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = result.GetTopProducts(10)
	}
}

func BenchmarkBatchSortResult_Validate(b *testing.B) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now()},
	}

	results := map[catalog.SortStrategy]*catalog.SortResult{
		catalog.SortByPriceAsc:  catalog.NewSortResult(products, catalog.SortByPriceAsc, time.Millisecond),
		catalog.SortByPriceDesc: catalog.NewSortResult(products, catalog.SortByPriceDesc, time.Millisecond),
	}

	batchResult := catalog.NewBatchSortResult(results, time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = batchResult.Validate()
	}
}
