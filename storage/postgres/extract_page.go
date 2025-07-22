package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type extractRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewExtractPageRepo(db *pgxpool.Pool, log logger.ILogger) storage.IExtractPageStorage {
	return &extractRepo{
		db:  db,
		log: log,
	}
}

func (r *extractRepo) Create(ctx context.Context, job *models.ExtractJob) error {
	query := `
		INSERT INTO extract_jobs (id, user_id, input_file_id, page_ranges, output_file_ids, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		job.ID, job.UserID, job.InputFileID, job.PageRanges, job.OutputFileIDs, job.Status, job.CreatedAt)
	return err
}

func (r *extractRepo) Update(ctx context.Context, job *models.ExtractJob) error {
	query := `
		UPDATE extract_jobs
		SET output_file_ids = $1, status = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, job.OutputFileIDs, job.Status, job.ID)
	return err
}

func (r *extractRepo) GetByID(ctx context.Context, id string) (*models.ExtractJob, error) {
	query := `
		SELECT id, user_id, input_file_id, page_ranges, output_file_ids, status, created_at
		FROM extract_jobs
		WHERE id = $1
	`

	var job models.ExtractJob
	err := r.db.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.UserID,
		&job.InputFileID,
		&job.PageRanges,
		&job.OutputFileIDs,
		&job.Status,
		&job.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}
