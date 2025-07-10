package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"product-catalog-sorting/internal/domain/catalog"
)

func TestProduct_SalesConversionRatio(t *testing.T) {
	tests := []struct {
		name           string
		salesCount     int
		viewsCount     int
		expectedRatio  float64
	}{
		{
			name:          "Normal conversion ratio",
			salesCount:    50,
			viewsCount:    200,
			expectedRatio: 0.25,
		},
		{
			name:          "Zero views should return zero",
			salesCount:    10,
			viewsCount:    0,
			expectedRatio: 0.0,
		},
		{
			name:          "Zero sales should return zero",
			salesCount:    0,
			viewsCount:    100,
			expectedRatio: 0.0,
		},
		{
			name:          "Perfect conversion",
			salesCount:    100,
			viewsCount:    100,
			expectedRatio: 1.0,
		},
		{
			name:          "High conversion ratio",
			salesCount:    75,
			viewsCount:    100,
			expectedRatio: 0.75,
		},
		{
			name:          "Low conversion ratio",
			salesCount:    1,
			viewsCount:    1000,
			expectedRatio: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := catalog.Product{
				SalesCount: tt.salesCount,
				ViewsCount: tt.viewsCount,
			}
			
			ratio := product.SalesConversionRatio()
			assert.Equal(t, tt.expectedRatio, ratio)
		})
	}
}

func TestProduct_IsHighPerformer(t *testing.T) {
	tests := []struct {
		name           string
		salesCount     int
		viewsCount     int
		expectedResult bool
	}{
		{
			name:           "High performer - good ratio and sales",
			salesCount:     60,
			viewsCount:     1000, // 6% conversion rate
			expectedResult: true,
		},
		{
			name:           "Low performer - low ratio",
			salesCount:     30,
			viewsCount:     1000, // 3% conversion rate
			expectedResult: false,
		},
		{
			name:           "Low performer - good ratio but low sales",
			salesCount:     40,
			viewsCount:     500, // 8% conversion rate but only 40 sales
			expectedResult: false,
		},
		{
			name:           "Edge case - exactly 5% and 50 sales",
			salesCount:     50,
			viewsCount:     1000, // exactly 5%
			expectedResult: false, // should be > 5%, not >=
		},
		{
			name:           "High performer - exactly 51 sales and >5%",
			salesCount:     51,
			viewsCount:     1000, // 5.1% conversion rate
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := catalog.Product{
				SalesCount: tt.salesCount,
				ViewsCount: tt.viewsCount,
			}
			
			result := product.IsHighPerformer()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestProduct_DaysOnMarket(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name         string
		createdAt    time.Time
		expectedDays int
	}{
		{
			name:         "Created today",
			createdAt:    now,
			expectedDays: 0,
		},
		{
			name:         "Created yesterday",
			createdAt:    now.AddDate(0, 0, -1),
			expectedDays: 1,
		},
		{
			name:         "Created a week ago",
			createdAt:    now.AddDate(0, 0, -7),
			expectedDays: 7,
		},
		{
			name:         "Created a month ago",
			createdAt:    now.AddDate(0, -1, 0),
			expectedDays: 30, // approximately
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := catalog.Product{
				CreatedAt: tt.createdAt,
			}
			
			days := product.DaysOnMarket()
			// Allow some tolerance for timing differences
			assert.InDelta(t, tt.expectedDays, days, 1)
		})
	}
}

func TestProduct_RevenueGenerated(t *testing.T) {
	tests := []struct {
		name            string
		price           catalog.Price
		salesCount      int
		expectedRevenue float64
	}{
		{
			name:            "Basic revenue calculation",
			price:           catalog.Price(25.50),
			salesCount:      10,
			expectedRevenue: 255.0,
		},
		{
			name:            "Zero sales",
			price:           catalog.Price(100.0),
			salesCount:      0,
			expectedRevenue: 0.0,
		},
		{
			name:            "High volume sales",
			price:           catalog.Price(9.99),
			salesCount:      1000,
			expectedRevenue: 9990.0,
		},
		{
			name:            "Expensive item low volume",
			price:           catalog.Price(1999.99),
			salesCount:      5,
			expectedRevenue: 9999.95,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := catalog.Product{
				Price:      tt.price,
				SalesCount: tt.salesCount,
			}
			
			revenue := product.RevenueGenerated()
			assert.Equal(t, tt.expectedRevenue, revenue)
		})
	}
}

