package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"tshongmart/models"
	"tshongmart/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.SendUpdateProfileJSON(w, http.StatusMethodNotAllowed, models.UpdateProfileResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var request models.UpdateProfileRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.SendUpdateProfileJSON(w, http.StatusBadRequest, models.UpdateProfileResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.ID == 0 ||
		request.FullName == "" ||
		request.Email == "" ||
		request.ContactInfo == "" ||
		request.CurrentPassword == "" {
		utils.SendUpdateProfileJSON(w, http.StatusBadRequest, models.UpdateProfileResponse{
			Success: false,
			Message: "All required fields must be filled",
		})
		return
	}

	var passwordHash string

	err = h.DB.QueryRow(
		`SELECT password_hash FROM users WHERE id = $1`,
		request.ID,
	).Scan(&passwordHash)

	if err != nil {
		utils.SendUpdateProfileJSON(w, http.StatusNotFound, models.UpdateProfileResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(request.CurrentPassword))
	if err != nil {
		utils.SendUpdateProfileJSON(w, http.StatusUnauthorized, models.UpdateProfileResponse{
			Success: false,
			Message: "Current password is incorrect",
		})
		return
	}

	finalPasswordHash := passwordHash

	if request.NewPassword != "" {
		newHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			utils.SendUpdateProfileJSON(w, http.StatusInternalServerError, models.UpdateProfileResponse{
				Success: false,
				Message: "Could not secure new password",
			})
			return
		}

		finalPasswordHash = string(newHash)
	}

	_, err = h.DB.Exec(
		`UPDATE users
		 SET full_name = $1,
		     email = $2,
		     contact_info = $3,
		     password_hash = $4
		 WHERE id = $5`,
		request.FullName,
		request.Email,
		request.ContactInfo,
		finalPasswordHash,
		request.ID,
	)

	if err != nil {
		utils.SendUpdateProfileJSON(w, http.StatusBadRequest, models.UpdateProfileResponse{
			Success: false,
			Message: "Could not update profile. Email may already be used.",
		})
		return
	}

	var updatedUser models.UserResponse

	err = h.DB.QueryRow(
		`SELECT id, full_name, email, contact_info, role, COALESCE(status, 'active')
         FROM users
         WHERE id = $1`,
		request.ID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.FullName,
		&updatedUser.Email,
		&updatedUser.ContactInfo,
		&updatedUser.Role,
		&updatedUser.Status,
	)

	if err != nil {
		utils.SendUpdateProfileJSON(w, http.StatusInternalServerError, models.UpdateProfileResponse{
			Success: false,
			Message: "Profile updated, but could not fetch updated user",
		})
		return
	}

	utils.SendUpdateProfileJSON(w, http.StatusOK, models.UpdateProfileResponse{
		Success: true,
		Message: "Profile updated successfully",
		User:    updatedUser,
	})
}