package main

import (
	"fmt"
	"log"
	"net/http"

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

	fmt.Println("Server started on port 8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}