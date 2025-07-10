package test

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/zap"

	"product-catalog-sorting/internal/application"
	"product-catalog-sorting/internal/domain/catalog"
	"product-catalog-sorting/test/testdata"
)

// TestDebugSorting helps debug the actual sorting behavior with the challenge products
func TestDebugSorting(t *testing.T) {
	logger := zap.NewNop()
	app, err := application.New(application.Config{
		Logger:  logger,
		Context: context.Background(),
	})
	if err != nil {
		t.Fatal(err)
	}

	products := testdata.GetTestProducts()
	ctx := context.Background()

	t.Run("Debug Challenge Products", func(t *testing.T) {
		fmt.Println("\n=== Challenge Products ===")
		for _, p := range products {
			ratio := p.SalesConversionRatio()
			revenue := p.RevenueGenerated()
			fmt.Printf("ID: %d, Name: %s, Price: $%.2f, Sales: %d, Views: %d, Ratio: %.4f, Revenue: $%.2f, Created: %s\n", 
				p.ID, p.Name, float64(p.Price), p.SalesCount, p.ViewsCount, ratio, revenue, p.CreatedAt.Format("2006-01-02"))
		}
	})

	t.Run("Debug Sales Conversion Ratio", func(t *testing.T) {
		result, err := app.SortProducts(ctx, products, catalog.SortBySalesConversionRatio)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println("\n=== Sorted by Sales Conversion Ratio ===")
		for i, p := range result.Products {
			ratio := p.SalesConversionRatio()
			fmt.Printf("%d. ID: %d, Name: %s, Sales: %d, Views: %d, Ratio: %.4f (%.2f%%)\n", 
				i+1, p.ID, p.Name, p.SalesCount, p.ViewsCount, ratio, ratio*100)
		}
	})

	t.Run("Debug Price Sorting", func(t *testing.T) {
		result, err := app.SortProducts(ctx, products, catalog.SortByPriceAsc)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println("\n=== Sorted by Price Ascending ===")
		for i, p := range result.Products {
			fmt.Printf("%d. ID: %d, Name: %s, Price: $%.2f\n", 
				i+1, p.ID, p.Name, float64(p.Price))
		}
	})

	t.Run("Debug Revenue Sorting", func(t *testing.T) {
		result, err := app.SortProducts(ctx, products, catalog.SortByRevenue)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println("\n=== Sorted by Revenue ===")
		for i, p := range result.Products {
			revenue := p.RevenueGenerated()
			fmt.Printf("%d. ID: %d, Name: %s, Revenue: $%.2f (Price: $%.2f × Sales: %d)\n", 
				i+1, p.ID, p.Name, revenue, float64(p.Price), p.SalesCount)
		}
	})
}
