package repository

import (
	"context"

	"gorm.io/gorm"
	"poc-collaborative-filter/internal/domain"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customerRepository) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.WithContext(ctx).First(&customer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) FindAll(ctx context.Context, limit, offset int) ([]*domain.Customer, error) {
	var customers []*domain.Customer
	query := r.db.WithContext(ctx)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&customers).Error
	return customers, err
}

func (r *customerRepository) FindByFilter(ctx context.Context, filter CustomerFilter, limit, offset int) ([]*domain.Customer, error) {
	var customers []*domain.Customer
	query := r.db.WithContext(ctx)

	query = r.applyFilter(query, filter)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&customers).Error
	return customers, err
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Customer{}, "id = ?", id).Error
}

func (r *customerRepository) Count(ctx context.Context, filter CustomerFilter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.Customer{})
	query = r.applyFilter(query, filter)
	err := query.Count(&count).Error
	return count, err
}

func (r *customerRepository) applyFilter(query *gorm.DB, filter CustomerFilter) *gorm.DB {
	if filter.Segment != nil {
		query = query.Where("segment = ?", *filter.Segment)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.RiskProfile != nil {
		query = query.Where("risk_profile = ?", *filter.RiskProfile)
	}
	if filter.Search != nil && *filter.Search != "" {
		searchPattern := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}
	return query
}
