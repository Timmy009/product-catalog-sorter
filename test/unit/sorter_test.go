package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"product-catalog-sorting/internal/domain/catalog"
	"product-catalog-sorting/internal/infrastructure/sorting"
)

func TestPriceSorter_Comprehensive(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Expensive", Price: 100.0, CreatedAt: time.Now()},
		{ID: 2, Name: "Cheap", Price: 10.0, CreatedAt: time.Now()},
		{ID: 3, Name: "Medium", Price: 50.0, CreatedAt: time.Now()},
		{ID: 4, Name: "Same Price 1", Price: 25.0, CreatedAt: time.Now()},
		{ID: 5, Name: "Same Price 2", Price: 25.0, CreatedAt: time.Now()},
	}

	ctx := context.Background()

	t.Run("Ascending Order", func(t *testing.T) {
		sorter := sorting.NewPriceSorter(true)
		sorted, err := sorter.Sort(ctx, products)

		require.NoError(t, err)
		require.Len(t, sorted, 5)

		// Verify ascending order
		expectedPrices := []catalog.Price{10.0, 25.0, 25.0, 50.0, 100.0}
		for i, expected := range expectedPrices {
			assert.Equal(t, expected, sorted[i].Price, "Position %d should have price %v", i, expected)
		}

		// Verify strategy and description
		assert.Equal(t, catalog.SortByPriceAsc, sorter.GetStrategy())
		assert.Contains(t, sorter.GetDescription(), "lowest to highest")
	})

	t.Run("Descending Order", func(t *testing.T) {
		sorter := sorting.NewPriceSorter(false)
		sorted, err := sorter.Sort(ctx, products)

		require.NoError(t, err)
		require.Len(t, sorted, 5)

		// Verify descending order
		expectedPrices := []catalog.Price{100.0, 50.0, 25.0, 25.0, 10.0}
		for i, expected := range expectedPrices {
			assert.Equal(t, expected, sorted[i].Price, "Position %d should have price %v", i, expected)
		}

		// Verify strategy and description
		assert.Equal(t, catalog.SortByPriceDesc, sorter.GetStrategy())
		assert.Contains(t, sorter.GetDescription(), "highest to lowest")
	})

	t.Run("Empty Collection", func(t *testing.T) {
		sorter := sorting.NewPriceSorter(true)
		sorted, err := sorter.Sort(ctx, catalog.ProductCollection{})

		require.NoError(t, err)
		assert.Empty(t, sorted)
	})

	t.Run("Single Item", func(t *testing.T) {
		singleProduct := catalog.ProductCollection{
			{ID: 1, Name: "Only One", Price: 42.0, CreatedAt: time.Now()},
		}

		sorter := sorting.NewPriceSorter(true)
		sorted, err := sorter.Sort(ctx, singleProduct)

		require.NoError(t, err)
		require.Len(t, sorted, 1)
		assert.Equal(t, catalog.Price(42.0), sorted[0].Price)
	})

	t.Run("Immutability Check", func(t *testing.T) {
		originalProducts := products.Copy()
		sorter := sorting.NewPriceSorter(true)

		_, err := sorter.Sort(ctx, products)
		require.NoError(t, err)

		// Verify original is unchanged
		assert.Equal(t, originalProducts, products)
	})
}

