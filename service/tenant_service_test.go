package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"iam/apperrors"
	"iam/models"
	"iam/repository"
	"iam/service"

	"github.com/google/uuid"
)

type fakeTenantRepository struct {
	createdTenant   models.Tenant
	updatedTenant   models.Tenant
	tenant          models.Tenant
	deletedID       string
	gotGetByID      string
	gotGetByName    string
	gotGetAllLimit  int
	gotGetAllCursor *models.TenantCursor
	tenants         []models.Tenant
	nextCursor      *models.TenantCursor

	createErr    error
	getByIDErr   error
	getByNameErr error
	updateErr    error
	deleteErr    error
	getAllErr    error
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
	f.gotGetByID = id
	return f.tenant, f.getByIDErr
}

func (f *fakeTenantRepository) GetByName(
	ctx context.Context,
	name string,
) (models.Tenant, error) {
	f.gotGetByName = name
	return f.tenant, f.getByNameErr
}

func (f *fakeTenantRepository) Update(
	ctx context.Context,
	tenant models.Tenant,
) (models.Tenant, error) {
	f.updatedTenant = tenant
	return f.updatedTenant, f.updateErr
}

func (f *fakeTenantRepository) Delete(
	ctx context.Context,
	id string,
) error {
	f.deletedID = id
	return f.deleteErr
}

func (f *fakeTenantRepository) GetAll(
	ctx context.Context,
	limit int,
	cursor *models.TenantCursor,
) ([]models.Tenant, *models.TenantCursor, error) {
	f.gotGetAllLimit = limit
	f.gotGetAllCursor = cursor
	return f.tenants, f.nextCursor, f.getAllErr
}

var _ repository.TenantRepository = (*fakeTenantRepository)(nil)

func TestCreateTenant_Success(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}
	tenantService := service.NewTenantService(fakeRepo)

	testTenantName := "Test Tenant A"
	testTenantDescription := "Test tenants description"

	tenant, err := tenantService.CreateTenant(
		context.Background(),
		testTenantName,
		testTenantDescription,
	)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if tenant.Name != testTenantName {
		t.Errorf("expected tenant name: %s, got: %s", testTenantName, tenant.Name)
	}

	if tenant.Description != testTenantDescription {
		t.Errorf("expected tenant description: %s, got: %s", testTenantDescription, tenant.Description)
	}

	if tenant.ID == "" {
		t.Error("expected tenant ID to be generated")
	}

	if tenant.CreatedAt.IsZero() {
		t.Error("expected tenant CreatedAt to be set")
	}

	if fakeRepo.createdTenant != tenant {
		t.Errorf("expected repository to receive created tenant %+v, got %+v", tenant, fakeRepo.createdTenant)
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

func TestCreateTenant_RepoError(t *testing.T) {
	createErr := errors.New("create tenant error")
	fakeRepo := &fakeTenantRepository{
		createErr: createErr,
	}

	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.CreateTenant(
		context.Background(),
		"TestTenantC",
		"Test Tenant C",
	)

	if err == nil {
		t.Fatal("expected error while creating tenant")
	}

	if !errors.Is(err, createErr) {
		t.Errorf("expected error: %v, got: %v", createErr, err)
	}
}

func TestGetTenantByID_Success(t *testing.T) {
	testTenant := models.Tenant{
		Name:        "Tenant XYZ",
		Description: "Tenant XYZs description",
		ID:          uuid.NewString(),
		CreatedAt:   time.Now().UTC(),
	}

	fakeRepo := &fakeTenantRepository{
		tenant: testTenant,
	}

	tenantService := service.NewTenantService(fakeRepo)

	tenant, err := tenantService.GetTenantByID(context.Background(), testTenant.ID)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if fakeRepo.gotGetByID != testTenant.ID {
		t.Errorf("expected repository to receive tenant ID %q, got %q", testTenant.ID, fakeRepo.gotGetByID)
	}

	if tenant != testTenant {
		t.Errorf("expected tenant %+v, got %+v", testTenant, tenant)
	}
}

func TestGetTenantByID_EmptyID(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}

	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.GetTenantByID(context.Background(), "")

	if err == nil || !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Fatalf("expected error: %v, got: %v", apperrors.ErrInvalidInput, err)
	}

	if fakeRepo.gotGetByID != "" {
		t.Error("repository should not be called for invalid input")
	}
}

