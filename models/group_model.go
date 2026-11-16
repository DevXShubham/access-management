package models

type Group struct {
	ID          string
	Name        string
	Description string
	// No tenantID because a group can be associated with multiple tenants
}
