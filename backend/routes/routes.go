package routes

import (
	"database/sql"
	"net/http"

	"tshongmart/handlers"
)

func RegisterRoutes(db *sql.DB) {
	authHandler := handlers.NewAuthHandler(db)
	productHandler := handlers.NewProductHandler(db)
	userHandler := handlers.NewUserHandler(db)
	contactHandler := handlers.NewContactHandler(db)
	adminHandler := handlers.NewAdminHandler(db)

	http.HandleFunc("/api/health", handlers.HealthHandler)

	http.HandleFunc("/api/signup", authHandler.SignupHandler)
	http.HandleFunc("/api/login", authHandler.LoginHandler)
	http.HandleFunc("/api/admin-login", handlers.AdminLoginHandler)
	http.HandleFunc("/api/users", adminHandler.GetUsersHandler)
	http.HandleFunc("/api/users/status", adminHandler.UpdateUserStatusHandler)

	http.HandleFunc("/api/profile", userHandler.UpdateProfileHandler)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	http.HandleFunc("/api/contact-messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			contactHandler.CreateContactMessageHandler(w, r)
			return
		}

		if r.Method == http.MethodGet {
			contactHandler.GetContactMessagesHandler(w, r)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			productHandler.GetProductsHandler(w, r)
			return
		}

		if r.Method == http.MethodPost {
			productHandler.AddProductHandler(w, r)
			return
		}

		if r.Method == http.MethodDelete {
			productHandler.DeleteProductHandler(w, r)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/api/my-products", productHandler.GetMyProductsHandler)

	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)
}