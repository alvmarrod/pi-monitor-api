package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"pi-monitor-api/internal/core/ports"
)

type CPUHandler struct {
	CPUService ports.CPUPort
}

func NewCPUHandler(service ports.CPUPort) *CPUHandler {
	return &CPUHandler{CPUService: service}
}

func (h *CPUHandler) GetCPULoad(w http.ResponseWriter, r *http.Request) {
	cpuLoad, err := h.CPUService.GetCPULoad()
	if err != nil {
		log.Printf("Error retrieving CPU load: %v", err)
		http.Error(w, "Failed to retrieve CPU load", http.StatusInternalServerError)
		return
	}

	log.Printf("CPU load retrieved successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cpuLoad)
}
