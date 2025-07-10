package utils

import (
	"fmt"
	"log"
	"time"
)

// ParseDate parses a date string in YYYY-MM-DD format
// Utility function for consistent date parsing across the application
func ParseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("Warning: Failed to parse date %s, using current time: %v", dateStr, err)
		return time.Now()
	}
	return date
}

// FormatDate formats a time.Time to YYYY-MM-DD string
// Consistent date formatting for output
func FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

// FormatPrice formats a price with proper currency formatting
// Utility for consistent price display
func FormatPrice(price float64) string {
	return fmt.Sprintf("$%.2f", price)
}

// CalculatePercentage calculates percentage with proper handling of zero denominators
// Safe percentage calculation utility
func CalculatePercentage(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0.0
	}
	return (float64(numerator) / float64(denominator)) * 100.0
}

// FormatPercentage formats a decimal as a percentage string
// Consistent percentage formatting
func FormatPercentage(ratio float64) string {
	return fmt.Sprintf("%.2f%%", ratio*100)
}

// SafeDivide performs division with zero-denominator protection
// Prevents division by zero errors
func SafeDivide(numerator, denominator float64) float64 {
	if denominator == 0 {
		return 0.0
	}
	return numerator / denominator
}

// Clamp constrains a value between min and max bounds
// Utility for value range validation
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// IsValidPrice checks if a price value is valid for business logic
// Business rule validation for prices
func IsValidPrice(price float64) bool {
	return price >= 0 && price <= 999999.99 // Reasonable upper bound
}

// IsValidCount checks if a count value is valid
// Validation for sales/views counts
func IsValidCount(count int) bool {
	return count >= 0 && count <= 1000000000 // Reasonable upper bound
}

// TimeAgo returns a human-readable "time ago" string
// User-friendly time formatting
func TimeAgo(t time.Time) string {
	duration := time.Since(t)
	
	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case duration < 30*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return FormatDate(t)
	}
}