func TestGetTenantByID_WhitespaceID(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}

	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.GetTenantByID(context.Background(), "   ")

	if err == nil || !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Fatalf("expected error: %v, got: %v", apperrors.ErrInvalidInput, err)
	}

	if fakeRepo.gotGetByID != "" {
		t.Error("repository should not be called for invalid input")
	}
}

func TestGetTenantByID_NotFound(t *testing.T) {
	fakeRepo := &fakeTenantRepository{
		getByIDErr: apperrors.ErrTenantNotFound,
	}

	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.GetTenantByID(context.Background(), "abcxyz")

	if err == nil {
		t.Fatal("expected error for missing tenant")
	}

	if !errors.Is(err, apperrors.ErrTenantNotFound) {
		t.Errorf("expected ErrTenantNotFound, got: %v", err)
	}
}

func TestGetTenantByID_RepoError(t *testing.T) {
	getByIDErr := errors.New("get tenant by id error")
	fakeRepo := &fakeTenantRepository{
		getByIDErr: getByIDErr,
	}

	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.GetTenantByID(context.Background(), "abcxyz")

	if err == nil {
		t.Fatal("expected error while fetching tenant by id")
	}

	if !errors.Is(err, getByIDErr) {
		t.Errorf("expected error: %v, got: %v", getByIDErr, err)
	}
}

func TestGetTenantByName_Success(t *testing.T) {
	testTenant := models.Tenant{
		Name:        "Tenant XYZ",
		Description: "Tenant XYZs description",
		ID:          uuid.NewString(),
		CreatedAt:   time.Now().UTC(),
	}

	fakeRepo := &fakeTenantRepository{
		tenant: testTenant,
	}

	tenantService := service.NewTenantService(fakeRepo)

	tenant, err := tenantService.GetTenantByName(context.Background(), "  Tenant XYZ  ")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if fakeRepo.gotGetByName != "Tenant XYZ" {
		t.Errorf("expected repository to receive tenant name %q, got %q", "Tenant XYZ", fakeRepo.gotGetByName)
	}

	if tenant != testTenant {
		t.Errorf("expected tenant %+v, got %+v", testTenant, tenant)
	}
}

func TestGetTenantByName_EmptyName(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}

	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.GetTenantByName(context.Background(), "  ")

	if err == nil || !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Fatalf("expected error: %v, got: %v", apperrors.ErrInvalidInput, err)
	}
}

func TestGetTenantByName_NotFound(t *testing.T) {
	fakeRepo := &fakeTenantRepository{
		getByNameErr: apperrors.ErrTenantNotFound,
	}

	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.GetTenantByName(context.Background(), "abcxyz")

	if err == nil {
		t.Fatal("expected error for missing tenant")
	}

	if !errors.Is(err, apperrors.ErrTenantNotFound) {
		t.Errorf("expected ErrTenantNotFound, got: %v", err)
	}
}

func TestGetTenantByName_RepoError(t *testing.T) {
	getByNameErr := errors.New("get tenant by name error")
	fakeRepo := &fakeTenantRepository{
		getByNameErr: getByNameErr,
	}

	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.GetTenantByName(context.Background(), "abcxyz")

	if err == nil {
		t.Fatal("expected error while fetching tenant by name")
	}

	if !errors.Is(err, getByNameErr) {
		t.Errorf("expected error: %v, got: %v", getByNameErr, err)
	}
}

