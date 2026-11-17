package handler

import (
	"encoding/json"
	"net/http"

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
