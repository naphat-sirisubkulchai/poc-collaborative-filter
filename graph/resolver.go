package graph

import (
	"poc-collaborative-filter/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	customerService             service.CustomerService
	recommendationService       service.RecommendationService
	collaborativeFilterService  service.CollaborativeFilterService
}

// NewResolver creates a new resolver with injected dependencies
func NewResolver(
	customerService service.CustomerService,
	recommendationService service.RecommendationService,
	collaborativeFilterService service.CollaborativeFilterService,
) *Resolver {
	return &Resolver{
		customerService:            customerService,
		recommendationService:      recommendationService,
		collaborativeFilterService: collaborativeFilterService,
	}
}
