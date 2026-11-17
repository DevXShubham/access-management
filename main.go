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

	log.Println("starting http server at http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
