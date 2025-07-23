package models

import "time"

type OrganizeJob struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	InputFileID  string    `json:"input_file_id"`
	NewOrder     []int     `json:"new_order"`      // Sahifalarning yangi tartibi, masalan: [3,1,2]
	OutputFileID string    `json:"output_file_id"` // Natija fayl ID
	Status       string    `json:"status"`         // pending, processing, done, failed
	CreatedAt    time.Time `json:"created_at"`     // Yaratilgan vaqt
}

type CreateOrganizeJobRequest struct {
	InputFileID string `json:"input_file_id" example:"f680eafd-4c56-435d-8a3b-60f0d15b7f4e"`
	NewOrder    []int  `json:"new_order" example:"[3,1,2]"` // JSON array shaklida keladi
}
