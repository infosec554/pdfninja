package models

import "time"

type MergeJob struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	OutputFileID *string   `db:"output_file_id"`
	Status       string    `db:"status"` // pending, done, failed
	CreatedAt    time.Time `db:"created_at"`
	InputFileIDs []string  `db:"-"` // populated manually via join from merge_job_input_files
}
type MergeJobInputFile struct {
	ID     string `db:"id"`
	JobID  string `db:"job_id"`
	FileID string `db:"file_id"`
}
type CreateMergeJobRequest struct {
	UserID       string   `json:"user_id"`
	InputFileIDs []string `json:"input_file_ids"`
}
