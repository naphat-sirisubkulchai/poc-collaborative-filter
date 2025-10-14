package service

import (
	"context"

	"poc-collaborative-filter/internal/domain"
	"poc-collaborative-filter/internal/repository"
)

type RecommendationService interface {
	GetRecommendations(ctx context.Context, customerID string, recType *domain.RecommendationType, limit int) ([]*domain.Recommendation, error)
	DismissRecommendation(ctx context.Context, id string) (*domain.Recommendation, error)
	AcceptRecommendation(ctx context.Context, id string) (*domain.Recommendation, error)
}

type recommendationService struct {
	repo repository.RecommendationRepository
}

func NewRecommendationService(repo repository.RecommendationRepository) RecommendationService {
	return &recommendationService{repo: repo}
}

func (s *recommendationService) GetRecommendations(ctx context.Context, customerID string, recType *domain.RecommendationType, limit int) ([]*domain.Recommendation, error) {
	return s.repo.FindByCustomerID(ctx, customerID, recType, limit)
}

func (s *recommendationService) DismissRecommendation(ctx context.Context, id string) (*domain.Recommendation, error) {
	rec, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	rec.Status = domain.RecommendationStatusDismissed
	if err := s.repo.Update(ctx, rec); err != nil {
		return nil, err
	}

	return rec, nil
}

func (s *recommendationService) AcceptRecommendation(ctx context.Context, id string) (*domain.Recommendation, error) {
	rec, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	rec.Status = domain.RecommendationStatusAccepted
	if err := s.repo.Update(ctx, rec); err != nil {
		return nil, err
	}

	return rec, nil
}
