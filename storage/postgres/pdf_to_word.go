package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type pdfToWordRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewPDFToWordRepo(db *pgxpool.Pool, log logger.ILogger) storage.IPDFToWordStorage {
	return &pdfToWordRepo{db: db}
}

func (r *pdfToWordRepo) Create(ctx context.Context, job *models.PDFToWordJob) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO pdf_to_word_jobs (id, user_id, input_file_id, status, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		job.ID, job.UserID, job.InputFileID, job.Status, job.CreatedAt)
	return err
}

func (r *pdfToWordRepo) GetByID(ctx context.Context, id string) (*models.PDFToWordJob, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, user_id, input_file_id, output_path, status, created_at FROM pdf_to_word_jobs WHERE id=$1`, id)

	var job models.PDFToWordJob
	err := row.Scan(&job.ID, &job.UserID, &job.InputFileID, &job.OutputPath, &job.Status, &job.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *pdfToWordRepo) Update(ctx context.Context, job *models.PDFToWordJob) error {
	_, err := r.db.Exec(ctx,
		`UPDATE pdf_to_word_jobs SET output_path=$1, status=$2 WHERE id=$3`,
		job.OutputPath, job.Status, job.ID)
	return err
}
