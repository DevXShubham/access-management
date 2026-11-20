package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	// "iam/repository"
	"iam/apperrors"
	"iam/models"
)

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(databaseConnection *sql.DB) *TenantRepository {
	return &TenantRepository{
		db: databaseConnection,
	}
}

func (r *TenantRepository) Create(ctx context.Context, tenant models.Tenant) error {
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS tenants (
			id STRING NOT NULL PRIMARY KEY,
			name STRING NOT NULL UNIQUE,
			description STRING NOT NULL,
			createdAt STRING NOT NULL
		)
	`
	_, err := r.db.ExecContext(ctx, sqlStmt)

	if err != nil {
		return fmt.Errorf("creating table 'tenants': %w", err)
	}

	_, err = r.db.ExecContext(ctx,
		"INSERT INTO tenants(id, name, description, createdAt) values(:id, :name, :description, :createdAt)",
		sql.Named("id", tenant.ID),
		sql.Named("name", tenant.Name),
		sql.Named("description", tenant.Description),
		sql.Named("createdAt", tenant.CreatedAt),
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("%w: tenant '%s' already exists", apperrors.ErrTenantAlreadyExists, tenant.Name)
		}

		return fmt.Errorf("inserting tenant: %w", err)
	}

	return nil
}

func parseTime(strTime string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339Nano, strTime)

	if err != nil {
		parsedTime, err = time.Parse("2006-01-02 15:04:05.999999999Z07:00", strTime)
	}

	return parsedTime, err
}

func (r *TenantRepository) GetByID(ctx context.Context, id string) (models.Tenant, error) {

	query := "SELECT id, name, description, createdAt FROM tenants WHERE id = :id"

	rows, err := r.db.QueryContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return models.Tenant{}, fmt.Errorf("getting tenant by id: %w", err)
	}

	if !rows.Next() {
		return models.Tenant{}, fmt.Errorf("%w (id: %s)", apperrors.ErrTenantNotFound, id)
	}

	var createdAtStr string
	var tenant models.Tenant
	err = rows.Scan(&tenant.ID, &tenant.Name, &tenant.Description, &createdAtStr)

	if err != nil {
		return models.Tenant{}, fmt.Errorf("getting tenant by id: %w", err)
	}

	createdAt, err := parseTime(createdAtStr)

	if err != nil {
		return models.Tenant{}, fmt.Errorf("parsing time: %w", err)
	}

	tenant.CreatedAt = createdAt

	return tenant, nil
}

func (r *TenantRepository) GetByName(ctx context.Context, name string) (models.Tenant, error) {

	query := "SELECT id, name, description, createdAt FROM tenants WHERE name = :name"

	rows, err := r.db.QueryContext(ctx, query, sql.Named("name", name))

	if err != nil {
		return models.Tenant{}, fmt.Errorf("getting tenant by id: %w", err)
	}

	if !rows.Next() {
		return models.Tenant{}, fmt.Errorf("%w (name: %s)", apperrors.ErrTenantNotFound, name)
	}

	var createdAtStr string
	var tenant models.Tenant
	err = rows.Scan(&tenant.ID, &tenant.Name, &tenant.Description, &createdAtStr)

	if err != nil {
		return models.Tenant{}, fmt.Errorf("getting tenant by name: %w", err)
	}

	createdAt, err := parseTime(createdAtStr)

	if err != nil {
		return models.Tenant{}, fmt.Errorf("parsing time: %w", err)
	}

	tenant.CreatedAt = createdAt

	return tenant, nil
}

func (r *TenantRepository) Update(ctx context.Context, tenantToUpdate models.Tenant) (models.Tenant, error) {
	query := `
		UPDATE tenants
		SET name=:name, description=:description
		WHERE id=:id
	`

	result, err := r.db.ExecContext(ctx, query,
		sql.Named("name", tenantToUpdate.Name),
		sql.Named("description", tenantToUpdate.Description),
		sql.Named("id", tenantToUpdate.ID),
	)

	if err != nil {
		// if strings.Contains(err.Error(), )
		return models.Tenant{}, fmt.Errorf("updating tenant: %w", err)
	}

	affectedRows, _ := result.RowsAffected()

	if affectedRows == 0 {
		return models.Tenant{}, fmt.Errorf("%w: (id: %s)", apperrors.ErrTenantNotFound, tenantToUpdate.ID)
	}

	return r.GetByID(ctx, tenantToUpdate.ID)
}

func (r *TenantRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM tenants WHERE id=:id"

	result, err := r.db.ExecContext(ctx, query, sql.Named("id", id))

	if err != nil {
		return fmt.Errorf("deleting tenant: %w", err)
	}

	affectedRows, _ := result.RowsAffected()

	if affectedRows == 0 {
		return fmt.Errorf("%w: (id: %s)", apperrors.ErrTenantNotFound, id)
	}

	return nil
}

func (r *TenantRepository) GetAll(ctx context.Context, limit int, cursor *models.TenantCursor) ([]models.Tenant, *models.TenantCursor, error) {
	return nil, nil, nil
}
