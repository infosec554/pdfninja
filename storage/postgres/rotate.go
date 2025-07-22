package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
)

type rotateRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewRotateRepo(db *pgxpool.Pool, log logger.ILogger) *rotateRepo {
	return &rotateRepo{db: db, log: log}
}

func (r *rotateRepo) Create(ctx context.Context, job *models.RotateJob) error {
	query := `
		INSERT INTO rotate_jobs (
			id, user_id, input_file_id, angle, output_file_id, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		job.ID, job.UserID, job.InputFileID, job.Angle,
		job.OutputFileID, job.Status, job.CreatedAt,
	)
	if err != nil {
		r.log.Error("rotateRepo.Create query error", logger.Error(err))
		return err
	}
	return nil
}

func (r *rotateRepo) GetByID(ctx context.Context, id string) (*models.RotateJob, error) {
	query := `
		SELECT id, user_id, input_file_id, angle, output_file_id, status, created_at
		FROM rotate_jobs WHERE id = $1
	`

	var job models.RotateJob
	err := r.db.QueryRow(ctx, query, id).Scan(
		&job.ID, &job.UserID, &job.InputFileID, &job.Angle,
		&job.OutputFileID, &job.Status, &job.CreatedAt,
	)
	if err != nil {
		r.log.Error("rotateRepo.GetByID query error", logger.Error(err))
		return nil, err
	}
	return &job, nil
}

func (r *rotateRepo) Update(ctx context.Context, job *models.RotateJob) error {
	query := `
		UPDATE rotate_jobs
		SET output_file_id = $1, status = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query,
		job.OutputFileID, job.Status, job.ID,
	)
	if err != nil {
		r.log.Error("rotateRepo.Update query error", logger.Error(err))
		return err
	}
	return nil
}
