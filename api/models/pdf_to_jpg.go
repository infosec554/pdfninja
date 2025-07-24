package models

import "time"

// PDFToJPGJob – PDF dan JPGga aylantirish uchun job modeli
type PDFToJPGJob struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	InputFileID   string    `json:"input_file_id"`
	OutputFileIDs []string  `json:"output_file_ids"`
	ZipFileID     *string   `json:"zip_file_id"` // <-- nil bo'lishi mumkin
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// PDFToJPGRequest – mijozdan keladigan request
type PDFToJPGRequest struct {
	InputFileID string `json:"input_file_id" binding:"required"`
}
