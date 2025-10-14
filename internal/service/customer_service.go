package service

import (
	"context"

	"poc-collaborative-filter/internal/domain"
	"poc-collaborative-filter/internal/repository"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	GetCustomer(ctx context.Context, id string) (*domain.Customer, error)
	GetCustomers(ctx context.Context, filter repository.CustomerFilter, limit, offset int) ([]*domain.Customer, int64, error)
	UpdateCustomer(ctx context.Context, id string, customer *domain.Customer) (*domain.Customer, error)
	DeleteCustomer(ctx context.Context, id string) error
}

type customerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	if err := s.repo.Create(ctx, customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *customerService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *customerService) GetCustomers(ctx context.Context, filter repository.CustomerFilter, limit, offset int) ([]*domain.Customer, int64, error) {
	customers, err := s.repo.FindByFilter(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return customers, count, nil
}

func (s *customerService) UpdateCustomer(ctx context.Context, id string, customer *domain.Customer) (*domain.Customer, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	customer.ID = existing.ID
	customer.CreatedAt = existing.CreatedAt

	if err := s.repo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *customerService) DeleteCustomer(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
