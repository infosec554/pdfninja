package models

import "time"

type File struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	FileType   string    `json:"file_type"`
	FileSize   int64     `json:"file_size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type FileUploadRequest struct {
	UserID   string
	FileName string
	FileData []byte
}
