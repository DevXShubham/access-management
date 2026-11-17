package main

import (
	"log"
	"net/http"

	"iam/handler"
	"iam/repository/memory"
	"iam/service"
)

func main() {
	tenantRepo := memory.NewTenantRepository()
	tenantService := service.NewTenantService(tenantRepo)
	tenantHandler := handler.NewTenantHandler(tenantService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tenants", tenantHandler.CreateTenant)
	mux.HandleFunc("GET /tenants/{id}", tenantHandler.GetTenantByID)
	mux.HandleFunc("GET /tenants/{id}", tenantHandler.GetTenantByName)
	mux.HandleFunc("PUT /tenants/{id}", tenantHandler.UpdateTenant)
	mux.HandleFunc("DELETE /tenants/{id}", tenantHandler.DeleteTenant)

	log.Println("starting http server at http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