func TestProduct_Validate(t *testing.T) {
	validProduct := catalog.Product{
		ID:         1,
		Name:       "Test Product",
		Price:      10.99,
		CreatedAt:  time.Now(),
		SalesCount: 5,
		ViewsCount: 50,
	}

	tests := []struct {
		name          string
		product       catalog.Product
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid product",
			product:     validProduct,
			expectError: false,
		},
		{
			name: "Invalid ID - zero",
			product: catalog.Product{
				ID:         0,
				Name:       "Test",
				Price:      10.0,
				CreatedAt:  time.Now(),
				SalesCount: 5,
				ViewsCount: 50,
			},
			expectError:   true,
			errorContains: "must be positive",
		},
		{
			name: "Invalid ID - negative",
			product: catalog.Product{
				ID:         -1,
				Name:       "Test",
				Price:      10.0,
				CreatedAt:  time.Now(),
				SalesCount: 5,
				ViewsCount: 50,
			},
			expectError:   true,
			errorContains: "must be positive",
		},
		{
			name: "Empty name",
			product: catalog.Product{
				ID:         1,
				Name:       "",
				Price:      10.0,
				CreatedAt:  time.Now(),
				SalesCount: 5,
				ViewsCount: 50,
			},
			expectError:   true,
			errorContains: "cannot be empty",
		},
		{
			name: "Name too long",
			product: catalog.Product{
				ID:         1,
				Name:       string(make([]byte, 256)), // 256 characters
				Price:      10.0,
				CreatedAt:  time.Now(),
				SalesCount: 5,
				ViewsCount: 50,
			},
			expectError:   true,
			errorContains: "cannot exceed 255 characters",
		},
		{
			name: "Negative price",
			product: catalog.Product{
				ID:         1,
				Name:       "Test",
				Price:      -10.0,
				CreatedAt:  time.Now(),
				SalesCount: 5,
				ViewsCount: 50,
			},
			expectError:   true,
			errorContains: "cannot be negative",
		},
		{
			name: "Price too high",
			product: catalog.Product{
				ID:         1,
				Name:       "Test",
				Price:      1000000.0,
				CreatedAt:  time.Now(),
				SalesCount: 5,
				ViewsCount: 50,
			},
			expectError:   true,
			errorContains: "exceeds maximum allowed value",
		},
		{
			name: "Zero created date",
			product: catalog.Product{
				ID:         1,
				Name:       "Test",
				Price:      10.0,
				CreatedAt:  time.Time{},
				SalesCount: 5,
				ViewsCount: 50,
			},
			expectError:   true,
			errorContains: "must be set",
		},
		{
			name: "Future created date",
			product: catalog.Product{
				ID:         1,
				Name:       "Test",
				Price:      10.0,
				CreatedAt:  time.Now().Add(24 * time.Hour),
				SalesCount: 5,
				ViewsCount: 50,
			},
			expectError:   true,
			errorContains: "cannot be in the future",
		},
		{
			name: "Negative sales count",
			product: catalog.Product{
				ID:         1,
				Name:       "Test",
				Price:      10.0,
				CreatedAt:  time.Now(),
				SalesCount: -1,
				ViewsCount: 50,
			},
			expectError:   true,
			errorContains: "cannot be negative",
		},
		{
			name: "Negative views count",
			product: catalog.Product{
				ID:         1,
				Name:       "Test",
				Price:      10.0,
				CreatedAt:  time.Now(),
				SalesCount: 5,
				ViewsCount: -1,
			},
			expectError:   true,
			errorContains: "cannot be negative",
		},
		{
			name: "Sales exceed views",
			product: catalog.Product{
				ID:         1,
				Name:       "Test",
				Price:      10.0,
				CreatedAt:  time.Now(),
				SalesCount: 100,
				ViewsCount: 50,
			},
			expectError:   true,
			errorContains: "sales count cannot exceed views count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			
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

func TestProduct_IsValid(t *testing.T) {
	validProduct := catalog.Product{
		ID:         1,
		Name:       "Test Product",
		Price:      10.99,
		CreatedAt:  time.Now(),
		SalesCount: 5,
		ViewsCount: 50,
	}
	
	invalidProduct := catalog.Product{
		ID:         0,
		Name:       "",
		Price:      -10.0,
		CreatedAt:  time.Time{},
		SalesCount: -1,
		ViewsCount: -1,
	}
	
	assert.True(t, validProduct.IsValid())
	assert.False(t, invalidProduct.IsValid())
}

func TestProduct_String(t *testing.T) {
	product := catalog.Product{
		ID:         1,
		Name:       "Test Product",
		Price:      25.99,
		SalesCount: 10,
		ViewsCount: 100,
		CreatedAt:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	
	str := product.String()
	assert.Contains(t, str, "Test Product")
	assert.Contains(t, str, "25.99")
	assert.Contains(t, str, "2023-01-01")
}

func TestProductCollection_Validate(t *testing.T) {
	validProducts := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now(), SalesCount: 5, ViewsCount: 50},
		{ID: 2, Name: "Product 2", Price: 20.0, CreatedAt: time.Now(), SalesCount: 10, ViewsCount: 100},
	}
	
	invalidProducts := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0, CreatedAt: time.Now(), SalesCount: 5, ViewsCount: 50},
		{ID: 0, Name: "", Price: -10.0, CreatedAt: time.Time{}, SalesCount: -1, ViewsCount: -1},
	}
	
	emptyProducts := catalog.ProductCollection{}
	
	t.Run("Valid collection", func(t *testing.T) {
		err := validProducts.Validate()
		assert.NoError(t, err)
	})
	
	t.Run("Invalid collection", func(t *testing.T) {
		err := invalidProducts.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product at index 1")
	})
	
	t.Run("Empty collection", func(t *testing.T) {
		err := emptyProducts.Validate()
		assert.NoError(t, err)
	})
}

