package memory

import (
	"context"
	"fmt"
	"iam/apperrors"
	"iam/models"
	"iam/repository"
	"strings"
	"sync"
)

type TenantRepository struct {
	mu      sync.RWMutex
	tenants map[string]models.Tenant
}

func NewTenantRepository() *TenantRepository {
	return &TenantRepository{
		tenants: make(map[string]models.Tenant),
	}
}

func (r *TenantRepository) Create(ctx context.Context, tenant models.Tenant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existingTenant := range r.tenants {
		if existingTenant.Name == tenant.Name {
			return fmt.Errorf("%w tenant '%s' already exists", apperrors.ErrTenantAlreadyExists, tenant.Name)
		}
	}

	r.tenants[tenant.ID] = tenant

	return nil
}

func (r *TenantRepository) GetByID(ctx context.Context, id string) (models.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tenant, exists := r.tenants[id]

	if !exists {
		return models.Tenant{}, fmt.Errorf("%w (id: %s)", apperrors.ErrTenantAlreadyExists, id)
	}

	return tenant, nil
}

func (r *TenantRepository) GetByName(ctx context.Context, name string) (models.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, tenant := range r.tenants {
		if tenant.Name == name {
			return tenant, nil
		}
	}

	return models.Tenant{}, fmt.Errorf("%w (name: %s)", apperrors.ErrTenantAlreadyExists, name)
}

func (r *TenantRepository) Update(ctx context.Context, tenantToUpdate models.Tenant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.tenants[tenantToUpdate.ID]

	if !exists {
		return fmt.Errorf("%w (id: %s)", apperrors.ErrTenantNotFound, tenantToUpdate.ID)
	}

	// check unique name before updating
	for _, existingTenant := range r.tenants {
		if strings.EqualFold(existingTenant.Name, tenantToUpdate.Name) &&
			existingTenant.ID != tenantToUpdate.ID {
			return fmt.Errorf(
				"%w: tenant name '%s' already exists",
				apperrors.ErrTenantAlreadyExists,
				tenantToUpdate.Name,
			)
		}
	}

	r.tenants[tenantToUpdate.ID] = tenantToUpdate

	return nil
}

func (r *TenantRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.tenants[id]

	if !exists {
		return fmt.Errorf("%w (id: %s)", apperrors.ErrTenantNotFound, id)
	}

	delete(r.tenants, id)

	return nil
}

var _ repository.TenantRepository = (*TenantRepository)(nil)
