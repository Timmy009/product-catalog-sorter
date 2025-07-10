package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"product-catalog-sorting/internal/application"
	"product-catalog-sorting/internal/domain/catalog"
	"product-catalog-sorting/test/testdata"
)

// CatalogTestSuite provides a comprehensive test suite for the catalog system
type CatalogTestSuite struct {
	suite.Suite
	app      *application.Application
	ctx      context.Context
	products []catalog.Product
	logger   *zap.Logger
}

// SetupSuite initializes the test suite
func (suite *CatalogTestSuite) SetupSuite() {
	suite.logger = zap.NewNop()
	suite.ctx = context.Background()

	app, err := application.New(application.Config{
		Logger:  suite.logger,
		Context: suite.ctx,
	})
	suite.Require().NoError(err)
	suite.app = app

	suite.products = testdata.GetTestProducts()
	suite.Require().Len(suite.products, 3, "Should have exactly 3 test products from the challenge")
}

// TestSortingStrategies tests all sorting strategies with the 3 challenge products
func (suite *CatalogTestSuite) TestSortingStrategies() {
	testCases := testdata.GetSortingTestCases()

	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			result, err := suite.app.SortProducts(suite.ctx, tc.Products, tc.Strategy)
			suite.Require().NoError(err)
			suite.Require().Len(result.Products, 3)

			// Verify the expected top product
			suite.Equal(tc.ExpectedTop, result.Products[0].Name,
				"Expected %s to be the top product for strategy %s", tc.ExpectedTop, tc.Strategy)

			// Verify result metadata
			suite.Equal(tc.Strategy, result.Strategy)
			suite.Equal(3, result.ProductCount)
			suite.Greater(result.ExecutionTime, time.Duration(0))
			suite.False(result.SortedAt.IsZero())

			// Verify result validation
			err = result.Validate()
			suite.NoError(err)
		})
	}
}

// TestBatchSortingWithAllStrategies tests batch sorting with all strategies
func (suite *CatalogTestSuite) TestBatchSortingWithAllStrategies() {
	allStrategies := catalog.NewSortStrategySet(catalog.AllSortStrategies()...)

	result, err := suite.app.BatchSort(suite.ctx, suite.products, allStrategies)
	suite.Require().NoError(err)

	// Verify batch result structure
	suite.Equal(len(catalog.AllSortStrategies()), result.StrategyCount)
	suite.Equal(3, result.ProductCount)
	suite.Greater(result.TotalTime, time.Duration(0))
	suite.False(result.ExecutedAt.IsZero())

	// Verify all strategies have results
	for _, strategy := range catalog.AllSortStrategies() {
		strategyResult, exists := result.GetResult(strategy)
		suite.True(exists, "Should have result for strategy %s", strategy)
		suite.NotNil(strategyResult)
		suite.Len(strategyResult.Products, 3)

		// Verify individual result validation
		err := strategyResult.Validate()
		suite.NoError(err)
	}

	// Verify batch result validation
	err = result.Validate()
	suite.NoError(err)
}

// TestProductValidationComprehensive tests product validation with various scenarios
func (suite *CatalogTestSuite) TestProductValidationComprehensive() {
	suite.Run("Valid Products", func() {
		err := suite.app.ValidateProducts(suite.ctx, suite.products)
		suite.NoError(err)
	})

	suite.Run("Individual Product Validation", func() {
		for i, product := range suite.products {
			err := product.Validate()
			suite.NoError(err, "Product at index %d should be valid", i)
			suite.True(product.IsValid())
		}
	})

	suite.Run("Product Business Logic", func() {
		for _, product := range suite.products {
			// Test conversion ratio calculation
			ratio := product.SalesConversionRatio()
			suite.GreaterOrEqual(ratio, 0.0)
			suite.LessOrEqual(ratio, 1.0)

			// Test revenue calculation
			revenue := product.RevenueGenerated()
			expectedRevenue := float64(product.Price) * float64(product.SalesCount)
			suite.Equal(expectedRevenue, revenue)

			// Test days on market
			days := product.DaysOnMarket()
			suite.GreaterOrEqual(days, 0)
		}
	})
}

