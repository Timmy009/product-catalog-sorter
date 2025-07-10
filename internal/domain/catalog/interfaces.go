package catalog

import (
	"context"
	"time"
)

// ProductRepository defines the contract for product data persistence
type ProductRepository interface {
	// FindByID retrieves a product by its unique identifier
	FindByID(ctx context.Context, id ProductID) (*Product, error)
	
	// FindAll retrieves all products with optional filtering
	FindAll(ctx context.Context, filter ProductFilter) (ProductCollection, error)
	
	// Save persists a product (create or update)
	Save(ctx context.Context, product *Product) error
	
	// Delete removes a product by ID
	Delete(ctx context.Context, id ProductID) error
	
	// Count returns the total number of products matching the filter
	Count(ctx context.Context, filter ProductFilter) (int, error)
}

// CatalogService defines the core business operations interface
type CatalogService interface {
	// SortProducts sorts a collection using the specified strategy
	SortProducts(ctx context.Context, products ProductCollection, strategy SortStrategy) (*SortResult, error)
	
	// BatchSort performs multiple sorting operations simultaneously
	BatchSort(ctx context.Context, products ProductCollection, strategies SortStrategySet) (*BatchSortResult, error)
	
	// GetSupportedStrategies returns all available sorting strategies
	GetSupportedStrategies() SortStrategySet
	
	// ValidateProducts ensures product data integrity
	ValidateProducts(ctx context.Context, products ProductCollection) error
	
	// AnalyzePerformance provides insights on product performance
	AnalyzePerformance(ctx context.Context, products ProductCollection) (*PerformanceAnalysis, error)
}

// SortingEngine defines the contract for sorting implementations
type SortingEngine interface {
	// Sort applies sorting logic to a product collection
	Sort(ctx context.Context, products ProductCollection) (ProductCollection, error)
	
	// GetStrategy returns the sorting strategy this engine implements
	GetStrategy() SortStrategy
	
	// GetMetadata returns information about the sorting algorithm
	GetMetadata() SortingMetadata
}

// SortingFactory creates sorting engines for different strategies
type SortingFactory interface {
	// CreateEngine creates a sorting engine for the specified strategy
	CreateEngine(strategy SortStrategy) (SortingEngine, error)
	
	// GetSupportedStrategies returns all strategies this factory can create
	GetSupportedStrategies() SortStrategySet
	
	// IsSupported checks if a strategy is supported
	IsSupported(strategy SortStrategy) bool
}

// MetricsCollector defines the contract for performance monitoring
type MetricsCollector interface {
	// RecordSortOperation records metrics for a sorting operation
	RecordSortOperation(ctx context.Context, strategy SortStrategy, duration time.Duration, productCount int)
	
	// RecordBatchOperation records metrics for a batch operation
	RecordBatchOperation(ctx context.Context, strategies SortStrategySet, duration time.Duration, productCount int)
	
	// GetMetrics retrieves collected metrics
	GetMetrics(ctx context.Context) (*OperationMetrics, error)
}

// CacheManager defines the contract for caching sorted results
type CacheManager interface {
	// Get retrieves cached sort results
	Get(ctx context.Context, key CacheKey) (*SortResult, error)
	
	// Set stores sort results in cache
	Set(ctx context.Context, key CacheKey, result *SortResult, ttl time.Duration) error
	
	// Invalidate removes cached results
	Invalidate(ctx context.Context, pattern string) error
	
	// Clear removes all cached results
	Clear(ctx context.Context) error
}

// EventPublisher defines the contract for publishing domain events
type EventPublisher interface {
	// PublishSortCompleted publishes when a sort operation completes
	PublishSortCompleted(ctx context.Context, event SortCompletedEvent) error
	
	// PublishBatchCompleted publishes when a batch operation completes
	PublishBatchCompleted(ctx context.Context, event BatchCompletedEvent) error
	
	// PublishPerformanceAlert publishes performance-related alerts
	PublishPerformanceAlert(ctx context.Context, event PerformanceAlertEvent) error
}

// Supporting types for interfaces

// SortingMetadata provides information about a sorting algorithm
type SortingMetadata struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	TimeComplexity string      `json:"time_complexity"`
	SpaceComplexity string     `json:"space_complexity"`
	StableSort   bool          `json:"stable_sort"`
	InPlace      bool          `json:"in_place"`
}

// PerformanceAnalysis provides insights on product performance
type PerformanceAnalysis struct {
	TotalProducts       int                    `json:"total_products"`
	HighPerformers      ProductCollection      `json:"high_performers"`
	LowPerformers       ProductCollection      `json:"low_performers"`
	AverageConversion   float64               `json:"average_conversion"`
	TotalRevenue        float64               `json:"total_revenue"`
	TopCategories       []CategoryMetrics      `json:"top_categories"`
	PerformanceMetrics  map[string]interface{} `json:"performance_metrics"`
	GeneratedAt         time.Time             `json:"generated_at"`
}

// CategoryMetrics represents performance metrics for a product category
type CategoryMetrics struct {
	Category    string  `json:"category"`
	ProductCount int    `json:"product_count"`
	TotalRevenue float64 `json:"total_revenue"`
	AvgConversion float64 `json:"avg_conversion"`
}

// OperationMetrics contains performance metrics for operations
type OperationMetrics struct {
	TotalOperations    int64                    `json:"total_operations"`
	AverageLatency     time.Duration           `json:"average_latency"`
	OperationsByStrategy map[SortStrategy]int64 `json:"operations_by_strategy"`
	ErrorRate          float64                 `json:"error_rate"`
	ThroughputPerSecond float64                `json:"throughput_per_second"`
	CollectedAt        time.Time               `json:"collected_at"`
}

// CacheKey represents a cache key for sorted results
type CacheKey struct {
	ProductHash  string       `json:"product_hash"`
	Strategy     SortStrategy `json:"strategy"`
	Version      string       `json:"version"`
}

// Domain Events

// SortCompletedEvent is published when a sort operation completes
type SortCompletedEvent struct {
	Strategy      SortStrategy  `json:"strategy"`
	ProductCount  int          `json:"product_count"`
	Duration      time.Duration `json:"duration"`
	Success       bool         `json:"success"`
	ErrorMessage  string       `json:"error_message,omitempty"`
	Timestamp     time.Time    `json:"timestamp"`
}

// BatchCompletedEvent is published when a batch operation completes
type BatchCompletedEvent struct {
	Strategies    SortStrategySet `json:"strategies"`
	ProductCount  int            `json:"product_count"`
	Duration      time.Duration  `json:"duration"`
	SuccessCount  int            `json:"success_count"`
	ErrorCount    int            `json:"error_count"`
	Timestamp     time.Time      `json:"timestamp"`
}

// PerformanceAlertEvent is published for performance-related alerts
type PerformanceAlertEvent struct {
	AlertType    string                 `json:"alert_type"`
	Severity     string                 `json:"severity"`
	Message      string                 `json:"message"`
	Metadata     map[string]interface{} `json:"metadata"`
	Timestamp    time.Time              `json:"timestamp"`
}