func TestSalesConversionRatioSorter_Comprehensive(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Low Ratio", SalesCount: 10, ViewsCount: 1000, CreatedAt: time.Now()},   // 0.01
		{ID: 2, Name: "High Ratio", SalesCount: 50, ViewsCount: 100, CreatedAt: time.Now()},   // 0.5
		{ID: 3, Name: "Medium Ratio", SalesCount: 30, ViewsCount: 200, CreatedAt: time.Now()}, // 0.15
		{ID: 4, Name: "Zero Views", SalesCount: 10, ViewsCount: 0, CreatedAt: time.Now()},     // 0.0
		{ID: 5, Name: "Zero Sales", SalesCount: 0, ViewsCount: 100, CreatedAt: time.Now()},    // 0.0
	}

	ctx := context.Background()
	sorter := sorting.NewSalesConversionRatioSorter()

	t.Run("Basic Sorting", func(t *testing.T) {
		sorted, err := sorter.Sort(ctx, products)

		require.NoError(t, err)
		require.Len(t, sorted, 5)

		// Verify descending order by ratio
		expectedOrder := []string{"High Ratio", "Medium Ratio", "Low Ratio", "Zero Views", "Zero Sales"}
		for i, expected := range expectedOrder {
			assert.Equal(t, expected, sorted[i].Name, "Position %d should be %s", i, expected)
		}

		// Verify ratios are in descending order
		for i := 0; i < len(sorted)-1; i++ {
			ratio1 := sorted[i].SalesConversionRatio()
			ratio2 := sorted[i+1].SalesConversionRatio()
			assert.GreaterOrEqual(t, ratio1, ratio2, "Ratio at position %d should be >= ratio at position %d", i, i+1)
		}
	})

	t.Run("Tie Breaking by Sales Count", func(t *testing.T) {
		tieProducts := catalog.ProductCollection{
			{ID: 1, Name: "Same Ratio A", SalesCount: 10, ViewsCount: 100, CreatedAt: time.Now()}, // 0.1
			{ID: 2, Name: "Same Ratio B", SalesCount: 20, ViewsCount: 200, CreatedAt: time.Now()}, // 0.1
			{ID: 3, Name: "Same Ratio C", SalesCount: 15, ViewsCount: 150, CreatedAt: time.Now()}, // 0.1
		}

		sorted, err := sorter.Sort(ctx, tieProducts)
		require.NoError(t, err)

		// With same ratio, should sort by sales count (higher first)
		assert.Equal(t, "Same Ratio B", sorted[0].Name) // 20 sales
		assert.Equal(t, "Same Ratio C", sorted[1].Name) // 15 sales
		assert.Equal(t, "Same Ratio A", sorted[2].Name) // 10 sales
	})

	t.Run("Tie Breaking by ID", func(t *testing.T) {
		identicalProducts := catalog.ProductCollection{
			{ID: 3, Name: "Product C", SalesCount: 10, ViewsCount: 100, CreatedAt: time.Now()},
			{ID: 1, Name: "Product A", SalesCount: 10, ViewsCount: 100, CreatedAt: time.Now()},
			{ID: 2, Name: "Product B", SalesCount: 10, ViewsCount: 100, CreatedAt: time.Now()},
		}

		sorted, err := sorter.Sort(ctx, identicalProducts)
		require.NoError(t, err)

		// With same ratio and sales, should sort by ID (ascending)
		assert.Equal(t, catalog.ProductID(1), sorted[0].ID)
		assert.Equal(t, catalog.ProductID(2), sorted[1].ID)
		assert.Equal(t, catalog.ProductID(3), sorted[2].ID)
	})

	t.Run("Strategy and Description", func(t *testing.T) {
		assert.Equal(t, catalog.SortBySalesConversionRatio, sorter.GetStrategy())
		assert.Contains(t, sorter.GetDescription(), "sales conversion ratio")
		assert.Contains(t, sorter.GetDescription(), "highest to lowest")
	})
}