// TestPerformanceWithLargeDataset tests performance with larger datasets
func (suite *CatalogTestSuite) TestPerformanceWithLargeDataset() {
	sizes := []int{100, 500, 1000}

	for _, size := range sizes {
		suite.Run(fmt.Sprintf("Dataset_%d", size), func() {
			largeDataset := generateTestProducts(size)

			// Test single sort performance
			start := time.Now()
			result, err := suite.app.SortProducts(suite.ctx, largeDataset, catalog.SortBySalesConversionRatio)
			duration := time.Since(start)

			suite.NoError(err)
			suite.Len(result.Products, size)
			suite.Less(duration, 5*time.Second, "Should complete within 5 seconds for %d products", size)

			suite.T().Logf("Sorted %d products in %v", size, duration)
		})
	}
}

// TestConcurrentOperations tests concurrent access to the sorting system
func (suite *CatalogTestSuite) TestConcurrentOperations() {
	const numGoroutines = 10
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			strategy := catalog.AllSortStrategies()[id%len(catalog.AllSortStrategies())]
			_, err := suite.app.SortProducts(suite.ctx, suite.products, strategy)
			results <- err
		}(i)
	}

	// Collect all results
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		suite.NoError(err, "Concurrent operation %d should succeed", i)
	}
}

// TestEdgeCasesAndErrorHandling tests various edge cases
func (suite *CatalogTestSuite) TestEdgeCasesAndErrorHandling() {
	suite.Run("Empty Product Collection", func() {
		emptyProducts := []catalog.Product{}
		result, err := suite.app.SortProducts(suite.ctx, emptyProducts, catalog.SortByPriceAsc)
		suite.NoError(err)
		suite.Empty(result.Products)
		suite.Equal(0, result.ProductCount)
	})

	suite.Run("Single Product", func() {
		singleProduct := []catalog.Product{suite.products[0]}
		result, err := suite.app.SortProducts(suite.ctx, singleProduct, catalog.SortByPriceAsc)
		suite.NoError(err)
		suite.Len(result.Products, 1)
		suite.Equal(suite.products[0].Name, result.Products[0].Name)
	})

	suite.Run("Invalid Strategy", func() {
		_, err := suite.app.SortProducts(suite.ctx, suite.products, catalog.SortStrategy("invalid"))
		suite.Error(err)
		suite.Contains(err.Error(), "invalid sort strategy")
	})

	suite.Run("Nil Products", func() {
		_, err := suite.app.SortProducts(suite.ctx, nil, catalog.SortByPriceAsc)
		suite.Error(err)
		suite.Contains(err.Error(), "cannot be nil")
	})
}

// TestSortStability tests that sorting is stable for equal elements
func (suite *CatalogTestSuite) TestSortStability() {
	// Create products with identical prices but different IDs
	identicalPriceProducts := []catalog.Product{
		{ID: 3, Name: "Product C", Price: 100.0, CreatedAt: time.Now(), SalesCount: 10, ViewsCount: 100},
		{ID: 1, Name: "Product A", Price: 100.0, CreatedAt: time.Now(), SalesCount: 10, ViewsCount: 100},
		{ID: 2, Name: "Product B", Price: 100.0, CreatedAt: time.Now(), SalesCount: 10, ViewsCount: 100},
	}

	result, err := suite.app.SortProducts(suite.ctx, identicalPriceProducts, catalog.SortByPriceAsc)
	suite.NoError(err)

	// Verify all products are present (order may vary due to Go's sort stability)
	foundIDs := make(map[catalog.ProductID]bool)
	for _, product := range result.Products {
		foundIDs[product.ID] = true
	}
	suite.True(foundIDs[1], "Should contain product ID 1")
	suite.True(foundIDs[2], "Should contain product ID 2")
	suite.True(foundIDs[3], "Should contain product ID 3")

	// All should have the same price
	for _, product := range result.Products {
		suite.Equal(catalog.Price(100.0), product.Price)
	}
}

// TestSupportedStrategies tests the supported strategies functionality
func (suite *CatalogTestSuite) TestSupportedStrategies() {
	strategies := suite.app.GetSupportedStrategies()

	suite.NotEmpty(strategies)
	suite.Equal(len(catalog.AllSortStrategies()), strategies.Len())

	// Verify all expected strategies are present
	expectedStrategies := catalog.AllSortStrategies()
	for _, expected := range expectedStrategies {
		suite.True(strategies.Contains(expected), "Should contain strategy %s", expected)
	}

	// Verify strategy validation
	err := strategies.Validate()
	suite.NoError(err)
}

// Helper function to generate test products
func generateTestProducts(count int) []catalog.Product {
	baseTime := time.Now()
	products := make([]catalog.Product, count)

	for i := 0; i < count; i++ {
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

// Run the test suite
func TestCatalogTestSuite(t *testing.T) {
	suite.Run(t, new(CatalogTestSuite))
}