func TestProductCollection_Copy(t *testing.T) {
	original := catalog.ProductCollection{
		{ID: 1, Name: "Product 1", Price: 10.0},
		{ID: 2, Name: "Product 2", Price: 20.0},
	}
	
	copied := original.Copy()
	
	// Verify copy is independent
	require.Equal(t, len(original), len(copied))
	assert.Equal(t, original[0].ID, copied[0].ID)
	
	// Modify original and verify copy is unchanged
	original[0].Name = "Modified"
	assert.NotEqual(t, original[0].Name, copied[0].Name)
	
	// Test nil collection
	var nilCollection catalog.ProductCollection
	nilCopy := nilCollection.Copy()
	assert.Nil(t, nilCopy)
}

func TestProductCollection_FilterHighPerformers(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "High Performer", SalesCount: 60, ViewsCount: 1000}, // 6%
		{ID: 2, Name: "Low Performer", SalesCount: 30, ViewsCount: 1000},  // 3%
		{ID: 3, Name: "Another High", SalesCount: 80, ViewsCount: 1000},   // 8%
	}
	
	highPerformers := products.FilterHighPerformers()
	
	assert.Len(t, highPerformers, 2)
	assert.Equal(t, "High Performer", highPerformers[0].Name)
	assert.Equal(t, "Another High", highPerformers[1].Name)
}

