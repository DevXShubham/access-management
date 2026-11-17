package service_test

import (
	"context"
	"errors"
	"testing"

	"iam/apperrors"
	"iam/models"
	"iam/repository"
	"iam/service"
)

type fakeTenantRepository struct {
	createdTenant models.Tenant
	updatedTenant models.Tenant
	tenant        models.Tenant
	deletedID     string

	createErr    error
	getByIDErr   error
	getByNameErr error
	updateErr    error
	deleteErr    error
}

func (f *fakeTenantRepository) Create(
	ctx context.Context,
	tenant models.Tenant,
) error {
	f.createdTenant = tenant
	return f.createErr
}

func (f *fakeTenantRepository) GetByID(
	ctx context.Context,
	id string,
) (models.Tenant, error) {
	return f.tenant, f.getByIDErr
}

func (f *fakeTenantRepository) GetByName(
	ctx context.Context,
	name string,
) (models.Tenant, error) {
	return f.tenant, f.getByNameErr
}

func (f *fakeTenantRepository) Update(
	ctx context.Context,
	tenant models.Tenant,
) error {
	f.updatedTenant = tenant
	return f.updateErr
}

func (f *fakeTenantRepository) Delete(
	ctx context.Context,
	id string,
) error {
	f.deletedID = id
	return f.deleteErr
}

var _ repository.TenantRepository = (*fakeTenantRepository)(nil)

func TestCreateTenant_Success(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}
	tenantService := service.NewTenantService(fakeRepo)

	tenant, err := tenantService.CreateTenant(
		context.Background(),
		"TestTenantA",
		"Test Tenant A",
	)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if tenant.Name != "TestTenantA" {
		t.Errorf("expected tenant name: 'TestTenantA', got: %s", tenant.Name)
	}

	if tenant.Description != "Test Tenant A" {
		t.Errorf("expected tenant description: 'Test Tenant A, got:%s ", tenant.Description)
	}

	if tenant.ID == "" {
		t.Error("expected tenant ID to be generated")
	}

	if fakeRepo.createdTenant.ID != tenant.ID {
		t.Error("expected repository to receive the generated tenant ID")
	}
}

func TestCreateTenant_EmptyName(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}
	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.CreateTenant(
		context.Background(),
		"  ",
		"Empty Tenant",
	)

	if err == nil {
		t.Fatal("expected error for empty tenant name")
	}

	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got: %v", err)
	}

	if fakeRepo.createdTenant.ID != "" {
		t.Error("repository should not be called for invalid input")
	}
}

func TestCreateTenant_DuplicateTenant(t *testing.T) {
	fakeRepo := &fakeTenantRepository{
		createErr: apperrors.ErrTenantAlreadyExists,
	}

	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.CreateTenant(
		context.Background(),
		"TestTenantB",
		"Test Tenant B",
	)

	if err == nil {
		t.Fatal("expected error for duplicate tenant")
	}

	if !errors.Is(err, apperrors.ErrTenantAlreadyExists) {
		t.Errorf("expected ErrTenantAlreadyExists, got: %v", err)
	}
}
