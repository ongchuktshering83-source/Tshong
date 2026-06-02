package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"tshongmart/models"
	"tshongmart/utils"
)

type AdminHandler struct {
	DB *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	rows, err := h.DB.Query(
		`SELECT id, full_name, email, contact_info, role, COALESCE(status, 'active'), created_at::text
		 FROM users
		 ORDER BY created_at DESC`,
	)

	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not fetch users",
		})
		return
	}

	defer rows.Close()

	users := []models.AdminUserResponse{}

	for rows.Next() {
		var user models.AdminUserResponse

		err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.ContactInfo,
			&user.Role,
			&user.Status,
			&user.CreatedAt,
		)

		if err != nil {
			utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Could not read user data",
			})
			return
		}

		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *AdminHandler) UpdateUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	userIDText := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")

	if userIDText == "" || status == "" {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "User ID and status are required",
		})
		return
	}

	if status != "active" && status != "banned" {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid status",
		})
		return
	}

	userID, err := strconv.Atoi(userIDText)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	_, err = h.DB.Exec(
		`UPDATE users SET status = $1 WHERE id = $2`,
		status,
		userID,
	)

	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not update user status",
		})
		return
	}

	message := "User activated successfully"
	if status == "banned" {
		message = "User banned successfully"
	}

	utils.SendJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: message,
	})
}

func (h *AdminHandler) ClearTestDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var request models.ClearTestDataRequest

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

	tx, err := h.DB.Begin()
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not start cleanup",
		})
		return
	}

	_, err = tx.Exec(`DELETE FROM contact_messages`)
	if err != nil {
		tx.Rollback()
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not clear contact messages",
		})
		return
	}

	_, err = tx.Exec(`DELETE FROM products`)
	if err != nil {
		tx.Rollback()
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not clear products",
		})
		return
	}

	_, err = tx.Exec(`DELETE FROM users`)
	if err != nil {
		tx.Rollback()
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not clear users",
		})
		return
	}

	err = tx.Commit()
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not finish cleanup",
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Test data cleared successfully",
	})
}