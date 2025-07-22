package models

import "time"

// ExtractJob – PDF sahifalarni ajratish jarayonining holatini ifodalaydi
type ExtractJob struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	InputFileID   string    `json:"input_file_id"`
	PageRanges    string    `json:"page_ranges"`     // Masalan: "1-2,5,7"
	OutputFileIDs []string  `json:"output_file_ids"` // Ajratilgan fayllarning ID’lari
	Status        string    `json:"status"`          // pending, done, failed
	CreatedAt     time.Time `json:"created_at"`
}

// ExtractPagesRequest – mijoz yuboradigan request struct
type ExtractPagesRequest struct {
	InputFileID string `json:"input_file_id" binding:"required"`
	PageRanges  string `json:"page_ranges" binding:"required"` // Masalan: "2-4,6"
}
