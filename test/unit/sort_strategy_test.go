package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"product-catalog-sorting/internal/domain/catalog"
)

func TestSortStrategy_String(t *testing.T) {
	tests := []struct {
		strategy catalog.SortStrategy
		expected string
	}{
		{catalog.SortByPriceAsc, "price_asc"},
		{catalog.SortByPriceDesc, "price_desc"},
		{catalog.SortBySalesConversionRatio, "sales_conversion_ratio"},
		{catalog.SortByCreatedAtDesc, "created_at_desc"},
		{catalog.SortByCreatedAtAsc, "created_at_asc"},
		{catalog.SortByPopularity, "popularity"},
		{catalog.SortByRevenue, "revenue"},
		{catalog.SortByName, "name"},
	}

	for _, tt := range tests {
		t.Run(string(tt.strategy), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.strategy.String())
		})
	}
}

func TestSortStrategy_IsValid(t *testing.T) {
	validStrategies := []catalog.SortStrategy{
		catalog.SortByPriceAsc,
		catalog.SortByPriceDesc,
		catalog.SortBySalesConversionRatio,
		catalog.SortByCreatedAtDesc,
		catalog.SortByCreatedAtAsc,
		catalog.SortByPopularity,
		catalog.SortByRevenue,
		catalog.SortByName,
	}

	for _, strategy := range validStrategies {
		t.Run(string(strategy), func(t *testing.T) {
			assert.True(t, strategy.IsValid())
		})
	}

	invalidStrategies := []catalog.SortStrategy{
		"invalid_strategy",
		"",
		"random_text",
	}

	for _, strategy := range invalidStrategies {
		t.Run(string(strategy), func(t *testing.T) {
			assert.False(t, strategy.IsValid())
		})
	}
}

func TestSortStrategy_Description(t *testing.T) {
	tests := []struct {
		strategy    catalog.SortStrategy
		description string
	}{
		{catalog.SortByPriceAsc, "Price (Low to High)"},
		{catalog.SortByPriceDesc, "Price (High to Low)"},
		{catalog.SortBySalesConversionRatio, "Sales Conversion Ratio (Best Performers First)"},
		{catalog.SortByCreatedAtDesc, "Creation Date (Newest First)"},
		{catalog.SortByCreatedAtAsc, "Creation Date (Oldest First)"},
		{catalog.SortByPopularity, "Popularity (Most Viewed First)"},
		{catalog.SortByRevenue, "Revenue Generated (Highest First)"},
		{catalog.SortByName, "Name (Alphabetical)"},
		{catalog.SortStrategy("unknown"), "Unknown Strategy (unknown)"},
	}

	for _, tt := range tests {
		t.Run(string(tt.strategy), func(t *testing.T) {
			assert.Equal(t, tt.description, tt.strategy.Description())
		})
	}
}

func TestSortStrategy_Priority(t *testing.T) {
	tests := []struct {
		strategy catalog.SortStrategy
		priority int
	}{
		{catalog.SortBySalesConversionRatio, 10},
		{catalog.SortByRevenue, 9},
		{catalog.SortByPopularity, 8},
		{catalog.SortByPriceAsc, 7},
		{catalog.SortByPriceDesc, 7},
		{catalog.SortByCreatedAtDesc, 6},
		{catalog.SortByCreatedAtAsc, 5},
		{catalog.SortByName, 4},
		{catalog.SortStrategy("unknown"), 1},
	}

	for _, tt := range tests {
		t.Run(string(tt.strategy), func(t *testing.T) {
			assert.Equal(t, tt.priority, tt.strategy.Priority())
		})
	}
}

func TestAllSortStrategies(t *testing.T) {
	strategies := catalog.AllSortStrategies()

	// Check that we have all expected strategies
	expectedStrategies := []catalog.SortStrategy{
		catalog.SortByPriceAsc,
		catalog.SortByPriceDesc,
		catalog.SortBySalesConversionRatio,
		catalog.SortByCreatedAtDesc,
		catalog.SortByCreatedAtAsc,
		catalog.SortByPopularity,
		catalog.SortByRevenue,
		catalog.SortByName,
	}

	assert.Len(t, strategies, len(expectedStrategies))

	for _, expected := range expectedStrategies {
		assert.Contains(t, strategies, expected)
	}

	// Verify all returned strategies are valid
	for _, strategy := range strategies {
		assert.True(t, strategy.IsValid())
	}
}

func TestSortStrategySet_NewSortStrategySet(t *testing.T) {
	strategies := catalog.NewSortStrategySet(
		catalog.SortByPriceAsc,
		catalog.SortBySalesConversionRatio,
		catalog.SortByPopularity,
	)

	assert.Len(t, strategies, 3)
	assert.Contains(t, strategies.ToSlice(), catalog.SortByPriceAsc)
	assert.Contains(t, strategies.ToSlice(), catalog.SortBySalesConversionRatio)
	assert.Contains(t, strategies.ToSlice(), catalog.SortByPopularity)
}

