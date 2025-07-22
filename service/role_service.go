package service

import (
	"context"

	"test/api/models"
	"test/pkg/logger"
	"test/storage"
)

type RoleService interface {
	Create(ctx context.Context, name, createdBy string) (string, error)
	Update(ctx context.Context, id, name string) error
	GetAll(ctx context.Context) ([]models.Role, error)
	Exists(ctx context.Context, id string) (bool, error)
}

type roleService struct {
	stg storage.IRoleStorage
	log logger.ILogger
}

func NewRoleService(stg storage.IStorage, log logger.ILogger) RoleService {
	return &roleService{
		stg: stg.Role(),
		log: log,
	}
}

func (s *roleService) Create(ctx context.Context, name, createdBy string) (string, error) {
	s.log.Info("RoleService.Create called",
		logger.String("name", name),
		logger.String("createdBy", createdBy),
	)

	id, err := s.stg.Create(ctx, name, createdBy)
	if err != nil {
		s.log.Error("Failed to create role", logger.Error(err))
		return "", err
	}

	s.log.Info("Role created successfully", logger.String("roleID", id))
	return id, nil
}

func (s *roleService) Update(ctx context.Context, id, name string) error {
	s.log.Info("RoleService.Update called",
		logger.String("id", id),
		logger.String("new_name", name),
	)

	err := s.stg.Update(ctx, id, name)
	if err != nil {
		s.log.Error("Failed to update role", logger.Error(err))
		return err
	}

	s.log.Info("Role updated successfully", logger.String("id", id))
	return nil
}

func (s *roleService) GetAll(ctx context.Context) ([]models.Role, error) {
	s.log.Info("RoleService.GetAll called")

	roles, err := s.stg.GetAll(ctx)
	if err != nil {
		s.log.Error("Failed to get roles", logger.Error(err))
		return nil, err
	}

	s.log.Info("Roles fetched successfully", logger.Int("count", len(roles)))
	return roles, nil
}

func (s *roleService) Exists(ctx context.Context, id string) (bool, error) {
	s.log.Info("RoleService.Exists called", logger.String("id", id))

	ok, err := s.stg.Exists(ctx, id)
	if err != nil {
		s.log.Error("Failed to check role existence", logger.Error(err))
		return false, err
	}

	return ok, nil
}
