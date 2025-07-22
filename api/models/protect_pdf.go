package models

import "time"

// ProtectPDFRequest – Kiruvchi so‘rov modeli
type ProtectPDFRequest struct {
	InputFileID string `json:"input_file_id" binding:"required"` // PDF fayl ID
	Password    string `json:"password" binding:"required"`      // Qo‘yiladigan parol
}

// ProtectPDFJob – Bazadagi ish (job) modeli
type ProtectPDFJob struct {
	ID           string    `json:"id"`             // Job ID
	UserID       string    `json:"user_id"`        // Foydalanuvchi IDsi
	InputFileID  string    `json:"input_file_id"`  // Kiruvchi fayl
	OutputFileID string    `json:"output_file_id"` // Himoyalangan fayl IDsi
	Password     string    `json:"password"`       // Qo‘yilgan parol
	Status       string    `json:"status"`         // pending, done, failed
	CreatedAt    time.Time `json:"created_at"`     // Yaralgan vaqti
}
