package models

import "time"

// PDFToJPGJob – PDF dan JPGga aylantirish uchun job modeli
type PDFToJPGJob struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	InputFileID string    `json:"input_file_id"`
	OutputPaths []string  `json:"output_paths"` // chiqarilgan JPG fayllar
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// PDFToJPGRequest – mijozdan keladigan request
type PDFToJPGRequest struct {
	InputFileID string `json:"input_file_id" binding:"required"`
}
