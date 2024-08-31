package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alvmarrod/pi-monitor-api/internal/core/ports"
)

type RAMHandler struct {
	RAMService ports.RAMPort
}

func NewRAMHandler(service ports.RAMPort) *RAMHandler {
	return &RAMHandler{RAMService: service}
}

func (h *RAMHandler) GetRAMInfo(w http.ResponseWriter, r *http.Request) {
	ramInfo, err := h.RAMService.GetRAMStats()
	if err != nil {
		log.Printf("Error retrieving RAM info: %v", err)
		http.Error(w, "Failed to retrieve RAM info", http.StatusInternalServerError)
		return
	}

	log.Printf("RAM info retrieved successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ramInfo)
}
