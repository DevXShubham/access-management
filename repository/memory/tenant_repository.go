package memory

import (
	"context"
	"fmt"
	"iam/apperrors"
	"iam/models"
	"iam/repository"
	"maps"
	"slices"
	"sort"
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

func (r *TenantRepository) GetAll(ctx context.Context, limit int, cursor *models.TenantCursor) ([]models.Tenant, *models.TenantCursor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	allTenants := slices.Collect(maps.Values(r.tenants))

	sort.Slice(allTenants, func(i int, j int) bool {
		if allTenants[i].CreatedAt.Equal(allTenants[j].CreatedAt) {
			return allTenants[i].ID < allTenants[j].ID
		}

		return allTenants[i].CreatedAt.Before(allTenants[j].CreatedAt)
	})

	startIndex := 0

	if cursor != nil {
		startIndex = sort.Search(len(allTenants), func(i int) bool {
			tenant := allTenants[i]

			return tenant.CreatedAt.After(cursor.CreatedAt) ||
				(tenant.CreatedAt.Equal(cursor.CreatedAt) && tenant.ID > cursor.ID)
		})
	}

	hasMore := len(allTenants)-startIndex > limit

	endIndex := startIndex + limit

	if hasMore {
		endIndex++
	}

	endIndex = min(endIndex, len(allTenants))

	page := allTenants[startIndex:endIndex]

	var nextCursor *models.TenantCursor

	if hasMore {
		lastElement := page[len(page)-2]
		nextCursor = &models.TenantCursor{
			CreatedAt: lastElement.CreatedAt,
			ID:        lastElement.ID,
		}

		page = page[:len(page)-1]
	}

	return page, nextCursor, nil
}

var _ repository.TenantRepository = (*TenantRepository)(nil)
