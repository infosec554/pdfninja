package models

import "time"

// RotatePDFRequest – HTTP orqali keladigan so‘rov (request) modeli
type RotatePDFRequest struct {
	InputFileID string `json:"input_file_id" binding:"required"`
	Angle       int    `json:"angle" binding:"required"` // Burchak: 90, 180, 270
	Pages       string `json:"pages" binding:"required"` // Sahifa raqamlari masalan: "1-3", "2", "odd"
}
type RotateJob struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	InputFileID  string `json:"input_file_id"`
	Angle        int    `json:"angle"`
	Pages        string `json:"pages"` // Aylantiriladigan sahifalar
	OutputFileID string `json:"output_file_id"`
	OutputPath   string `json:"output_path"` // ✅ YANGI QO‘SHILGAN FIELD
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
