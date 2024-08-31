package main

import (
	"log"
	"net/http"

	"github.com/alvmarrod/pi-monitor-api/internal/adapters/handler"
	"github.com/alvmarrod/pi-monitor-api/internal/adapters/repository"
	"github.com/alvmarrod/pi-monitor-api/internal/core/services"

	"github.com/gorilla/mux"
)

// RegisterV1Routes sets up the routes for version 1 of the API
func RegisterV1Routes(r *mux.Router) {

	//Instantiate the real components that can be mocked during testing
	fileReader := &repository.RealFileReader{}
	execFinder := &repository.RealToolInstalled{}
	cmd := &repository.RealCmdExecutor{}

	// Initialize repositories and services
	cpuRepo := repository.NewCPURepository(fileReader)
	ramRepo := repository.NewRAMRepository(fileReader)
	storageRepo := repository.NewStorageRepository(fileReader, execFinder, cmd)
	networkRepo := repository.NewNetworkRepository(fileReader, execFinder, cmd)

	cpuService := services.NewCPUService(cpuRepo)
	ramService := services.NewRAMService(ramRepo)
	storageService := services.NewStorageService(storageRepo)
	networkService := services.NewNetworkService(networkRepo)

	// Initialize handlers
	cpuHandler := handler.NewCPUHandler(cpuService)
	ramHandler := handler.NewRAMHandler(ramService)
	storageHandler := handler.NewStorageHandler(storageService)
	networkHandler := handler.NewNetworkHandler(networkService)

	// Create a subrouter for version 1 of the API
	v1 := r.PathPrefix("/v1").Subrouter()

	// Define endpoints under the /v1 prefix
	v1.HandleFunc("/cpu", cpuHandler.GetCPULoad).Methods("GET")
	v1.HandleFunc("/ram", ramHandler.GetRAMInfo).Methods("GET")
	v1.HandleFunc("/storage", storageHandler.GetStorageInfo).Methods("GET")
	v1.HandleFunc("/network", networkHandler.GetNetworkInfo).Methods("GET")
}

func main() {
	// Set up the main router
	r := mux.NewRouter()

	// Register all v1 routes
	RegisterV1Routes(r)

	// Start the HTTP server
	log.Println("Starting API server on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