func TestCreatedAtSorter_Comprehensive(t *testing.T) {
	now := time.Now()
	products := catalog.ProductCollection{
		{ID: 1, Name: "Old", CreatedAt: now.AddDate(-1, 0, 0)},           // 1 year ago
		{ID: 2, Name: "New", CreatedAt: now},                             // now
		{ID: 3, Name: "Medium", CreatedAt: now.AddDate(0, -6, 0)},        // 6 months ago
		{ID: 4, Name: "Same Time 1", CreatedAt: now.AddDate(0, -3, -10)}, // 3 months 10 days ago
		{ID: 5, Name: "Same Time 2", CreatedAt: now.AddDate(0, -3, -20)}, // 3 months 20 days ago
	}

	ctx := context.Background()

	t.Run("Descending Order (Newest First)", func(t *testing.T) {
		sorter := sorting.NewCreatedAtSorter(false)
		sorted, err := sorter.Sort(ctx, products)

		require.NoError(t, err)
		require.Len(t, sorted, 5)

		// Let's debug what we actually get
		t.Logf("Actual sorted order (descending):")
		for i, p := range sorted {
			t.Logf("  %d. %s - %s", i+1, p.Name, p.CreatedAt.Format("2006-01-02"))
		}

		// Just verify the dates are in descending order, don't assume specific names
		for i := 0; i < len(sorted)-1; i++ {
			time1 := sorted[i].CreatedAt
			time2 := sorted[i+1].CreatedAt
			assert.True(t, time1.After(time2) || time1.Equal(time2),
				"Time at position %d (%s) should be >= time at position %d (%s)",
				i, time1.Format("2006-01-02"), i+1, time2.Format("2006-01-02"))
		}

		// Verify the newest is first and oldest is last
		assert.Equal(t, "New", sorted[0].Name, "Newest should be first")
		assert.Equal(t, "Old", sorted[len(sorted)-1].Name, "Oldest should be last")

		assert.Equal(t, catalog.SortByCreatedAtDesc, sorter.GetStrategy())
	})

	t.Run("Ascending Order (Oldest First)", func(t *testing.T) {
		sorter := sorting.NewCreatedAtSorter(true)
		sorted, err := sorter.Sort(ctx, products)

		require.NoError(t, err)
		require.Len(t, sorted, 5)

		// Let's debug what we actually get
		t.Logf("Actual sorted order (ascending):")
		for i, p := range sorted {
			t.Logf("  %d. %s - %s", i+1, p.Name, p.CreatedAt.Format("2006-01-02"))
		}

		// Just verify the dates are in ascending order, don't assume specific names
		for i := 0; i < len(sorted)-1; i++ {
			time1 := sorted[i].CreatedAt
			time2 := sorted[i+1].CreatedAt
			assert.True(t, time1.Before(time2) || time1.Equal(time2),
				"Time at position %d (%s) should be <= time at position %d (%s)",
				i, time1.Format("2006-01-02"), i+1, time2.Format("2006-01-02"))
		}

		// Verify the oldest is first and newest is last
		assert.Equal(t, "Old", sorted[0].Name, "Oldest should be first")
		assert.Equal(t, "New", sorted[len(sorted)-1].Name, "Newest should be last")

		assert.Equal(t, catalog.SortByCreatedAtAsc, sorter.GetStrategy())
	})

	t.Run("Tie Breaking by ID", func(t *testing.T) {
		sameTimeProducts := catalog.ProductCollection{
			{ID: 3, Name: "Product C", CreatedAt: now},
			{ID: 1, Name: "Product A", CreatedAt: now},
			{ID: 2, Name: "Product B", CreatedAt: now},
		}

		sorter := sorting.NewCreatedAtSorter(false)
		sorted, err := sorter.Sort(ctx, sameTimeProducts)
		require.NoError(t, err)

		// With same creation time, should sort by ID (ascending)
		assert.Equal(t, catalog.ProductID(1), sorted[0].ID)
		assert.Equal(t, catalog.ProductID(2), sorted[1].ID)
		assert.Equal(t, catalog.ProductID(3), sorted[2].ID)
	})
}

