package models

import "time"

type Laptop struct {
	ID        int    `db:"id" json:"id"`
	Serial    string `db:"serial" json:"serial"`
	Brand     string `db:"brand" json:"brand"`
	Model     string `db:"model" json:"model"`
	Available bool   `db:"available" json:"available"`
}
type LaptopRequest struct {
	Serial string `json:"serial" binding:"required,max=100"`
	Brand  string `json:"brand" binding:"required,max=100"`
	Model  string `json:"model" binding:"required,max=100"`
}
type LaptopLoanRequest struct {
	LaptopID int `json:"laptop_id" binding:"required,gt=0"`
}
type LaptopLoan struct {
	ID         int        `db:"id" json:"id"`
	UserID     int        `db:"user_id" json:"user_id"`
	LaptopID   int        `db:"laptop_id" json:"laptop_id"`
	LoanedAt   time.Time  `db:"loaned_at" json:"loaned_at"`
	DueAt      time.Time  `db:"due_at" json:"due_at"`
	ReturnedAt *time.Time `db:"returned_at" json:"returned_at,omitempty"`
	Status     string     `db:"status" json:"status"`
}
