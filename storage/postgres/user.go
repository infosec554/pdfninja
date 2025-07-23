package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type userRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewUserRepo(db *pgxpool.Pool, log logger.ILogger) storage.IUserStorage {
	return &userRepo{
		db:  db,
		log: log,
	}
}

func (r *userRepo) Create(ctx context.Context, req models.CreateUser) (string, error) {
	id := uuid.New().String()

	query := `
		INSERT INTO users (id, name, email, password, status)
		VALUES ($1, $2, $3, $4, 'active')
	`

	_, err := r.db.Exec(ctx, query, id, req.Name, req.Email, req.Password)
	if err != nil {
		r.log.Error("error inserting user", logger.Error(err))
		return "", err
	}

	return id, nil
}

func (r *userRepo) GetForLoginByEmail(ctx context.Context, email string) (models.LoginUser, error) {
	var user models.LoginUser

	query := `
		SELECT id, password, status
		FROM users
		WHERE email = $1 AND status = 'active'
	`

	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Password, &user.Status)
	if err != nil {
		r.log.Error("failed to get user by email", logger.Error(err))
		return models.LoginUser{}, err
	}

	return user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, name, email, status, created_at
		FROM users
		WHERE id = $1 AND status = 'active'
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
	)
	if err != nil {
		r.log.Error("failed to get user by ID", logger.Error(err))
		return nil, err
	}

	return &user, nil
}