func TestPopularitySorter_Comprehensive(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Low Views", ViewsCount: 100, SalesCount: 10, CreatedAt: time.Now()},
		{ID: 2, Name: "High Views", ViewsCount: 1000, SalesCount: 50, CreatedAt: time.Now()},
		{ID: 3, Name: "Medium Views", ViewsCount: 500, SalesCount: 25, CreatedAt: time.Now()},
		{ID: 4, Name: "Same Views A", ViewsCount: 300, SalesCount: 20, CreatedAt: time.Now()},
		{ID: 5, Name: "Same Views B", ViewsCount: 300, SalesCount: 30, CreatedAt: time.Now()},
	}

	ctx := context.Background()
	sorter := sorting.NewPopularitySorter()

	t.Run("Basic Sorting", func(t *testing.T) {
		sorted, err := sorter.Sort(ctx, products)

		require.NoError(t, err)
		require.Len(t, sorted, 5)

		// Verify descending order by views
		expectedOrder := []string{"High Views", "Medium Views", "Same Views B", "Same Views A", "Low Views"}
		for i, expected := range expectedOrder {
			assert.Equal(t, expected, sorted[i].Name, "Position %d should be %s", i, expected)
		}

		// Verify views are in descending order
		for i := 0; i < len(sorted)-1; i++ {
			views1 := sorted[i].ViewsCount
			views2 := sorted[i+1].ViewsCount
			assert.GreaterOrEqual(t, views1, views2, "Views at position %d should be >= views at position %d", i, i+1)
		}

		assert.Equal(t, catalog.SortByPopularity, sorter.GetStrategy())
		assert.Contains(t, sorter.GetDescription(), "popularity")
	})

	t.Run("Tie Breaking by Sales Count", func(t *testing.T) {
		// Same Views B should come before Same Views A due to higher sales count
		sorted, err := sorter.Sort(ctx, products)
		require.NoError(t, err)

		// Find the products with same view count
		var sameViewsProducts []catalog.Product
		for _, product := range sorted {
			if product.ViewsCount == 300 {
				sameViewsProducts = append(sameViewsProducts, product)
			}
		}

		require.Len(t, sameViewsProducts, 2)
		assert.Equal(t, "Same Views B", sameViewsProducts[0].Name) // 30 sales
		assert.Equal(t, "Same Views A", sameViewsProducts[1].Name) // 20 sales
	})
}

func TestRevenueSorter_Comprehensive(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Low Revenue", Price: 10.0, SalesCount: 5, CreatedAt: time.Now()},     // 50.0
		{ID: 2, Name: "High Revenue", Price: 100.0, SalesCount: 10, CreatedAt: time.Now()},  // 1000.0
		{ID: 3, Name: "Medium Revenue", Price: 25.0, SalesCount: 8, CreatedAt: time.Now()},  // 200.0
		{ID: 4, Name: "Same Revenue A", Price: 20.0, SalesCount: 5, CreatedAt: time.Now()},  // 100.0
		{ID: 5, Name: "Same Revenue B", Price: 10.0, SalesCount: 10, CreatedAt: time.Now()}, // 100.0
	}

	ctx := context.Background()
	sorter := sorting.NewRevenueSorter()

	t.Run("Basic Sorting", func(t *testing.T) {
		sorted, err := sorter.Sort(ctx, products)

		require.NoError(t, err)
		require.Len(t, sorted, 5)

		// Verify descending order by revenue
		expectedRevenues := []float64{1000.0, 200.0, 100.0, 100.0, 50.0}
		for i, expected := range expectedRevenues {
			actual := sorted[i].RevenueGenerated()
			assert.Equal(t, expected, actual, "Position %d should have revenue %v", i, expected)
		}

		assert.Equal(t, catalog.SortByRevenue, sorter.GetStrategy())
		assert.Contains(t, sorter.GetDescription(), "revenue")
	})

	t.Run("Tie Breaking by Sales Count", func(t *testing.T) {
		sorted, err := sorter.Sort(ctx, products)
		require.NoError(t, err)

		// Find products with same revenue (100.0)
		var sameRevenueProducts []catalog.Product
		for _, product := range sorted {
			if product.RevenueGenerated() == 100.0 {
				sameRevenueProducts = append(sameRevenueProducts, product)
			}
		}

		require.Len(t, sameRevenueProducts, 2)
		assert.Equal(t, "Same Revenue B", sameRevenueProducts[0].Name) // 10 sales
		assert.Equal(t, "Same Revenue A", sameRevenueProducts[1].Name) // 5 sales
	})
}

