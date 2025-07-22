package models

import "time"

type AddPageNumberJob struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	InputFileID  string    `json:"input_file_id"`
	OutputFileID string    `json:"output_file_id"`
	Status       string    `json:"status"` // pending, done, failed
	CreatedAt    time.Time `json:"created_at"`
	Font         string    `json:"font"`
	FontSize     int       `json:"font_size"`
	Position     string    `json:"position"` // bottom-right, top-left, etc.
	FirstNumber  int       `json:"first_number"`
}

type AddPageNumbersRequest struct {
	InputFileID string `json:"input_file_id" binding:"required"`
	Font        string `json:"font"`                   // Arial, Helvetica, va hokazo
	FontSize    int    `json:"font_size"`              // 12, 14, ...
	Position    string `json:"position"`               // bottom-center, top-right, ...
	PageRange   string `json:"page_range"`             // 1-5,7
	Color       string `json:"color"`                  // #000000
	FirstNumber int    `json:"first_number"`           // Boshlang‘ich raqam
}
