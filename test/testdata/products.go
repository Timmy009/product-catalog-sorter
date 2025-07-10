package testdata

import (
	"time"

	"product-catalog-sorting/internal/domain/catalog"
)

// GetTestProducts returns the exact 3 products from the code challenge
func GetTestProducts() []catalog.Product {
	return []catalog.Product{
		{
			ID:         1,
			Name:       "Alabaster Table",
			Price:      12.99,
			CreatedAt:  parseDate("2019-01-04"),
			SalesCount: 32,
			ViewsCount: 730, // 4.38% conversion ratio
		},
		{
			ID:         2,
			Name:       "Zebra Table",
			Price:      44.49,
			CreatedAt:  parseDate("2012-01-04"),
			SalesCount: 301,
			ViewsCount: 3279, // 9.18% conversion ratio
		},
		{
			ID:         3,
			Name:       "Coffee Table",
			Price:      10.00,
			CreatedAt:  parseDate("2014-05-28"),
			SalesCount: 1048,
			ViewsCount: 20123, // 5.21% conversion ratio
		},
	}
}

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic("Failed to parse date: " + dateStr)
	}
	return date
}

// GetTestProductCollection returns products as a ProductCollection
func GetTestProductCollection() catalog.ProductCollection {
	products := GetTestProducts()
	collection := make(catalog.ProductCollection, len(products))
	copy(collection, products)
	return collection
}

// ProductTestCase represents a test case for product operations
type ProductTestCase struct {
	Name        string
	Products    []catalog.Product
	Strategy    catalog.SortStrategy
	ExpectedTop string // Name of expected top product
}

// GetSortingTestCases returns predefined test cases for sorting
func GetSortingTestCases() []ProductTestCase {
	products := GetTestProducts()
	
	return []ProductTestCase{
		{
			Name:        "Price Ascending",
			Products:    products,
			Strategy:    catalog.SortByPriceAsc,
			ExpectedTop: "Coffee Table", // $10.00
		},
		{
			Name:        "Price Descending", 
			Products:    products,
			Strategy:    catalog.SortByPriceDesc,
			ExpectedTop: "Zebra Table", // $44.49
		},
		{
			Name:        "Sales Conversion Ratio",
			Products:    products,
			Strategy:    catalog.SortBySalesConversionRatio,
			ExpectedTop: "Zebra Table", // 9.18% conversion, highest
		},
		{
			Name:        "Creation Date Newest",
			Products:    products,
			Strategy:    catalog.SortByCreatedAtDesc,
			ExpectedTop: "Alabaster Table", // 2019-01-04, most recent
		},
		{
			Name:        "Creation Date Oldest",
			Products:    products,
			Strategy:    catalog.SortByCreatedAtAsc,
			ExpectedTop: "Zebra Table", // 2012-01-04, oldest
		},
		{
			Name:        "Popularity",
			Products:    products,
			Strategy:    catalog.SortByPopularity,
			ExpectedTop: "Coffee Table", // 20123 views, highest
		},
		{
			Name:        "Revenue",
			Products:    products,
			Strategy:    catalog.SortByRevenue,
			ExpectedTop: "Zebra Table", // 301 * $44.49 = $13,391.49, highest revenue
		},
		{
			Name:        "Name",
			Products:    products,
			Strategy:    catalog.SortByName,
			ExpectedTop: "Alabaster Table", // Alphabetically first
		},
	}
}
