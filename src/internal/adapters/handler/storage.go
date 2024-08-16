package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"pi-monitor-api/internal/core/ports"
)

type StorageHandler struct {
	StorageService ports.StoragePort
}

func NewStorageHandler(service ports.StoragePort) *StorageHandler {
	return &StorageHandler{StorageService: service}
}

func (h *StorageHandler) GetStorageInfo(w http.ResponseWriter, r *http.Request) {
	storageInfo, err := h.StorageService.GetDevices()
	if err != nil {
		log.Printf("Error retrieving storage info: %v", err)
		http.Error(w, "Failed to retrieve storage info", http.StatusInternalServerError)
		return
	}

	log.Printf("Storage info retrieved successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storageInfo)
}