func TestNameSorter_Comprehensive(t *testing.T) {
	products := catalog.ProductCollection{
		{ID: 1, Name: "Zebra Product", CreatedAt: time.Now()},
		{ID: 2, Name: "Apple Product", CreatedAt: time.Now()},
		{ID: 3, Name: "banana product", CreatedAt: time.Now()}, // lowercase
		{ID: 4, Name: "Cherry Product", CreatedAt: time.Now()},
		{ID: 5, Name: "apple product", CreatedAt: time.Now()}, // lowercase, same as #2
	}

	ctx := context.Background()
	sorter := sorting.NewNameSorter()

	t.Run("Case Insensitive Sorting", func(t *testing.T) {
		sorted, err := sorter.Sort(ctx, products)

		require.NoError(t, err)
		require.Len(t, sorted, 5)

		// Verify alphabetical order (case-insensitive)
		expectedOrder := []string{"Apple Product", "apple product", "banana product", "Cherry Product", "Zebra Product"}
		for i, expected := range expectedOrder {
			assert.Equal(t, expected, sorted[i].Name, "Position %d should be %s", i, expected)
		}

		assert.Equal(t, catalog.SortByName, sorter.GetStrategy())
		assert.Contains(t, sorter.GetDescription(), "alphabetically")
		assert.Contains(t, sorter.GetDescription(), "case-insensitive")
	})

	t.Run("Tie Breaking by ID", func(t *testing.T) {
		sameNameProducts := catalog.ProductCollection{
			{ID: 3, Name: "Same Name", CreatedAt: time.Now()},
			{ID: 1, Name: "Same Name", CreatedAt: time.Now()},
			{ID: 2, Name: "Same Name", CreatedAt: time.Now()},
		}

		sorted, err := sorter.Sort(ctx, sameNameProducts)
		require.NoError(t, err)

		// With same name, should sort by ID (ascending)
		assert.Equal(t, catalog.ProductID(1), sorted[0].ID)
		assert.Equal(t, catalog.ProductID(2), sorted[1].ID)
		assert.Equal(t, catalog.ProductID(3), sorted[2].ID)
	})
}

func TestSorterFactory_Comprehensive(t *testing.T) {
	factory := sorting.NewSorterFactory()

	t.Run("Create All Supported Sorters", func(t *testing.T) {
		strategies := catalog.AllSortStrategies()

		for _, strategy := range strategies {
			t.Run(string(strategy), func(t *testing.T) {
				sorter, err := factory.CreateSorter(strategy)
				require.NoError(t, err)
				assert.NotNil(t, sorter)
				assert.Equal(t, strategy, sorter.GetStrategy())
				assert.NotEmpty(t, sorter.GetDescription())
			})
		}
	})

	t.Run("Unsupported Strategy", func(t *testing.T) {
		_, err := factory.CreateSorter(catalog.SortStrategy("unsupported"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported sort strategy")
	})

	t.Run("Get Supported Strategies", func(t *testing.T) {
		strategies := factory.GetSupportedStrategies()
		assert.NotEmpty(t, strategies)

		// Verify all returned strategies are valid
		for _, strategy := range strategies.ToSlice() {
			assert.True(t, strategy.IsValid())
		}
	})

	t.Run("Is Supported", func(t *testing.T) {
		assert.True(t, factory.IsSupported(catalog.SortByPriceAsc))
		assert.True(t, factory.IsSupported(catalog.SortBySalesConversionRatio))
		assert.False(t, factory.IsSupported(catalog.SortStrategy("invalid")))
	})
}

func TestSorter_ContextCancellation(t *testing.T) {
	products := generateLargeProductSet(1000)

	t.Run("Context Cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		sorter := sorting.NewPriceSorter(true)
		_, err := sorter.Sort(ctx, products)

		// Note: Our current implementation doesn't check context cancellation
		// This test documents the current behavior
		assert.NoError(t, err)
	})

	t.Run("Context Timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(1 * time.Millisecond) // Ensure timeout

		sorter := sorting.NewSalesConversionRatioSorter()
		_, err := sorter.Sort(ctx, products)

		// Note: Our current implementation doesn't check context timeout
		// This test documents the current behavior
		assert.NoError(t, err)
	})
}