func TestDeleteTenant_Success(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}

	tenantService := service.NewTenantService(fakeRepo)

	tTenantID := "abcxyz " // white space validation
	err := tenantService.DeleteTenant(context.Background(), tTenantID)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if fakeRepo.deletedID != strings.TrimSpace(tTenantID) {
		t.Errorf("expected deleted tenant id: %s, got: %s", tTenantID, fakeRepo.deletedID)
	}
}

func TestDeleteTenant_EmptyID(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}

	tenantService := service.NewTenantService(fakeRepo)

	err := tenantService.DeleteTenant(context.Background(), "  ")

	if err == nil {
		t.Fatal("expected error while deleting using empty id")
	}

	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("expected error: %v, got: %v", apperrors.ErrInvalidInput, err)
	}
}

func TestDeleteTenant_NotFound(t *testing.T) {
	fakeRepo := &fakeTenantRepository{
		deleteErr: apperrors.ErrTenantNotFound,
	}

	tenantService := service.NewTenantService(fakeRepo)

	err := tenantService.DeleteTenant(context.Background(), "abcxyz")

	if err == nil {
		t.Fatal("expected error while deleting tenant")
	}

	if !errors.Is(err, apperrors.ErrTenantNotFound) {
		t.Errorf("expected error: %v, got %v", apperrors.ErrTenantNotFound, err)
	}
}

func TestDeleteTenant_RepoError(t *testing.T) {
	tDeletionError := errors.New("tenant deletion error")

	fakeRepo := &fakeTenantRepository{
		deleteErr: tDeletionError,
	}

	tenantService := service.NewTenantService(fakeRepo)

	err := tenantService.DeleteTenant(context.Background(), "abcxyz")

	if err == nil {
		t.Fatal("expected error while deleting tenant")
	}

	if !errors.Is(err, tDeletionError) {
		t.Errorf("expected error: %v, got %v", tDeletionError, err)
	}
}

func TestUpdateTenant_Success(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}
	tenantService := service.NewTenantService(fakeRepo)

	tenantToUpdate := models.Tenant{
		ID:          "abcxyz   ", // white space validation
		Name:        "Tenant ABCXYZ   ",
		Description: "ABCXYZ description",
	}

	updatedTenant, err := tenantService.UpdateTenant(context.Background(), tenantToUpdate)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// remove whitespace before validation
	tenantToUpdate.ID = strings.TrimSpace(tenantToUpdate.ID)
	tenantToUpdate.Name = strings.TrimSpace(tenantToUpdate.Name)

	if tenantToUpdate != fakeRepo.updatedTenant {
		t.Errorf("expect tenant repo to receive tenant: %+v, got %+v", tenantToUpdate, fakeRepo.updatedTenant)
	}

	if updatedTenant.ID != tenantToUpdate.ID {
		t.Errorf("expected updated tenant id: %s, got: %s", tenantToUpdate.ID, updatedTenant.ID)
	}

	if updatedTenant.Name != tenantToUpdate.Name {
		t.Errorf("expected updated tenant name: %s, got: %s", tenantToUpdate.Name, updatedTenant.Name)
	}

	if updatedTenant.Description != tenantToUpdate.Description {
		t.Errorf("expected updated tenant description: %s, got: %s", tenantToUpdate.Name, updatedTenant.Description)
	}
}

func TestUpdateTenant_EmptyIDAndName(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}
	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.UpdateTenant(context.Background(),
		models.Tenant{
			ID:          "  ",
			Name:        "  ",
			Description: "Tenant with empty id and name",
		},
	)

	if err == nil {
		t.Fatal("expected error for empty tenant id and name")
	}

	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got: %v", err)
	}

	if fakeRepo.updatedTenant.ID != "" {
		t.Error("tenant repository should not be called for invalid input")
	}
}

