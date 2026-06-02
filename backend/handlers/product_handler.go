package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"tshongmart/models"
	"tshongmart/utils"
)

type ProductHandler struct {
	DB *sql.DB
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

func (h *ProductHandler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid form data",
		})
		return
	}

	userID, err := strconv.Atoi(r.FormValue("userId"))
	if err != nil || userID == 0 {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Valid user ID is required",
		})
		return
	}

	title := r.FormValue("title")
	category := r.FormValue("category")
	price := r.FormValue("price")
	location := r.FormValue("location")
	contact := r.FormValue("contact")
	description := r.FormValue("description")

	if title == "" || category == "" || price == "" || location == "" || contact == "" || description == "" {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "All required product fields must be filled",
		})
		return
	}

	imagePath := ""

	file, fileHeader, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		uploadDir := "uploads"
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Could not create upload folder",
			})
			return
		}

		extension := filepath.Ext(fileHeader.Filename)
		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		savePath := filepath.Join(uploadDir, fileName)

		destination, err := os.Create(savePath)
		if err != nil {
			utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Could not save image",
			})
			return
		}
		defer destination.Close()

		_, err = io.Copy(destination, file)
		if err != nil {
			utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Could not upload image",
			})
			return
		}

		imagePath = "/uploads/" + fileName
	}

	_, err = h.DB.Exec(
		`INSERT INTO products 
		(user_id, title, category, price, location, contact, description, image_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		userID,
		title,
		category,
		price,
		location,
		contact,
		description,
		imagePath,
	)

	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not add product",
		})
		return
	}

	utils.SendJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Product added successfully",
	})
}

func (h *ProductHandler) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	rows, err := h.DB.Query(
		`SELECT 
			products.id,
			products.user_id,
			products.title,
			products.category,
			products.price,
			products.location,
			products.contact,
			products.description,
			COALESCE(products.image_path, ''),
			users.full_name,
			products.status,
			products.created_at::text
		FROM products
		JOIN users ON products.user_id = users.id
		WHERE COALESCE(users.status, 'active') = 'active'
		ORDER BY products.created_at DESC`,
	)

	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not fetch products",
		})
		return
	}

	defer rows.Close()

	products := []models.ProductResponse{}

	for rows.Next() {
		var product models.ProductResponse

		err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.Title,
			&product.Category,
			&product.Price,
			&product.Location,
			&product.Contact,
			&product.Description,
			&product.ImagePath,
			&product.Seller,
			&product.Status,
			&product.CreatedAt,
		)

		if err != nil {
			utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Could not read product data",
			})
			return
		}

		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetMyProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	userIDText := r.URL.Query().Get("userId")

	if userIDText == "" {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "User ID is required",
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

	rows, err := h.DB.Query(
		`SELECT 
			products.id,
			products.user_id,
			products.title,
			products.category,
			products.price,
			products.location,
			products.contact,
			products.description,
			COALESCE(products.image_path, ''),
			users.full_name,
			products.status,
			products.created_at::text
		FROM products
		JOIN users ON products.user_id = users.id
		WHERE products.user_id = $1
		ORDER BY products.created_at DESC`,
		userID,
	)

	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not fetch your products",
		})
		return
	}

	defer rows.Close()

	products := []models.ProductResponse{}

	for rows.Next() {
		var product models.ProductResponse

		err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.Title,
			&product.Category,
			&product.Price,
			&product.Location,
			&product.Contact,
			&product.Description,
			&product.ImagePath,
			&product.Seller,
			&product.Status,
			&product.CreatedAt,
		)

		if err != nil {
			utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Could not read your product data",
			})
			return
		}

		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	productIDText := r.URL.Query().Get("id")

	if productIDText == "" {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Product ID is required",
		})
		return
	}

	productID, err := strconv.Atoi(productIDText)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid product ID",
		})
		return
	}

	_, err = h.DB.Exec(`DELETE FROM products WHERE id = $1`, productID)

	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not delete product",
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product deleted successfully",
	})
}