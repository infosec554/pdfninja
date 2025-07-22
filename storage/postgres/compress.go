package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
)

type compressRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewCompressRepo(db *pgxpool.Pool, log logger.ILogger) *compressRepo {
	return &compressRepo{db: db, log: log}
}

func (r *compressRepo) Create(ctx context.Context, job *models.CompressJob) error {
	query := `
		INSERT INTO compress_jobs (id, user_id, input_file_id, compression, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		job.ID, job.UserID, job.InputFileID, job.Compression, job.Status, job.CreatedAt)
	return err
}

func (r *compressRepo) Update(ctx context.Context, job *models.CompressJob) error {
	query := `
		UPDATE compress_jobs
		SET output_file_id = $1, status = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, job.OutputFileID, job.Status, job.ID)
	return err
}

func (r *compressRepo) GetByID(ctx context.Context, id string) (*models.CompressJob, error) {
	query := `
		SELECT id, user_id, input_file_id, compression, output_file_id, status, created_at
		FROM compress_jobs
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var job models.CompressJob
	err := row.Scan(
		&job.ID,
		&job.UserID,
		&job.InputFileID,
		&job.Compression,
		&job.OutputFileID,
		&job.Status,
		&job.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
