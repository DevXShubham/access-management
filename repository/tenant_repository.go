package repository

import (
	"context"
	"iam/models"
)

type TenantRepository interface {
	Create(ctx context.Context, tenant models.Tenant) error
	GetByID(ctx context.Context, id string) (models.Tenant, error)
	GetByName(ctx context.Context, name string) (models.Tenant, error)
	Update(ctx context.Context, tenant models.Tenant) (models.Tenant, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, limit int, cursor *models.TenantCursor) ([]models.Tenant, *models.TenantCursor, error)
}
