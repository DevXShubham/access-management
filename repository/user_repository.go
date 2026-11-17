package repository

import (
	"context"
	"iam/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	GetByID(ctx context.Context, id string) (models.User, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]models.User, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id string) error
}
