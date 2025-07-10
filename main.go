package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"

	"product-catalog-sorting/internal/application"
	"product-catalog-sorting/internal/domain/catalog"
)

func main() {
	// Initialize logger
	logger := initializeLogger()
	defer func() {
		if err := logger.Sync(); err != nil {
			// Ignore sync errors on stdout/stderr
		}
	}()

	logger.Info("Starting Product Catalog Sorting Demo")

	// Initialize application
	app, err := application.New(application.Config{
		Logger:  logger,
		Context: context.Background(),
	})
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	// Initialize sample product data - using the exact products from the code challenge
	products := []catalog.Product{
		{
			ID:         1,
			Name:       "Alabaster Table",
			Price:      12.99,
			CreatedAt:  parseDate("2019-01-04"),
			SalesCount: 32,
			ViewsCount: 730,
		},
		{
			ID:         2,
			Name:       "Zebra Table",
			Price:      44.49,
			CreatedAt:  parseDate("2012-01-04"),
			SalesCount: 301,
			ViewsCount: 3279,
		},
		{
			ID:         3,
			Name:       "Coffee Table",
			Price:      10.00,
			CreatedAt:  parseDate("2014-05-28"),
			SalesCount: 1048,
			ViewsCount: 20123,
		},
	}

	ctx := context.Background()

	// Demonstrate different sorting strategies
	fmt.Println("=== Product Catalog Sorting Demo ===")

	// Sort by price (ascending)
	fmt.Println("\n📊 Sorted by Price (Low to High):")
	sortedByPrice, err := app.SortProducts(ctx, products, catalog.SortByPriceAsc)
	if err != nil {
		log.Fatal(err)
	}
	displayProducts(sortedByPrice.Products)

	// Sort by sales conversion ratio (descending)
	fmt.Println("\n🎯 Sorted by Sales Conversion Ratio (High to Low):")
	sortedByRatio, err := app.SortProducts(ctx, products, catalog.SortBySalesConversionRatio)
	if err != nil {
		log.Fatal(err)
	}
	displayProducts(sortedByRatio.Products)

	// Sort by creation date (newest first)
	fmt.Println("\n📅 Sorted by Creation Date (Newest First):")
	sortedByDate, err := app.SortProducts(ctx, products, catalog.SortByCreatedAtDesc)
	if err != nil {
		log.Fatal(err)
	}
	displayProducts(sortedByDate.Products)

	// Sort by revenue
	fmt.Println("\n💰 Sorted by Revenue (Highest First):")
	sortedByRevenue, err := app.SortProducts(ctx, products, catalog.SortByRevenue)
	if err != nil {
		log.Fatal(err)
	}
	displayProducts(sortedByRevenue.Products)

	// Demonstrate A/B testing with batch sorting
	fmt.Println("\n🧪 A/B Testing - Batch Sort Results:")
	strategies := catalog.NewSortStrategySet(
		catalog.SortByPriceAsc,
		catalog.SortBySalesConversionRatio,
		catalog.SortByPopularity,
	)

	batchResults, err := app.BatchSort(ctx, products, strategies)
	if err != nil {
		log.Fatal(err)
	}

	for strategy, result := range batchResults.Results {
		if len(result.Products) > 0 {
			fmt.Printf("  %s: %s\n", strategy.Description(), result.Products[0].Name)
		}
	}

	// Demonstrate performance with larger dataset
	fmt.Println("\n⚡ Performance Test with 10,000 products:")
	largeDataset := generateLargeDataset(10000)
	
	start := time.Now()
	_, err = app.SortProducts(ctx, largeDataset, catalog.SortBySalesConversionRatio)
	if err != nil {
		log.Fatal(err)
	}
	duration := time.Since(start)
	fmt.Printf("Sorted 10,000 products in %v\n", duration)

	logger.Info("Demo completed successfully")
}

func initializeLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	
	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	
	return logger
}

func parseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatalf("Failed to parse date %s: %v", dateStr, err)
	}
	return date
}

func displayProducts(products []catalog.Product) {
	for _, product := range products {
		ratio := product.SalesConversionRatio()
		revenue := product.RevenueGenerated()
		fmt.Printf("  • %s - $%.2f (Sales: %d, Views: %d, Ratio: %.4f, Revenue: $%.2f)\n",
			product.Name, float64(product.Price), product.SalesCount, product.ViewsCount, ratio, revenue)
	}
}

func generateLargeDataset(size int) []catalog.Product {
	products := make([]catalog.Product, size)
	for i := 0; i < size; i++ {
		products[i] = catalog.Product{
			ID:         catalog.ProductID(i + 1),
			Name:       fmt.Sprintf("Product %d", i+1),
			Price:      catalog.Price(10 + float64(i%100)),
			CreatedAt:  time.Now().AddDate(0, 0, -i%365),
			SalesCount: i%1000 + 1,
			ViewsCount: (i%5000 + 100),
		}
	}
	return products
}
