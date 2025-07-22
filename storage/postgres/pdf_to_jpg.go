package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type pdfToJPGRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewPDFToJPGRepo(db *pgxpool.Pool, log logger.ILogger) storage.IPDFToJPGStorage {
	return &pdfToJPGRepo{db: db, log: log}
}

func (r *pdfToJPGRepo) Create(ctx context.Context, job *models.PDFToJPGJob) error {
	query := `
		INSERT INTO pdf_to_jpg_jobs (id, user_id, input_file_id, output_paths, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, job.ID, job.UserID, job.InputFileID, job.OutputPaths, job.Status, job.CreatedAt)
	return err
}

func (r *pdfToJPGRepo) Update(ctx context.Context, job *models.PDFToJPGJob) error {
	query := `
		UPDATE pdf_to_jpg_jobs
		SET output_paths = $1, status = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, job.OutputPaths, job.Status, job.ID)
	return err
}

func (r *pdfToJPGRepo) GetByID(ctx context.Context, id string) (*models.PDFToJPGJob, error) {
	query := `
		SELECT id, user_id, input_file_id, output_paths, status, created_at
		FROM pdf_to_jpg_jobs
		WHERE id = $1
	`

	var job models.PDFToJPGJob
	err := r.db.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.UserID,
		&job.InputFileID,
		&job.OutputPaths,
		&job.Status,
		&job.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
