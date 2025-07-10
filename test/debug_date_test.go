package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"product-catalog-sorting/internal/domain/catalog"
	"product-catalog-sorting/internal/infrastructure/sorting"
)

// TestDebugDateSorting helps debug the actual date sorting behavior
func TestDebugDateSorting(t *testing.T) {
	now := time.Now()
	products := catalog.ProductCollection{
		{ID: 1, Name: "Old", CreatedAt: now.AddDate(-1, 0, 0)},      // 1 year ago
		{ID: 2, Name: "New", CreatedAt: now},                        // now
		{ID: 3, Name: "Medium", CreatedAt: now.AddDate(0, -6, 0)},   // 6 months ago
		{ID: 4, Name: "Same Time 1", CreatedAt: now.AddDate(0, -3, -10)}, // 3 months 10 days ago
		{ID: 5, Name: "Same Time 2", CreatedAt: now.AddDate(0, -3, -20)}, // 3 months 20 days ago
	}

	ctx := context.Background()

	t.Run("Debug Dates", func(t *testing.T) {
		fmt.Println("\n=== Original Products with Dates ===")
		for _, p := range products {
			fmt.Printf("ID: %d, Name: %s, Date: %s\n", p.ID, p.Name, p.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	})

	t.Run("Debug Descending Sort", func(t *testing.T) {
		sorter := sorting.NewCreatedAtSorter(false) // newest first
		sorted, err := sorter.Sort(ctx, products)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println("\n=== Sorted Descending (Newest First) ===")
		for i, p := range sorted {
			fmt.Printf("%d. ID: %d, Name: %s, Date: %s\n", i+1, p.ID, p.Name, p.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	})

	t.Run("Debug Ascending Sort", func(t *testing.T) {
		sorter := sorting.NewCreatedAtSorter(true) // oldest first
		sorted, err := sorter.Sort(ctx, products)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println("\n=== Sorted Ascending (Oldest First) ===")
		for i, p := range sorted {
			fmt.Printf("%d. ID: %d, Name: %s, Date: %s\n", i+1, p.ID, p.Name, p.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	})
}
