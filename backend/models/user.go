package models

type SignupRequest struct {
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	ContactInfo string `json:"contactInfo"`
	Password    string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID          int    `json:"id"`
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	ContactInfo string `json:"contactInfo"`
	Role        string `json:"role"`
	Status      string `json:"status"`
}

type UpdateProfileRequest struct {
	ID              int    `json:"id"`
	FullName        string `json:"fullName"`
	Email           string `json:"email"`
	ContactInfo     string `json:"contactInfo"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type UpdateProfileResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}

type AdminUserResponse struct {
	ID          int    `json:"id"`
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	ContactInfo string `json:"contactInfo"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
}