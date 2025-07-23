package models

import "time"

// ExtractJob – PDF sahifalarni ajratish jarayonining holatini ifodalaydi
type ExtractJob struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	InputFileID    string    `json:"input_file_id"`
	PagesToExtract string    `json:"pages_to_extract"`
	OutputFileID   *string   `json:"output_file_id,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// ExtractPagesRequest – mijoz yuboradigan request struct
type ExtractPagesRequest struct {
	InputFileID string `json:"input_file_id" binding:"required"`
	PageRanges  string `json:"page_ranges" binding:"required"` // Masalan: "2-4,6"
}
