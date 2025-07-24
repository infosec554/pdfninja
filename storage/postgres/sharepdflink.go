package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type sharedLinkRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewSharedLinkRepo(db *pgxpool.Pool, log logger.ILogger) storage.ISharedLinkStorage {
	return &sharedLinkRepo{db: db, log: log}
}

func (r *sharedLinkRepo) Create(ctx context.Context, req *models.SharedLink) error {
	query := `
		INSERT INTO shared_links (id, file_id, shared_token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query,
		req.ID,
		req.FileID,
		req.SharedToken,
		req.ExpiresAt,
		req.CreatedAt,
	)
	return err
}

func (r *sharedLinkRepo) GetByToken(ctx context.Context, token string) (*models.SharedLink, error) {
	query := `
		SELECT id, file_id, shared_token, expires_at, created_at
		FROM shared_links
		WHERE shared_token = $1
	`

	var link models.SharedLink
	err := r.db.QueryRow(ctx, query, token).Scan(
		&link.ID,
		&link.FileID,
		&link.SharedToken,
		&link.ExpiresAt,
		&link.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &link, nil
}