func TestUpdateTenant_EmptyID(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}
	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.UpdateTenant(context.Background(),
		models.Tenant{
			ID:          "  ",
			Name:        "Tenant ABCXYZ",
			Description: "Tenant with empty id",
		},
	)

	if err == nil {
		t.Fatal("expected error for empty tenant id")
	}

	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got: %v", err)
	}

	if fakeRepo.updatedTenant.ID != "" {
		t.Error("tenant repository should not be called for invalid input")
	}
}

func TestUpdateTenant_EmptyName(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}
	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.UpdateTenant(context.Background(),
		models.Tenant{
			ID:          "abcxyz",
			Name:        "  ",
			Description: "Tenant with empty name",
		},
	)

	if err == nil {
		t.Fatal("expected error for empty tenant name")
	}

	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got: %v", err)
	}

	if fakeRepo.updatedTenant.ID != "" {
		t.Error("tenant repository should not be called for invalid input")
	}
}

func TestUpdateTenant_NotFound(t *testing.T) {
	fakeRepo := &fakeTenantRepository{
		updateErr: apperrors.ErrTenantNotFound,
	}

	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.UpdateTenant(context.Background(), models.Tenant{
		ID:          "abcxyz",
		Name:        "Tenant ABCXYZ",
		Description: "ABCXYZ description",
	})

	if err == nil {
		t.Fatal("expected error for missing tenant")
	}

	if !errors.Is(err, apperrors.ErrTenantNotFound) {
		t.Errorf("expected ErrTenantNotFound, got: %v", err)
	}
}

func TestUpdateTenant_RepoError(t *testing.T) {
	updateErr := errors.New("tenant update error")

	fakeRepo := &fakeTenantRepository{
		updateErr: updateErr,
	}

	tenantService := service.NewTenantService(fakeRepo)

	_, err := tenantService.UpdateTenant(context.Background(), models.Tenant{
		ID:          "abcxyz",
		Name:        "Tenant ABCXYZ",
		Description: "ABCXYZ description",
	})

	if err == nil {
		t.Fatal("expected error while updating tenant")
	}

	if !errors.Is(err, updateErr) {
		t.Errorf("expected error: %v, got: %v", updateErr, err)
	}
}

func TestGetAllTenants_SuccessWithCursor(t *testing.T) {
	tTenants := []models.Tenant{
		{
			ID:          uuid.NewString(),
			Name:        "Tenant A",
			Description: "First page tenant",
			CreatedAt:   time.Now().UTC(),
		},
		{
			ID:          uuid.NewString(),
			Name:        "Tenant A",
			Description: "First page tenant",
			CreatedAt:   time.Now().UTC(),
		},
	}

	tInputCursor := &models.TenantCursor{
		ID:        "previous-tenant",
		CreatedAt: time.Now().UTC(),
	}

	tNextCursor := &models.TenantCursor{
		ID:        tTenants[1].ID,
		CreatedAt: tTenants[1].CreatedAt,
	}

	fakeRepo := &fakeTenantRepository{
		tenants:    tTenants,
		nextCursor: tNextCursor,
	}

	tenantService := service.NewTenantService(fakeRepo)

	tenants, nextCursor, err :=
		tenantService.GetAllTenants(context.Background(), 2, tInputCursor)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(tenants) != 2 || tenants[0] != tTenants[0] || tenants[1] != tTenants[1] {
		t.Errorf("expected tenants %+v, got %+v", tTenants, tenants)
	}

	if nextCursor != tNextCursor {
		t.Errorf("expected next cursor %+v, got %+v", tNextCursor, nextCursor)
	}

	// called parameters
	if fakeRepo.gotGetAllLimit != 2 {
		t.Errorf("expected repository limit 2, got %d", fakeRepo.gotGetAllLimit)
	}

	if fakeRepo.gotGetAllCursor != tInputCursor {
		t.Errorf("expected repository cursor %+v, got %+v", tInputCursor, fakeRepo.gotGetAllCursor)
	}
}

