package di

import (
	"poc-collaborative-filter/graph"
	"poc-collaborative-filter/internal/config"
	"poc-collaborative-filter/internal/database"
	"poc-collaborative-filter/internal/domain"
	"poc-collaborative-filter/internal/repository"
	"poc-collaborative-filter/internal/service"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module provides all application dependencies
var Module = fx.Options(
	// Configuration
	fx.Provide(config.Load),

	// Logger
	fx.Provide(NewLogger),

	// Database
	fx.Provide(database.NewDatabase),

	// Repositories
	fx.Provide(
		fx.Annotate(
			repository.NewCustomerRepository,
			fx.As(new(repository.CustomerRepository)),
		),
		fx.Annotate(
			repository.NewRecommendationRepository,
			fx.As(new(repository.RecommendationRepository)),
		),
	),

	// Strategy Selector
	fx.Provide(domain.NewStrategySelector),

	// Services
	fx.Provide(
		fx.Annotate(
			service.NewCustomerService,
			fx.As(new(service.CustomerService)),
		),
		fx.Annotate(
			service.NewRecommendationService,
			fx.As(new(service.RecommendationService)),
		),
		fx.Annotate(
			service.NewCollaborativeFilterService,
			fx.As(new(service.CollaborativeFilterService)),
		),
	),

	// GraphQL Resolver
	fx.Provide(graph.NewResolver),

	// Lifecycle hooks
	fx.Invoke(RunMigrations),
)

// NewLogger creates a new zap logger
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if cfg.Server.GinMode == "release" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return logger, nil
}

// RunMigrations runs database migrations on startup
func RunMigrations(db *gorm.DB, logger *zap.Logger) error {
	return database.AutoMigrate(db, logger)
}
