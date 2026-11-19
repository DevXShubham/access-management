package models

import "time"

type Tenant struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
}

type TenantCursor struct {
	CreatedAt time.Time
	ID        string
}
