package repository

import (
	"context"

	"gorm.io/gorm"
	"poc-collaborative-filter/internal/domain"
)

type recommendationRepository struct {
	db *gorm.DB
}

func NewRecommendationRepository(db *gorm.DB) RecommendationRepository {
	return &recommendationRepository{db: db}
}

func (r *recommendationRepository) Create(ctx context.Context, recommendation *domain.Recommendation) error {
	return r.db.WithContext(ctx).Create(recommendation).Error
}

func (r *recommendationRepository) FindByID(ctx context.Context, id string) (*domain.Recommendation, error) {
	var recommendation domain.Recommendation
	err := r.db.WithContext(ctx).First(&recommendation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &recommendation, nil
}

func (r *recommendationRepository) FindByCustomerID(ctx context.Context, customerID string, recType *domain.RecommendationType, limit int) ([]*domain.Recommendation, error) {
	var recommendations []*domain.Recommendation
	query := r.db.WithContext(ctx).Where("customer_id = ?", customerID)

	if recType != nil {
		query = query.Where("type = ?", *recType)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("score DESC").Find(&recommendations).Error
	return recommendations, err
}

func (r *recommendationRepository) Update(ctx context.Context, recommendation *domain.Recommendation) error {
	return r.db.WithContext(ctx).Save(recommendation).Error
}

func (r *recommendationRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Recommendation{}, "id = ?", id).Error
}

func (r *recommendationRepository) DeleteByCustomerID(ctx context.Context, customerID string) error {
	return r.db.WithContext(ctx).Delete(&domain.Recommendation{}, "customer_id = ?", customerID).Error
}
