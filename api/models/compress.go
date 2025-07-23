package models

import "time"

type CompressJob struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	InputFileID  string    `json:"input_file_id"`
	OutputFileID string    `json:"output_file_id"`
	Compression  string    `json:"compression"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type CompressRequest struct {
	InputFileID string `json:"input_file_id" binding:"required"`
	Compression string `json:"compression" binding:"required"` // "low", "medium", "high"
}
