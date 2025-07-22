package service

import (
	"context"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type UserService interface {
	Create(ctx context.Context, req models.CreateUser) (string, error)
	GetForLoginByEmail(context.Context, string) (models.LoginUser, error) // ✅ yangi nom
}

type userService struct {
	stg storage.IUserStorage
	log logger.ILogger
}

func NewUserService(stg storage.IStorage, log logger.ILogger) UserService {
	return &userService{
		stg: stg.User(),
		log: log,
	}
}

func (s *userService) Create(ctx context.Context, req models.CreateUser) (string, error) {
	s.log.Info("UserService.Create called", logger.String("email", req.Email))

	id, err := s.stg.Create(ctx, req)
	if err != nil {
		s.log.Error("failed to create user", logger.Error(err))
		return "", err
	}

	s.log.Info("user successfully created", logger.String("userID", id))
	return id, nil
}

func (s *userService) GetForLoginByEmail(ctx context.Context, email string) (models.LoginUser, error) {
	s.log.Info("UserService.GetForLoginByEmail called", logger.String("email", email))

	user, err := s.stg.GetForLoginByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user for login", logger.Error(err))
		return models.LoginUser{}, err
	}

	s.log.Info("user fetched for login", logger.String("userID", user.ID))
	return user, nil
}
