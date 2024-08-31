package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alvmarrod/pi-monitor-api/internal/core/ports"
)

type NetworkHandler struct {
	NetworkService ports.NetworkPort
}

func NewNetworkHandler(service ports.NetworkPort) *NetworkHandler {
	return &NetworkHandler{NetworkService: service}
}

func (h *NetworkHandler) GetNetworkInfo(w http.ResponseWriter, r *http.Request) {
	networkInfo, err := h.NetworkService.GetNetworkInterfaces()
	if err != nil {
		log.Printf("Error retrieving network info: %v", err)
		http.Error(w, "Failed to retrieve network info", http.StatusInternalServerError)
		return
	}

	log.Printf("Network info retrieved successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(networkInfo)
}
