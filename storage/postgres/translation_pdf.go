package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type translatePDFRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewTranslatePDFRepo(db *pgxpool.Pool, log logger.ILogger) storage.ITranslatePDFStorage {
	return &translatePDFRepo{db: db, log: log}
}

func (r *translatePDFRepo) Create(ctx context.Context, job *models.TranslatePDFJob) error {
	query := `
		INSERT INTO translate_pdf_jobs (
			id, user_id, input_file_id, source_lang, target_lang, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, job.ID, job.UserID, job.InputFileID, job.SourceLang, job.TargetLang, job.Status, job.CreatedAt)
	return err
}

func (r *translatePDFRepo) Update(ctx context.Context, job *models.TranslatePDFJob) error {
	query := `
		UPDATE translate_pdf_jobs
		SET output_file_id = $1, status = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, job.OutputFileID, job.Status, job.ID)
	return err
}

func (r *translatePDFRepo) GetByID(ctx context.Context, id string) (*models.TranslatePDFJob, error) {
	query := `
		SELECT id, user_id, input_file_id, source_lang, target_lang, output_file_id, status, created_at
		FROM translate_pdf_jobs
		WHERE id = $1
	`

	var job models.TranslatePDFJob
	err := r.db.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.UserID,
		&job.InputFileID,
		&job.SourceLang,
		&job.TargetLang,
		&job.OutputFileID,
		&job.Status,
		&job.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
