package models

type ProductRequest struct {
	UserID      int    `json:"userId"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Price       string `json:"price"`
	Location    string `json:"location"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
}

type ProductResponse struct {
	ID          int    `json:"id"`
	UserID      int    `json:"userId"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Price       string `json:"price"`
	Location    string `json:"location"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
	ImagePath   string `json:"imagePath"`
	Seller      string `json:"seller"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
}