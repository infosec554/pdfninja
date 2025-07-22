package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type jpgToPdfRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewJpgToPdfRepo(db *pgxpool.Pool, log logger.ILogger) storage.IJpgToPdfStorage {
	return &jpgToPdfRepo{db: db, log: log}
}

func (r *jpgToPdfRepo) Create(ctx context.Context, job *models.JpgToPdfJob) error {
	query := `
		INSERT INTO jpg_to_pdf_jobs (id, user_id, image_file_ids, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, job.ID, job.UserID, job.ImageFileIDs, job.Status, job.CreatedAt)
	return err
}

func (r *jpgToPdfRepo) Update(ctx context.Context, job *models.JpgToPdfJob) error {
	query := `
		UPDATE jpg_to_pdf_jobs
		SET output_file_id = $1, status = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, job.OutputFileID, job.Status, job.ID)
	return err
}

func (r *jpgToPdfRepo) GetByID(ctx context.Context, id string) (*models.JpgToPdfJob, error) {
	query := `
		SELECT id, user_id, image_file_ids, output_file_id, status, created_at
		FROM jpg_to_pdf_jobs
		WHERE id = $1
	`

	var job models.JpgToPdfJob
	err := r.db.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.UserID,
		&job.ImageFileIDs,
		&job.OutputFileID,
		&job.Status,
		&job.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}
