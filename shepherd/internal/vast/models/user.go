package models

type User struct {
	Email    string  `json:"email"`
	FullName string  `json:"fullname"`
	Credit   float64 `json:"credit"`
}
