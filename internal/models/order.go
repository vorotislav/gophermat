package models

import "time"

type Order struct {
	ID         int64     `json:"ID"`
	UserID     int       `json:"user_id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}
