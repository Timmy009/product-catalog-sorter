package application

import (
	"context"

	"go.uber.org/zap"

	"product-catalog-sorting/internal/domain/catalog"
	"product-catalog-sorting/internal/infrastructure/sorting"
)

// Config holds the application configuration
type Config struct {
	Logger  *zap.Logger
	Context context.Context
}

// Application represents the main application
type Application struct {
	catalogService catalog.Service
	logger         *zap.Logger
}

// New creates a new application instance
func New(config Config) (*Application, error) {
	// Create sorter factory
	sorterFactory := sorting.NewSorterFactory()

	// Create catalog service
	catalogService := catalog.NewService(sorterFactory, config.Logger)

	return &Application{
		catalogService: catalogService,
		logger:         config.Logger,
	}, nil
}

// SortProducts sorts products using the specified strategy
func (a *Application) SortProducts(ctx context.Context, products []catalog.Product, strategy catalog.SortStrategy) (*catalog.SortResult, error) {
	productCollection := catalog.ProductCollection(products)
	return a.catalogService.SortProducts(ctx, productCollection, strategy)
}

// BatchSort sorts products using multiple strategies
func (a *Application) BatchSort(ctx context.Context, products []catalog.Product, strategies catalog.SortStrategySet) (*catalog.BatchSortResult, error) {
	productCollection := catalog.ProductCollection(products)
	return a.catalogService.BatchSort(ctx, productCollection, strategies)
}

// GetSupportedStrategies returns all supported sorting strategies
func (a *Application) GetSupportedStrategies() catalog.SortStrategySet {
	return a.catalogService.GetSupportedStrategies()
}

// ValidateProducts validates a collection of products
func (a *Application) ValidateProducts(ctx context.Context, products []catalog.Product) error {
	productCollection := catalog.ProductCollection(products)
	return a.catalogService.ValidateProducts(ctx, productCollection)
}