func TestSorter_LargeDatasets(t *testing.T) {
	sizes := []int{100, 1000, 5000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Dataset size %d", size), func(t *testing.T) {
			products := generateLargeProductSet(size)
			ctx := context.Background()

			sorters := []struct {
				name   string
				sorter catalog.Sorter
			}{
				{"PriceSorter", sorting.NewPriceSorter(true)},
				{"SalesConversionRatioSorter", sorting.NewSalesConversionRatioSorter()},
				{"PopularitySorter", sorting.NewPopularitySorter()},
				{"NameSorter", sorting.NewNameSorter()},
			}

			for _, s := range sorters {
				t.Run(s.name, func(t *testing.T) {
					start := time.Now()
					sorted, err := s.sorter.Sort(ctx, products)
					duration := time.Since(start)

					require.NoError(t, err)
					assert.Len(t, sorted, size)

					// Performance assertion - should complete within reasonable time
					assert.Less(t, duration, 5*time.Second, "Sorting %d products should complete within 5 seconds", size)

					t.Logf("%s sorted %d products in %v", s.name, size, duration)
				})
			}
		})
	}
}

func TestSorter_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("Nil Context", func(t *testing.T) {
		products := catalog.ProductCollection{
			{ID: 1, Name: "Product", Price: 10.0, CreatedAt: time.Now()},
		}

		sorter := sorting.NewPriceSorter(true)
		// This should not panic even with context.TODO()
		sorted, err := sorter.Sort(context.TODO(), products)
		require.NoError(t, err)
		assert.Len(t, sorted, 1)
	})

	t.Run("Extreme Values", func(t *testing.T) {
		products := catalog.ProductCollection{
			{ID: 1, Name: "Max Price", Price: 999999.99, SalesCount: 1000000, ViewsCount: 1000000, CreatedAt: time.Now()},
			{ID: 2, Name: "Min Price", Price: 0.01, SalesCount: 1, ViewsCount: 1, CreatedAt: time.Now()},
		}

		sorters := []catalog.Sorter{
			sorting.NewPriceSorter(true),
			sorting.NewSalesConversionRatioSorter(),
			sorting.NewPopularitySorter(),
			sorting.NewRevenueSorter(),
		}

		for _, sorter := range sorters {
			sorted, err := sorter.Sort(ctx, products)
			require.NoError(t, err)
			assert.Len(t, sorted, 2)
		}
	})
}

// Helper function to generate large product sets for testing
func generateLargeProductSet(size int) catalog.ProductCollection {
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

// Benchmark tests for performance validation
func BenchmarkPriceSorter_Small(b *testing.B) {
	products := generateLargeProductSet(100)
	sorter := sorting.NewPriceSorter(true)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sorter.Sort(ctx, products)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPriceSorter_Large(b *testing.B) {
	products := generateLargeProductSet(10000)
	sorter := sorting.NewPriceSorter(true)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sorter.Sort(ctx, products)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSalesConversionRatioSorter_Small(b *testing.B) {
	products := generateLargeProductSet(100)
	sorter := sorting.NewSalesConversionRatioSorter()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sorter.Sort(ctx, products)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSalesConversionRatioSorter_Large(b *testing.B) {
	products := generateLargeProductSet(10000)
	sorter := sorting.NewSalesConversionRatioSorter()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sorter.Sort(ctx, products)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAllSorters_Medium(b *testing.B) {
	products := generateLargeProductSet(1000)
	ctx := context.Background()

	sorters := map[string]catalog.Sorter{
		"PriceAsc":             sorting.NewPriceSorter(true),
		"PriceDesc":            sorting.NewPriceSorter(false),
		"SalesConversionRatio": sorting.NewSalesConversionRatioSorter(),
		"CreatedAtDesc":        sorting.NewCreatedAtSorter(false),
		"CreatedAtAsc":         sorting.NewCreatedAtSorter(true),
		"Popularity":           sorting.NewPopularitySorter(),
		"Revenue":              sorting.NewRevenueSorter(),
		"Name":                 sorting.NewNameSorter(),
	}

	for name, sorter := range sorters {
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := sorter.Sort(ctx, products)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
