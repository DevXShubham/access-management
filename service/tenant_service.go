package service

import (
	"context"
	"fmt"
	"iam/apperrors"
	"iam/models"
	"iam/repository"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TenantService struct {
	repo repository.TenantRepository
}

func NewTenantService(repo repository.TenantRepository) *TenantService {
	return &TenantService{
		repo: repo,
	}
}

func (s *TenantService) CreateTenant(
	ctx context.Context,
	name string,
	description string,
) (models.Tenant, error) {

	if strings.TrimSpace(name) == "" {
		return models.Tenant{}, fmt.Errorf("%w: tenant name cannot be empty", apperrors.ErrInvalidInput)
	}

	newTenant := models.Tenant{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now().UTC(),
	}

	err := s.repo.Create(ctx, newTenant)

	if err != nil {
		return models.Tenant{}, fmt.Errorf("creating tenant: %w", err)
	}

	return newTenant, nil
}

func (s *TenantService) GetTenantByID(
	ctx context.Context,
	id string,
) (models.Tenant, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return models.Tenant{}, fmt.Errorf("%w: tenant id cannot be empty", apperrors.ErrInvalidInput)
	}

	tenant, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return models.Tenant{}, fmt.Errorf("fetching tenant by id: %w", err)
	}

	return tenant, nil
}

func (s *TenantService) GetTenantByName(
	ctx context.Context,
	name string,
) (models.Tenant, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.Tenant{}, fmt.Errorf("%w: tenant name cannot be empty", apperrors.ErrInvalidInput)
	}

	tenant, err := s.repo.GetByName(ctx, name)

	if err != nil {
		return models.Tenant{}, fmt.Errorf("fetching tenant by name: %w", err)
	}

	return tenant, nil
}

func (s *TenantService) UpdateTenant(
	ctx context.Context,
	tenant models.Tenant,
) (models.Tenant, error) {

	tenant.ID = strings.TrimSpace(tenant.ID)
	tenant.Name = strings.TrimSpace(tenant.Name)

	if tenant.ID == "" {
		return models.Tenant{}, fmt.Errorf("%w: tenant id cannot be empty", apperrors.ErrInvalidInput)
	}

	if tenant.Name == "" {
		return models.Tenant{}, fmt.Errorf("%w: tenant name cannot be empty", apperrors.ErrInvalidInput)
	}

	updatedTenant, err := s.repo.Update(ctx, tenant)

	if err != nil {
		return models.Tenant{}, fmt.Errorf("updating tenant: %w", err)
	}

	return updatedTenant, nil
}

func (s *TenantService) DeleteTenant(
	ctx context.Context,
	id string,
) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: tenant id cannot be empty", apperrors.ErrInvalidInput)
	}

	err := s.repo.Delete(ctx, id)

	if err != nil {
		return fmt.Errorf("deleting tenant: %w", err)
	}

	return nil
}

func (s *TenantService) GetAllTenants(
	ctx context.Context,
	limit int,
	cursor *models.TenantCursor,
) ([]models.Tenant, *models.TenantCursor, error) {

	if limit <= 0 {
		return nil, nil, fmt.Errorf("%w: limit must be greater than 0", apperrors.ErrInvalidInput)
	}

	tenants, nextCursor, err := s.repo.GetAll(ctx, limit, cursor)

	if err != nil {
		return nil, nil, fmt.Errorf("retrieving tenants: %w", err)
	}

	return tenants, nextCursor, nil
}
