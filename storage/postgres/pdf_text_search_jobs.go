package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type pdfTextSearchRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewPDFTextSearchRepo(db *pgxpool.Pool, log logger.ILogger) storage.IPDFTextSearchStorage {
	return &pdfTextSearchRepo{db: db, log: log}
}

func (r *pdfTextSearchRepo) Create(ctx context.Context, job *models.PDFTextSearchJob) error {
	query := `INSERT INTO pdf_text_search_jobs (id, user_id, input_file_id, status, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, job.ID, job.UserID, job.InputFileID, job.Status, job.CreatedAt)
	return err
}

func (r *pdfTextSearchRepo) Update(ctx context.Context, job *models.PDFTextSearchJob) error {
	query := `UPDATE pdf_text_search_jobs SET extracted_text = $1, status = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, job.ExtractedText, job.Status, job.ID)
	return err
}

func (r *pdfTextSearchRepo) GetByID(ctx context.Context, id string) (*models.PDFTextSearchJob, error) {
	var job models.PDFTextSearchJob
	query := `SELECT id, user_id, input_file_id, extracted_text, status, created_at FROM pdf_text_search_jobs WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&job.ID, &job.UserID, &job.InputFileID, &job.ExtractedText, &job.Status, &job.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