func TestGetAllTenants_SuccessWithoutCursor(t *testing.T) {
	tTenants := []models.Tenant{
		{
			ID:          uuid.NewString(),
			Name:        "Tenant A",
			Description: "First page tenant",
			CreatedAt:   time.Now().UTC(),
		},
		{
			ID:          uuid.NewString(),
			Name:        "Tenant B",
			Description: "First page tenant",
			CreatedAt:   time.Now().UTC(),
		},
	}

	tNextCursor := &models.TenantCursor{
		ID:        tTenants[1].ID,
		CreatedAt: tTenants[1].CreatedAt,
	}

	fakeRepo := &fakeTenantRepository{
		tenants:    tTenants,
		nextCursor: tNextCursor,
	}
	tenantService := service.NewTenantService(fakeRepo)

	tenants, nextCursor, err :=
		tenantService.GetAllTenants(context.Background(), 2, nil)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(tenants) != 2 || tenants[0] != tTenants[0] || tenants[1] != tTenants[1] {
		t.Errorf("expected tenants %+v, got %+v", tTenants, tenants)
	}

	if nextCursor != tNextCursor {
		t.Errorf("expected next cursor %+v, got %+v", tNextCursor, nextCursor)
	}

	// called parameters
	if fakeRepo.gotGetAllLimit != 2 {
		t.Errorf("expected repository limit 2, got %d", fakeRepo.gotGetAllLimit)
	}

	if fakeRepo.gotGetAllCursor != nil {
		t.Errorf("expected repository cursor nil, got %+v", fakeRepo.gotGetAllCursor)
	}
}

func TestGetAllTenants_ZeroLimit(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}
	tenantService := service.NewTenantService(fakeRepo)

	tenants, nextCursor, err := tenantService.GetAllTenants(context.Background(), 0, nil)

	if err == nil {
		t.Fatal("expected error for invalid limit")
	}

	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got: %v", err)
	}

	if tenants != nil {
		t.Errorf("expected nil tenants, got %+v", tenants)
	}

	if nextCursor != nil {
		t.Errorf("expected nil next cursor, got %+v", nextCursor)
	}

	if fakeRepo.gotGetAllLimit != 0 || fakeRepo.gotGetAllCursor != nil {
		t.Error("tenant repository should not be called for invalid limit")
	}
}

func TestGetAllTenants_NegativeLimit(t *testing.T) {
	fakeRepo := &fakeTenantRepository{}
	tenantService := service.NewTenantService(fakeRepo)

	tenants, nextCursor, err := tenantService.GetAllTenants(context.Background(), -1, nil)

	if err == nil {
		t.Fatal("expected error for negative limit")
	}

	if !errors.Is(err, apperrors.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got: %v", err)
	}

	if tenants != nil {
		t.Errorf("expected nil tenants, got %+v", tenants)
	}

	if nextCursor != nil {
		t.Errorf("expected nil next cursor, got %+v", nextCursor)
	}

	if fakeRepo.gotGetAllLimit != 0 || fakeRepo.gotGetAllCursor != nil {
		t.Error("tenant repository should not be called for invalid limit")
	}
}

func TestGetAllTenants_Error(t *testing.T) {
	getAllErr := errors.New("get all tenants error")

	fakeRepo := &fakeTenantRepository{
		getAllErr: getAllErr,
	}
	tenantService := service.NewTenantService(fakeRepo)

	_, _, err := tenantService.GetAllTenants(context.Background(), 2, nil)

	if err == nil {
		t.Fatal("expected error while retrieving tenants")
	}

	if !errors.Is(err, getAllErr) {
		t.Errorf("expected error: %v, got: %v", getAllErr, err)
	}

	if fakeRepo.gotGetAllLimit != 2 {
		t.Errorf("expected repository limit 2, got %d", fakeRepo.gotGetAllLimit)
	}

	if fakeRepo.gotGetAllCursor != nil {
		t.Errorf("expected repository cursor nil, got %+v", fakeRepo.gotGetAllCursor)
	}
}
