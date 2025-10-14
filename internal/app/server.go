package app

import (
	"context"

	"poc-collaborative-filter/graph"
	"poc-collaborative-filter/graph/generated"
	"poc-collaborative-filter/internal/config"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewGraphQLServer creates a new GraphQL server
func NewGraphQLServer(resolver *graph.Resolver) *handler.Server {
	c := generated.Config{Resolvers: resolver}
	return handler.NewDefaultServer(generated.NewExecutableSchema(c))
}

// NewGinRouter creates a new Gin router with GraphQL endpoints
func NewGinRouter(cfg *config.Config, graphqlServer *handler.Server, logger *zap.Logger) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	router := gin.Default()

	// Apply middleware
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// GraphQL playground (development only)
	if cfg.Server.GinMode != "release" {
		router.GET("/", gin.WrapH(playground.Handler("GraphQL Playground", "/query")))
		logger.Info("GraphQL Playground enabled at http://localhost:" + cfg.Server.Port)
	}

	// GraphQL endpoint
	router.POST("/query", gin.WrapH(graphqlServer))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "collaborative-filter-service",
		})
	})

	return router
}

// CORSMiddleware adds CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RunServer starts the HTTP server
func RunServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	router *gin.Engine,
	logger *zap.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting collaborative filter service",
				zap.String("mode", cfg.Server.GinMode),
				zap.String("port", cfg.Server.Port),
			)

			// Start server in a goroutine
			go func() {
				addr := cfg.Server.GetServerAddr()
				logger.Info("Server listening", zap.String("address", addr))
				if err := router.Run(addr); err != nil {
					logger.Fatal("Failed to start server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down collaborative filter service...")
			return nil
		},
	})
}

// Module provides app-related dependencies
var Module = fx.Options(
	fx.Provide(
		NewGraphQLServer,
		NewGinRouter,
	),
	fx.Invoke(RunServer),
)
