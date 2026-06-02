package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tshongmart/config"
	"tshongmart/routes"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Println("Connected to PostgreSQL database successfully.")
	defer db.Close()

	routes.RegisterRoutes(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server started on port " + port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}