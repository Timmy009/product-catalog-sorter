package catalog

import (
	"fmt"
	"strings"
)

// SortStrategy represents different product sorting strategies
// Using type-safe enumeration pattern with rich behavior
type SortStrategy string

// Predefined sorting strategies
const (
	SortByPriceAsc              SortStrategy = "price_asc"
	SortByPriceDesc             SortStrategy = "price_desc"
	SortBySalesConversionRatio  SortStrategy = "sales_conversion_ratio"
	SortByCreatedAtDesc         SortStrategy = "created_at_desc"
	SortByCreatedAtAsc          SortStrategy = "created_at_asc"
	SortByPopularity            SortStrategy = "popularity"
	SortByRevenue               SortStrategy = "revenue"
	SortByName                  SortStrategy = "name"
)

// AllSortStrategies returns all available sort strategies
func AllSortStrategies() []SortStrategy {
	return []SortStrategy{
		SortByPriceAsc,
		SortByPriceDesc,
		SortBySalesConversionRatio,
		SortByCreatedAtDesc,
		SortByCreatedAtAsc,
		SortByPopularity,
		SortByRevenue,
		SortByName,
	}
}

// String returns the string representation of the sort strategy
func (s SortStrategy) String() string {
	return string(s)
}

// IsValid checks if the sort strategy is supported
func (s SortStrategy) IsValid() bool {
	for _, strategy := range AllSortStrategies() {
		if s == strategy {
			return true
		}
	}
	return false
}

// Description returns a human-readable description of the sort strategy
func (s SortStrategy) Description() string {
	switch s {
	case SortByPriceAsc:
		return "Price (Low to High)"
	case SortByPriceDesc:
		return "Price (High to Low)"
	case SortBySalesConversionRatio:
		return "Sales Conversion Ratio (Best Performers First)"
	case SortByCreatedAtDesc:
		return "Creation Date (Newest First)"
	case SortByCreatedAtAsc:
		return "Creation Date (Oldest First)"
	case SortByPopularity:
		return "Popularity (Most Viewed First)"
	case SortByRevenue:
		return "Revenue Generated (Highest First)"
	case SortByName:
		return "Name (Alphabetical)"
	default:
		return fmt.Sprintf("Unknown Strategy (%s)", s)
	}
}

// Priority returns the business priority of this sort strategy
// Higher values indicate higher business importance
func (s SortStrategy) Priority() int {
	switch s {
	case SortBySalesConversionRatio:
		return 10 // Highest priority - directly impacts revenue
	case SortByRevenue:
		return 9
	case SortByPopularity:
		return 8
	case SortByPriceAsc, SortByPriceDesc:
		return 7
	case SortByCreatedAtDesc:
		return 6
	case SortByCreatedAtAsc:
		return 5
	case SortByName:
		return 4
	default:
		return 1
	}
}

// SortStrategySet represents a collection of sort strategies with utility methods
type SortStrategySet []SortStrategy

// NewSortStrategySet creates a new set of sort strategies
func NewSortStrategySet(strategies ...SortStrategy) SortStrategySet {
	return SortStrategySet(strategies)
}

// Contains checks if the set contains a specific strategy
func (s SortStrategySet) Contains(strategy SortStrategy) bool {
	for _, existing := range s {
		if existing == strategy {
			return true
		}
	}
	return false
}

// Validate checks if all strategies in the set are valid
func (s SortStrategySet) Validate() error {
	var invalidStrategies []string
	
	for _, strategy := range s {
		if !strategy.IsValid() {
			invalidStrategies = append(invalidStrategies, string(strategy))
		}
	}
	
	if len(invalidStrategies) > 0 {
		return fmt.Errorf("invalid sort strategies: %s", 
			strings.Join(invalidStrategies, ", "))
	}
	
	return nil
}

// String returns a string representation of the strategy set
func (s SortStrategySet) String() string {
	if len(s) == 0 {
		return "SortStrategySet{empty}"
	}
	
	strategies := make([]string, len(s))
	for i, strategy := range s {
		strategies[i] = string(strategy)
	}
	
	return fmt.Sprintf("SortStrategySet{%s}", strings.Join(strategies, ", "))
}

// Len returns the number of strategies in the set
func (s SortStrategySet) Len() int {
	return len(s)
}

// ToSlice returns the underlying slice
func (s SortStrategySet) ToSlice() []SortStrategy {
	return []SortStrategy(s)
}
