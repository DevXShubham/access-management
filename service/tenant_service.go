package service

import (
	"context"
	"fmt"
	"iam/apperrors"
	"iam/models"
	"iam/repository"
	"strings"

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
) error {

	tenant.ID = strings.TrimSpace(tenant.ID)
	tenant.Name = strings.TrimSpace(tenant.Name)

	if tenant.ID == "" {
		return fmt.Errorf("%w: tenant id cannot be empty", apperrors.ErrInvalidInput)
	}

	if tenant.Name == "" {
		return fmt.Errorf("%w: tenant name cannot be empty", apperrors.ErrInvalidInput)
	}

	err := s.repo.Update(ctx, tenant)

	if err != nil {
		return fmt.Errorf("updating tenant: %w", err)
	}

	return nil
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
