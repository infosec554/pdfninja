package models

import "time"

type PDFToWordRequest struct {
	InputFileID string `json:"input_file_id" binding:"required"`
}

type PDFToWordJob struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	InputFileID  string    `json:"input_file_id"`
	OutputFileID string    `json:"output_file_id"` // to‘g‘rilangan
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
