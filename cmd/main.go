package main

import (
	"context"

	"test/api"
	"test/config"
	"test/pkg/logger"
	"test/pkg/mailer"
	"test/service"
	"test/storage/postgres"
	"test/storage/redis"
)

func main() {
	// 1. Load config
	cfg := config.Load()

	// 2. Logger yaratish
	log := logger.New(cfg.ServiceName)

	// 3. Postgresga ulanish
	pgStore, err := postgres.New(context.Background(), cfg, log, nil)
	if err != nil {
		log.Error("error while connecting to db", logger.Error(err))
		return
	}
	defer pgStore.Close()

	// 4. Mailer yaratish
	mailService := mailer.New(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
		cfg.SMTPSenderName,
	)

	redis := redis.New(cfg) // storage/redis package'idan

	// 5. Servicelarni ulash (Redis hozircha kerak emas)
	services := service.New(pgStore, log, mailService, redis)

	// 6. API serverni ishga tushurish
	server := api.New(services, log)

	log.Info("Service is running on", logger.Int("port", 8080))
	if err = server.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
