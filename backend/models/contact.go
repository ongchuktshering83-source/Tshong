package models

type ContactMessageRequest struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

type ContactMessageResponse struct {
	ID        int    `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
	CreatedAt string `json:"createdAt"`
}