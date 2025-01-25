package common

import "time"

type UserCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID string `json:"id"`
	UserCredentials
	Balance
}

type Order struct {
	ID         string    `json:"id"`
	Number     string    `json:"number"`
	UploadedAt time.Time `json:"uploaded_at"`
	UserID     string    `json:"user_id"`
	Accrual    float64   `json:"accrual"`
	Status     string    `json:"status"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawal struct {
	ID          string    `json:"id"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
