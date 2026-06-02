package models

type AdminLoginRequest struct {
	Password string `json:"password"`
}
type ClearTestDataRequest struct {
	Password string `json:"password"`
}