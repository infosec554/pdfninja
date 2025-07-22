package models

import "time"

// CreateJpgToPdfRequest – mijoz yuboradigan so‘rov
type CreateJpgToPdfRequest struct {
	ImageFileIDs []string `json:"image_file_ids" binding:"required"` // JPG fayllarning ID’lari
}

// JpgToPdfJob – JPG fayllarni PDF’ga aylantirish jarayonining modeli
type JpgToPdfJob struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	ImageFileIDs []string  `json:"image_file_ids"`
	OutputFileID string    `json:"output_file_id"`
	Status       string    `json:"status"` // pending, done, failed
	CreatedAt    time.Time `json:"created_at"`
}
