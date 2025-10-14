package repository

import (
	"context"

	"poc-collaborative-filter/internal/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	FindByID(ctx context.Context, id string) (*domain.Customer, error)
	FindAll(ctx context.Context, limit, offset int) ([]*domain.Customer, error)
	FindByFilter(ctx context.Context, filter CustomerFilter, limit, offset int) ([]*domain.Customer, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter CustomerFilter) (int64, error)
}

type RecommendationRepository interface {
	Create(ctx context.Context, recommendation *domain.Recommendation) error
	FindByID(ctx context.Context, id string) (*domain.Recommendation, error)
	FindByCustomerID(ctx context.Context, customerID string, recType *domain.RecommendationType, limit int) ([]*domain.Recommendation, error)
	Update(ctx context.Context, recommendation *domain.Recommendation) error
	Delete(ctx context.Context, id string) error
	DeleteByCustomerID(ctx context.Context, customerID string) error
}

type CustomerFilter struct {
	Segment     *domain.CustomerSegment
	Status      *domain.CustomerStatus
	RiskProfile *domain.CustomerRiskProfile
	Search      *string
}
