package main

import (
	"log"
	"net/http"

	"iam/handler"
	"iam/repository/sqlite"
	"iam/service"
)

func logMiddleWare(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving request %s %s", r.Method, r.URL.Path)
		handlerFunc(w, r)
		log.Printf("served request %s %s", r.Method, r.URL.Path)
	}
}

func main() {
	databaseConnection := sqlite.Connect()

	tenantRepo := sqlite.NewTenantRepository(databaseConnection)
	tenantService := service.NewTenantService(tenantRepo)
	tenantHandler := handler.NewTenantHandler(tenantService)

	mux := http.NewServeMux()

	handlerFunctions := map[string]http.HandlerFunc{
		"POST /tenants":            tenantHandler.CreateTenant,
		"GET /tenants":             tenantHandler.GetAllTenants,
		"GET /tenants/{id}":        tenantHandler.GetTenantByID,
		"GET /tenants/name/{name}": tenantHandler.GetTenantByName,
		"PUT /tenants":             tenantHandler.UpdateTenant,
		"DELETE /tenants/{id}":     tenantHandler.DeleteTenant,
	}

	for route := range handlerFunctions {
		mux.HandleFunc(route, logMiddleWare(handlerFunctions[route]))
	}

	log.Println("starting http server at http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
