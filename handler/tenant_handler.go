package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"iam/apperrors"
	"iam/models"
	"iam/service"
)

type TenantHandler struct {
	service *service.TenantService
}

func NewTenantHandler(service *service.TenantService) *TenantHandler {
	return &TenantHandler{
		service: service,
	}
}

type CreateTenantRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateTenantResponse struct {
	Message string        `json:"message"`
	Tenant  models.Tenant `json:"tenant"`
}

func (h *TenantHandler) CreateTenant(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateTenantRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tenant, err := h.service.CreateTenant(r.Context(), req.Name, req.Description)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateTenantResponse{
		Message: "tenant created",
		Tenant:  tenant,
	})
}

type GetTenantResponse struct {
	Message string        `json:"message"`
	Tenant  models.Tenant `json:"tenant"`
}

func (h *TenantHandler) GetTenantByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")

	tenant, err := h.service.GetTenantByID(r.Context(), id)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetTenantResponse{
		Message: "tenant found",
		Tenant:  tenant,
	})
}

func (h *TenantHandler) GetTenantByName(
	w http.ResponseWriter,
	r *http.Request,
) {
	name := r.PathValue("name")

	tenant, err := h.service.GetTenantByName(r.Context(), name)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetTenantResponse{
		Message: "tenant found",
		Tenant:  tenant,
	})
}

type UpdateTenantRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateTenantResponse struct {
	Message string        `json:"message"`
	Tenant  models.Tenant `json:"tenant"`
}

func (h *TenantHandler) UpdateTenant(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req UpdateTenantRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tenantToUpdate := models.Tenant{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}

	updatedTenant, err := h.service.UpdateTenant(r.Context(), tenantToUpdate)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(UpdateTenantResponse{
		Message: "tenant updated",
		Tenant:  updatedTenant,
	})
}

type DeleteTenantResponse struct {
	Message  string `json:"message"`
	TenantID string `json:"id"`
}

func (h *TenantHandler) DeleteTenant(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")

	err := h.service.DeleteTenant(r.Context(), id)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DeleteTenantResponse{
		Message:  "tenant deleted",
		TenantID: id,
	})
}

type GetAllTenantsResponse struct {
	Message    string          `json:"message"`
	Tenants    []models.Tenant `json:"tenants"`
	NextCursor *string         `json:"nextCursor"`
}

func (h *TenantHandler) GetAllTenants(
	w http.ResponseWriter,
	r *http.Request,
) {
	limitStr := r.URL.Query().Get("limit")

	limit := 10

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)

		if err != nil {
			WriteError(w, apperrors.ErrInvalidInput)
			return
		}

		limit = parsedLimit
	}

	var cursor *models.TenantCursor

	cursorStr := r.URL.Query().Get("cursor")

	if cursorStr != "" {
		decoded, err := base64.StdEncoding.DecodeString(cursorStr)

		if err != nil {
			WriteError(w, apperrors.ErrInvalidInput)
			return
		}

		var parsedCursor models.TenantCursor

		if err := json.Unmarshal(decoded, &parsedCursor); err != nil {
			WriteError(w, apperrors.ErrInvalidInput)
			return
		}

		cursor = &parsedCursor
	}

	tenants, nextCursor, err := h.service.GetAllTenants(r.Context(), limit, cursor)

	if err != nil {
		WriteError(w, err)
		return
	}

	var nextCursorEncoded *string

	if nextCursor != nil {
		nextCursorJson, err := json.Marshal(nextCursor)
		if err != nil {
			WriteError(w, err)
			return
		}

		enc := base64.StdEncoding.EncodeToString(nextCursorJson)
		nextCursorEncoded = &enc
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetAllTenantsResponse{
		Message:    "tenants retrieved",
		Tenants:    tenants,
		NextCursor: nextCursorEncoded,
	})
}
