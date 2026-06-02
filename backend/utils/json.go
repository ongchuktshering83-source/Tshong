package utils

import (
	"encoding/json"
	"net/http"

	"tshongmart/models"
)

func SendJSON(w http.ResponseWriter, statusCode int, response models.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func SendLoginJSON(w http.ResponseWriter, statusCode int, response models.LoginResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func SendUpdateProfileJSON(w http.ResponseWriter, statusCode int, response models.UpdateProfileResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}