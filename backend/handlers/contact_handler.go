package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"tshongmart/models"
	"tshongmart/utils"
)

type ContactHandler struct {
	DB *sql.DB
}

func NewContactHandler(db *sql.DB) *ContactHandler {
	return &ContactHandler{DB: db}
}

func (h *ContactHandler) CreateContactMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var request models.ContactMessageRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.FullName == "" || request.Email == "" || request.Subject == "" || request.Message == "" {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "All contact fields are required",
		})
		return
	}

	_, err = h.DB.Exec(
		`INSERT INTO contact_messages (full_name, email, subject, message)
		 VALUES ($1, $2, $3, $4)`,
		request.FullName,
		request.Email,
		request.Subject,
		request.Message,
	)

	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not send message",
		})
		return
	}

	utils.SendJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Message sent successfully",
	})
}

func (h *ContactHandler) GetContactMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	rows, err := h.DB.Query(
		`SELECT id, full_name, email, subject, message, created_at::text
		 FROM contact_messages
		 ORDER BY created_at DESC`,
	)

	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not fetch contact messages",
		})
		return
	}

	defer rows.Close()

	messages := []models.ContactMessageResponse{}

	for rows.Next() {
		var message models.ContactMessageResponse

		err := rows.Scan(
			&message.ID,
			&message.FullName,
			&message.Email,
			&message.Subject,
			&message.Message,
			&message.CreatedAt,
		)

		if err != nil {
			utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Could not read contact message data",
			})
			return
		}

		messages = append(messages, message)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}