func TestProductCollection_TotalRevenue(t *testing.T) {
	products := catalog.ProductCollection{
		{Price: 10.0, SalesCount: 5},  // 50.0
		{Price: 20.0, SalesCount: 3},  // 60.0
		{Price: 15.0, SalesCount: 2},  // 30.0
	}
	
	totalRevenue := products.TotalRevenue()
	assert.Equal(t, 140.0, totalRevenue)
	
	// Test empty collection
	emptyProducts := catalog.ProductCollection{}
	assert.Equal(t, 0.0, emptyProducts.TotalRevenue())
}

func TestProductCollection_AverageConversionRatio(t *testing.T) {
	products := catalog.ProductCollection{
		{SalesCount: 10, ViewsCount: 100}, // 0.1
		{SalesCount: 20, ViewsCount: 100}, // 0.2
		{SalesCount: 30, ViewsCount: 100}, // 0.3
	}
	
	avgRatio := products.AverageConversionRatio()
	assert.InDelta(t, 0.2, avgRatio, 0.0001) // Allow small floating point differences
	
	// Test empty collection
	emptyProducts := catalog.ProductCollection{}
	assert.Equal(t, 0.0, emptyProducts.AverageConversionRatio())
}

func TestProductCollection_SortInterface(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 3, Name: "Third"},
		{ID: 1, Name: "First"},
		{ID: 2, Name: "Second"},
	}
	
	// Test Len
	assert.Equal(t, 3, products.Len())
	
	// Test Less
	assert.True(t, products.Less(1, 0))  // ID 1 < ID 3
	assert.False(t, products.Less(0, 1)) // ID 3 > ID 1
	
	// Test Swap
	products.Swap(0, 1)
	assert.Equal(t, catalog.ProductID(1), products[0].ID)
	assert.Equal(t, catalog.ProductID(3), products[1].ID)
}

func TestProductID_Methods(t *testing.T) {
	validID := catalog.ProductID(123)
	invalidID := catalog.ProductID(0)
	negativeID := catalog.ProductID(-1)
	
	// Test String
	assert.Equal(t, "ProductID(123)", validID.String())
	
	// Test IsValid
	assert.True(t, validID.IsValid())
	assert.False(t, invalidID.IsValid())
	assert.False(t, negativeID.IsValid())
}

func TestPrice_Methods(t *testing.T) {
	validPrice := catalog.Price(25.99)
	zeroPrice := catalog.Price(0.0)
	negativePrice := catalog.Price(-10.0)
	tooHighPrice := catalog.Price(1000000.0)
	
	// Test String
	assert.Equal(t, "$25.99", validPrice.String())
	
	// Test IsValid
	assert.True(t, validPrice.IsValid())
	assert.True(t, zeroPrice.IsValid())
	assert.False(t, negativePrice.IsValid())
	assert.False(t, tooHighPrice.IsValid())
	
	// Test ToFloat64
	assert.Equal(t, 25.99, validPrice.ToFloat64())
}

// Benchmark tests
func BenchmarkProduct_SalesConversionRatio(b *testing.B) {
	product := catalog.Product{
		SalesCount: 50,
		ViewsCount: 200,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = product.SalesConversionRatio()
	}
}

func BenchmarkProduct_Validate(b *testing.B) {
	product := catalog.Product{
		ID:         1,
		Name:       "Test Product",
		Price:      10.99,
		CreatedAt:  time.Now(),
		SalesCount: 5,
		ViewsCount: 50,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = product.Validate()
	}
}

func BenchmarkProductCollection_Copy(b *testing.B) {
	products := make(catalog.ProductCollection, 1000)
	for i := 0; i < 1000; i++ {
		products[i] = catalog.Product{
			ID:         catalog.ProductID(i + 1),
			Name:       "Product",
			Price:      10.0,
			CreatedAt:  time.Now(),
			SalesCount: 5,
			ViewsCount: 50,
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = products.Copy()
	}
}