func TestSortStrategySet_Contains(t *testing.T) {
	strategies := catalog.NewSortStrategySet(
		catalog.SortByPriceAsc,
		catalog.SortBySalesConversionRatio,
	)

	assert.True(t, strategies.Contains(catalog.SortByPriceAsc))
	assert.True(t, strategies.Contains(catalog.SortBySalesConversionRatio))
	assert.False(t, strategies.Contains(catalog.SortByPopularity))
	assert.False(t, strategies.Contains(catalog.SortByName))
}

func TestSortStrategySet_Validate(t *testing.T) {
	t.Run("Valid strategies", func(t *testing.T) {
		strategies := catalog.NewSortStrategySet(
			catalog.SortByPriceAsc,
			catalog.SortBySalesConversionRatio,
		)

		err := strategies.Validate()
		assert.NoError(t, err)
	})

	t.Run("Invalid strategies", func(t *testing.T) {
		strategies := catalog.SortStrategySet{
			catalog.SortByPriceAsc,
			catalog.SortStrategy("invalid"),
			catalog.SortBySalesConversionRatio,
		}

		err := strategies.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid sort strategies")
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("Empty strategies", func(t *testing.T) {
		strategies := catalog.SortStrategySet{}

		err := strategies.Validate()
		assert.NoError(t, err)
	})
}

func TestSortStrategySet_String(t *testing.T) {
	t.Run("Non-empty set", func(t *testing.T) {
		strategies := catalog.NewSortStrategySet(
			catalog.SortByPriceAsc,
			catalog.SortBySalesConversionRatio,
		)

		str := strategies.String()
		assert.Contains(t, str, "SortStrategySet{")
		assert.Contains(t, str, "price_asc")
		assert.Contains(t, str, "sales_conversion_ratio")
	})

	t.Run("Empty set", func(t *testing.T) {
		strategies := catalog.SortStrategySet{}

		str := strategies.String()
		assert.Equal(t, "SortStrategySet{empty}", str)
	})
}

func TestSortStrategySet_Len(t *testing.T) {
	strategies := catalog.NewSortStrategySet(
		catalog.SortByPriceAsc,
		catalog.SortBySalesConversionRatio,
		catalog.SortByPopularity,
	)

	assert.Equal(t, 3, strategies.Len())

	emptyStrategies := catalog.SortStrategySet{}
	assert.Equal(t, 0, emptyStrategies.Len())
}

func TestSortStrategySet_ToSlice(t *testing.T) {
	originalStrategies := []catalog.SortStrategy{
		catalog.SortByPriceAsc,
		catalog.SortBySalesConversionRatio,
		catalog.SortByPopularity,
	}

	strategySet := catalog.NewSortStrategySet(originalStrategies...)
	slice := strategySet.ToSlice()

	assert.Equal(t, len(originalStrategies), len(slice))
	for _, strategy := range originalStrategies {
		assert.Contains(t, slice, strategy)
	}
}

// Edge cases and error conditions
func TestSortStrategy_EdgeCases(t *testing.T) {
	t.Run("Empty string strategy", func(t *testing.T) {
		strategy := catalog.SortStrategy("")
		assert.False(t, strategy.IsValid())
		assert.Equal(t, "Unknown Strategy ()", strategy.Description())
		assert.Equal(t, 1, strategy.Priority())
	})

	t.Run("Very long strategy name", func(t *testing.T) {
		longName := string(make([]byte, 1000))
		strategy := catalog.SortStrategy(longName)
		assert.False(t, strategy.IsValid())
	})
}

func TestSortStrategySet_EdgeCases(t *testing.T) {
	t.Run("Nil slice conversion", func(t *testing.T) {
		var strategies catalog.SortStrategySet
		slice := strategies.ToSlice()
		assert.Nil(t, slice)
	})

	t.Run("Duplicate strategies", func(t *testing.T) {
		strategies := catalog.NewSortStrategySet(
			catalog.SortByPriceAsc,
			catalog.SortByPriceAsc, // duplicate
			catalog.SortBySalesConversionRatio,
		)

		assert.Equal(t, 3, strategies.Len()) // duplicates are allowed
		assert.True(t, strategies.Contains(catalog.SortByPriceAsc))
	})
}

// Benchmark tests
func BenchmarkSortStrategy_IsValid(b *testing.B) {
	strategy := catalog.SortByPriceAsc

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strategy.IsValid()
	}
}

func BenchmarkSortStrategy_Description(b *testing.B) {
	strategy := catalog.SortBySalesConversionRatio

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strategy.Description()
	}
}

func BenchmarkSortStrategySet_Contains(b *testing.B) {
	strategies := catalog.NewSortStrategySet(catalog.AllSortStrategies()...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strategies.Contains(catalog.SortByPriceAsc)
	}
}

func BenchmarkSortStrategySet_Validate(b *testing.B) {
	strategies := catalog.NewSortStrategySet(catalog.AllSortStrategies()...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strategies.Validate()
	}
}
