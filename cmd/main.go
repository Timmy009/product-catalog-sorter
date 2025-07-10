package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"product-catalog-sorting/internal/application"
	"product-catalog-sorting/internal/domain/catalog"
	"product-catalog-sorting/pkg/version"
)

// Build information set by ldflags
var (
	Version    = "dev"
	CommitHash = "unknown"
	BuildTime  = "unknown"
	GoVersion  = "unknown"
)

func main() {
	// Initialize build info
	buildInfo := version.BuildInfo{
		Version:    Version,
		CommitHash: CommitHash,
		BuildTime:  BuildTime,
		GoVersion:  GoVersion,
	}

	// Initialize logger
	logger := initializeLogger()
	defer func() {
		if err := logger.Sync(); err != nil {
			// Ignore sync errors on stdout/stderr
		}
	}()

	// Log startup information
	logger.Info("Starting Product Catalog Sorting System",
		zap.String("version", buildInfo.Version),
		zap.String("commit", buildInfo.CommitHash),
		zap.String("build_time", buildInfo.BuildTime),
		zap.String("go_version", buildInfo.GoVersion),
		zap.Int("cpu_count", runtime.NumCPU()),
		zap.String("os", runtime.GOOS),
		zap.String("arch", runtime.GOARCH),
	)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown
	setupGracefulShutdown(ctx, cancel, logger)

	// Initialize application
	app, err := application.New(application.Config{
		Logger:  logger,
		Context: ctx,
	})
	if err != nil {
		logger.Fatal("Failed to initialize application", zap.Error(err))
	}

	// Run demonstration
	if err := runDemonstration(ctx, app, logger); err != nil {
		logger.Error("Demonstration failed", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("Application completed successfully")
}

// initializeLogger creates a production-ready structured logger
func initializeLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logger, err := config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	return logger
}

// setupGracefulShutdown configures graceful shutdown handling
func setupGracefulShutdown(ctx context.Context, cancel context.CancelFunc, logger *zap.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigChan:
			logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
			cancel()
		case <-ctx.Done():
			return
		}
	}()
}

// runDemonstration executes the main application demonstration
func runDemonstration(ctx context.Context, app *application.Application, logger *zap.Logger) error {
	logger.Info("Starting product catalog sorting demonstration")

	// Load sample data - using the exact products from the code challenge
	products, err := loadSampleProducts()
	if err != nil {
		return fmt.Errorf("failed to load sample products: %w", err)
	}

	logger.Info("Loaded sample products", zap.Int("count", len(products)))

	// Demonstrate sorting strategies
	if err := demonstrateSortingStrategies(ctx, app, products, logger); err != nil {
		return fmt.Errorf("sorting demonstration failed: %w", err)
	}

	// Demonstrate A/B testing capabilities
	if err := demonstrateABTesting(ctx, app, products, logger); err != nil {
		return fmt.Errorf("A/B testing demonstration failed: %w", err)
	}

	return nil
}

// loadSampleProducts creates the exact products from the code challenge
func loadSampleProducts() ([]catalog.Product, error) {
	products := []catalog.Product{
		{
			ID:         catalog.ProductID(1),
			Name:       "Alabaster Table",
			Price:      catalog.Price(12.99),
			CreatedAt:  parseDate("2019-01-04"),
			SalesCount: 32,
			ViewsCount: 730,
		},
		{
			ID:         catalog.ProductID(2),
			Name:       "Zebra Table",
			Price:      catalog.Price(44.49),
			CreatedAt:  parseDate("2012-01-04"),
			SalesCount: 301,
			ViewsCount: 3279,
		},
		{
			ID:         catalog.ProductID(3),
			Name:       "Coffee Table",
			Price:      catalog.Price(10.00),
			CreatedAt:  parseDate("2014-05-28"),
			SalesCount: 1048,
			ViewsCount: 20123,
		},
	}

	// Validate all products
	for i, product := range products {
		if err := product.Validate(); err != nil {
			return nil, fmt.Errorf("invalid product at index %d: %w", i, err)
		}
	}

	return products, nil
}

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatalf("Failed to parse date %s: %v", dateStr, err)
	}
	return date
}

// demonstrateSortingStrategies shows different sorting capabilities
func demonstrateSortingStrategies(ctx context.Context, app *application.Application, products []catalog.Product, logger *zap.Logger) error {
	logger.Info("Demonstrating sorting strategies")

	strategies := []catalog.SortStrategy{
		catalog.SortByPriceAsc,
		catalog.SortByPriceDesc,
		catalog.SortBySalesConversionRatio,
		catalog.SortByCreatedAtDesc,
		catalog.SortByPopularity,
		catalog.SortByRevenue,
	}

	for _, strategy := range strategies {
		start := time.Now()
		
		result, err := app.SortProducts(ctx, products, strategy)
		if err != nil {
			return fmt.Errorf("failed to sort by %s: %w", strategy, err)
		}

		duration := time.Since(start)
		
		logger.Info("Sorting completed",
			zap.String("strategy", string(strategy)),
			zap.Duration("duration", duration),
			zap.Int("product_count", len(result.Products)),
		)

		// Display all results for the 3 products
		fmt.Printf("\n🎯 %s:\n", strategy.Description())
		for i, product := range result.Products {
			fmt.Printf("  %d. %s - $%.2f (Sales: %d, Views: %d, Ratio: %.4f, Revenue: $%.2f)\n",
				i+1, product.Name, float64(product.Price), product.SalesCount, product.ViewsCount,
				product.SalesConversionRatio(), product.RevenueGenerated())
		}
	}

	return nil
}

// demonstrateABTesting shows batch sorting for A/B testing scenarios
func demonstrateABTesting(ctx context.Context, app *application.Application, products []catalog.Product, logger *zap.Logger) error {
	logger.Info("Demonstrating A/B testing capabilities")

	strategies := catalog.NewSortStrategySet(
		catalog.SortByPriceAsc,
		catalog.SortBySalesConversionRatio,
		catalog.SortByPopularity,
	)

	start := time.Now()
	results, err := app.BatchSort(ctx, products, strategies)
	if err != nil {
		return fmt.Errorf("batch sort failed: %w", err)
	}
	duration := time.Since(start)

	logger.Info("A/B testing batch sort completed",
		zap.Duration("duration", duration),
		zap.Int("strategy_count", len(strategies)),
		zap.Int("product_count", len(products)),
	)

	fmt.Printf("\n🧪 A/B Testing Results (Top product for each strategy):\n")
	for strategy, result := range results.Results {
		if len(result.Products) > 0 {
			top := result.Products[0]
			fmt.Printf("  %s: %s - $%.2f\n", 
				strategy.Description(), top.Name, float64(top.Price))
		}
	}

	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
