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
	_, err := r.db.ExecContext(ctx,
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

		return fmt.Errorf("inserting tenant: %w\n%w", apperrors.ErrDBExecution, err)
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
		return models.Tenant{}, fmt.Errorf("getting tenant by id: %w\n%w", apperrors.ErrDBExecution, err)
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
		return models.Tenant{}, fmt.Errorf("parsing time: %w\n%w", apperrors.ErrParseTime, err)
	}

	tenant.CreatedAt = createdAt

	return tenant, nil
}

func (r *TenantRepository) GetByName(ctx context.Context, name string) (models.Tenant, error) {

	query := "SELECT id, name, description, createdAt FROM tenants WHERE name = :name"

	rows, err := r.db.QueryContext(ctx, query, sql.Named("name", name))

	if err != nil {
		return models.Tenant{}, fmt.Errorf("getting tenant by id: %w\n%w", apperrors.ErrDBExecution, err)
	}

	if !rows.Next() {
		return models.Tenant{}, fmt.Errorf("%w (name: %s)", apperrors.ErrTenantNotFound, name)
	}

	var createdAtStr string
	var tenant models.Tenant
	err = rows.Scan(&tenant.ID, &tenant.Name, &tenant.Description, &createdAtStr)

	if err != nil {
		return models.Tenant{}, fmt.Errorf("getting tenant by name: %w\n%w", apperrors.ErrDBExecution, err)
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
		return models.Tenant{}, fmt.Errorf("updating tenant: %w\n%w", apperrors.ErrDBExecution, err)
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
		return fmt.Errorf("deleting tenant: %w\n%w", apperrors.ErrDBExecution, err)
	}

	affectedRows, _ := result.RowsAffected()

	if affectedRows == 0 {
		return fmt.Errorf("%w: (id: %s)", apperrors.ErrTenantNotFound, id)
	}

	return nil
}

func (r *TenantRepository) GetAll(ctx context.Context, limit int, cursor *models.TenantCursor) (
	[]models.Tenant, *models.TenantCursor, error,
) {
	query := "SELECT id, name, description, createdAt FROM tenants "

	if cursor != nil {
		query += `
			WHERE createdAt < :createdAtCursor
			OR (createdAt = :createdAtCursor
			AND id < :idCursor)
		`
	}

	query += `
		ORDER BY createdAt DESC, id DESC
		LIMIT :limit
	`

	var rows *sql.Rows
	var err error

	if cursor != nil {
		rows, err = r.db.QueryContext(ctx, query,
			sql.Named("createdAtCursor", cursor.CreatedAt),
			sql.Named("idCursor", cursor.ID),
			sql.Named("limit", limit+1),
		)
	} else {
		rows, err = r.db.QueryContext(ctx, query, sql.Named("limit", limit+1))
	}

	if err != nil {
		return nil, nil, fmt.Errorf("fetching tenants: %w\n%w", apperrors.ErrDBExecution, err)
	}

	var tenants []models.Tenant

	for rows.Next() {
		var id string
		var name string
		var description string
		var createdAtStr string

		rows.Scan(&id, &name, &description, &createdAtStr)

		createdAt, err := parseTime(createdAtStr)

		if err != nil {
			return nil, nil, fmt.Errorf("fetching tenants: %w", apperrors.ErrParseTime)
		}

		tenant := models.Tenant{
			ID:          id,
			Name:        name,
			Description: description,
			CreatedAt:   createdAt,
		}

		tenants = append(tenants, tenant)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("reading rows: %w\n%w", apperrors.ErrDBExecution, err)
	}

	hasMore := len(tenants) > limit

	var nextCursor *models.TenantCursor

	if hasMore {
		lastTenant := tenants[len(tenants)-2]

		nextCursor = &models.TenantCursor{
			CreatedAt: lastTenant.CreatedAt,
			ID:        lastTenant.ID,
		}

		tenants = tenants[:len(tenants)-1]
	}

	return tenants, nextCursor, nil
}
