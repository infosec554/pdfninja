package service

import (
	"context"

	"test/pkg/logger"
	"test/storage"
)

type SysUserService interface {
	GetByPhone(ctx context.Context, phone string) (id, hashedPassword, status string, err error)
	Create(ctx context.Context, name, phone, hashedPassword, createdBy string) (string, error)
	AssignRoles(ctx context.Context, sysuserID string, roleIDs []string) error
}


type sysuserService struct {
	stg storage.ISysuserStorage
	log logger.ILogger
}


func NewSysUserService(stg storage.IStorage, log logger.ILogger) SysUserService {
	return &sysuserService{
		stg: stg.Sysuser(),
		log: log,
	}
}


func (s *sysuserService) GetByPhone(ctx context.Context, phone string) (string, string, string, error) {
	s.log.Info("SysUserService.GetByPhone called", logger.String("phone", phone))

	id, hashedPassword, status, err := s.stg.GetByPhone(ctx, phone)
	if err != nil {
		s.log.Error("Failed to get sysuser by phone", logger.Error(err))
	}
	return id, hashedPassword, status, err
}

func (s *sysuserService) Create(ctx context.Context, name, phone, hashedPassword, createdBy string) (string, error) {
	s.log.Info("SysUserService.Create called",
		logger.String("name", name),
		logger.String("phone", phone),
		logger.String("createdBy", createdBy),
	)

	id, err := s.stg.Create(ctx, name, phone, hashedPassword, createdBy)
	if err != nil {
		s.log.Error("Failed to create sysuser", logger.Error(err))
		return "", err
	}

	s.log.Info("SysUser created successfully", logger.String("sysuserID", id))
	return id, nil
}

func (s *sysuserService) AssignRoles(ctx context.Context, sysuserID string, roleIDs []string) error {
	s.log.Info("SysUserService.AssignRoles called",
		logger.String("sysuserID", sysuserID),
		logger.Int("roleCount", len(roleIDs)),
	)

	err := s.stg.AssignRoles(ctx, sysuserID, roleIDs)
	if err != nil {
		s.log.Error("Failed to assign roles to sysuser", logger.Error(err))
		return err
	}

	s.log.Info("Roles successfully assigned to sysuser", logger.String("sysuserID", sysuserID))
	return nil
}
