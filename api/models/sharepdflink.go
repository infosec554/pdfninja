package models

import "time"

type CreateSharedLinkRequest struct {
	FileID    string    `json:"file_id" binding:"required"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

type SharedLink struct {
	ID          string    `json:"id"`
	FileID      string    `json:"file_id"`
	SharedToken string    `json:"shared_token"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
