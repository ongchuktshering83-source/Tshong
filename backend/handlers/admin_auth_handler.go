package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"tshongmart/models"
	"tshongmart/utils"
)

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var request models.AdminLoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminPassword == "" {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Admin password is not configured",
		})
		return
	}

	if request.Password != adminPassword {
		utils.SendJSON(w, http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid admin password",
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Admin access granted",
	})
}