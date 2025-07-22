package models

import "time"

type OrganizeJob struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	InputFileID  string    `json:"input_file_id"`
	NewOrder     []int     `json:"new_order"` // sahifa tartibi masalan: [3,1,2]
	OutputFileID string    `json:"output_file_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateOrganizeJobRequest struct {
	InputFileID string `json:"input_file_id" example:"abc-123"`
	NewOrder    []int  `json:"new_order" example:"3,1,2"` // JSON orqali keladi, keyin []int ga aylantiriladi
}
