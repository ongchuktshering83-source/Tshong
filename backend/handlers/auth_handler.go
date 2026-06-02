package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"tshongmart/models"
	"tshongmart/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var request models.SignupRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.FullName == "" || request.Email == "" || request.ContactInfo == "" || request.Password == "" {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "All required fields must be filled",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not secure password",
		})
		return
	}

	_, err = h.DB.Exec(
		`INSERT INTO users (full_name, email, contact_info, password_hash, status)
		 VALUES ($1, $2, $3, $4, 'active')`,
		request.FullName,
		request.Email,
		request.ContactInfo,
		string(hashedPassword),
	)

	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Email may already be registered",
		})
		return
	}

	utils.SendJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Account created successfully",
	})
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendLoginJSON(w, http.StatusMethodNotAllowed, models.LoginResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var request models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.SendLoginJSON(w, http.StatusBadRequest, models.LoginResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.Email == "" || request.Password == "" {
		utils.SendLoginJSON(w, http.StatusBadRequest, models.LoginResponse{
			Success: false,
			Message: "Email and password are required",
		})
		return
	}

	var user models.UserResponse
	var passwordHash string

	err = h.DB.QueryRow(
		`SELECT 
			id,
			full_name,
			email,
			contact_info,
			role,
			COALESCE(status, 'active'),
			password_hash
		 FROM users
		 WHERE email = $1`,
		request.Email,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.ContactInfo,
		&user.Role,
		&user.Status,
		&passwordHash,
	)

	if err != nil {
		utils.SendLoginJSON(w, http.StatusUnauthorized, models.LoginResponse{
			Success: false,
			Message: "Invalid email or password",
		})
		return
	}

	if user.Status == "banned" {
		utils.SendLoginJSON(w, http.StatusForbidden, models.LoginResponse{
			Success: false,
			Message: "Your account has been banned. Please contact the admin.",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(request.Password))
	if err != nil {
		utils.SendLoginJSON(w, http.StatusUnauthorized, models.LoginResponse{
			Success: false,
			Message: "Invalid email or password",
		})
		return
	}

	utils.SendLoginJSON(w, http.StatusOK, models.LoginResponse{
		Success: true,
		Message: "Login successful",
		User:    user,
	})
}