package models

import "time"

type PDFToWordRequest struct {
	InputFileID string `json:"input_file_id" binding:"required"`
}

type PDFToWordJob struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	InputFileID string    `json:"input_file_id"`
	OutputPath  string    `json:"output_path"`
	Status      string    `json:"status"` // pending, done, failed
	CreatedAt   time.Time `json:"created_at"`
}
