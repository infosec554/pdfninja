package service

import (
	"context"
	"fmt"
	"time"

	"test/pkg/logger"
	"test/pkg/mailer"
	"test/storage"
)

type OtpService interface {
	SendOtp(ctx context.Context, email string) (string, error)
	GetUnconfirmedByID(ctx context.Context, id string) (string, string, time.Time, error)
	UpdateStatusToConfirmed(ctx context.Context, id string) error
	GetByIDAndEmail(ctx context.Context, id, email string) (bool, error)
}

type otpService struct {
	stg    storage.IOTPStorage
	log    logger.ILogger
	mailer *mailer.Mailer
	redis  storage.IRedisStorage
}

func NewOtpService(stg storage.IStorage, log logger.ILogger, mailer *mailer.Mailer, redis storage.IRedisStorage) OtpService {
	return &otpService{
		stg:    stg.OTP(),
		log:    log,
		mailer: mailer,
		redis:  redis,
	}
}

func (s *otpService) SendOtp(ctx context.Context, email string) (string, error) {
	s.log.Info("OtpService.SendOtp called", logger.String("email", email))

	code := generateCode()
	expiresAt := time.Now().Add(5 * time.Minute)

	id, err := s.stg.Create(ctx, email, code, expiresAt)
	if err != nil {
		s.log.Error("Failed to create OTP", logger.Error(err))
		return "", err
	}

	redisKey := "otp:" + id
	if err := s.redis.SetX(ctx, redisKey, code, 5*time.Minute); err != nil {
		s.log.Error("Failed to set Redis", logger.Error(err))
	}

	body := fmt.Sprintf(`
<html>
  <head>
    <style>
      .container {
        max-width: 500px;
        margin: auto;
        padding: 20px;
        border: 1px solid #eee;
        border-radius: 10px;
        font-family: Arial, sans-serif;
        background-color: #f9f9f9;
      }
      .otp-code {
        font-size: 24px;
        font-weight: bold;
        color: #2c3e50;
        margin: 20px 0;
      }
      .footer {
        font-size: 12px;
        color: #888;
        margin-top: 30px;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <h2>Assalomu alaykum!</h2>
      <p>Sizning OTP kodingiz:</p>
      <div class="otp-code">%s</div>
      <p>Iltimos, bu kodni 5 daqiqa ichida kiriting.</p>
      <div class="footer">Agar siz bu so‘rovni yubormagan bo‘lsangiz, xabarni e’tiborsiz qoldiring.</div>
    </div>
  </body>
</html>`, code)
	if err := s.mailer.Send(email, "Your OTP Code", body); err != nil {
		s.log.Error("Failed to send email", logger.Error(err))
		return "", err
	}

	s.log.Info("OTP created and sent", logger.String("otpID", id))
	return id, nil
}

func (s *otpService) GetUnconfirmedByID(ctx context.Context, id string) (string, string, time.Time, error) {
	s.log.Info("OtpService.GetUnconfirmedByID called", logger.String("otpID", id))

	redisKey := "otp:" + id
	if val, err := s.redis.Get(ctx, redisKey); err == nil && val != "" {
		email, _, expiresAt, err := s.stg.GetUnconfirmedByID(ctx, id)
		if err != nil {
			s.log.Error("DB fallback failed in GetUnconfirmedByID", logger.Error(err))
			return "", "", time.Time{}, err
		}
		return email, val, expiresAt, nil
	}

	return s.stg.GetUnconfirmedByID(ctx, id)
}

func (s *otpService) UpdateStatusToConfirmed(ctx context.Context, id string) error {
	s.log.Info("OtpService.UpdateStatusToConfirmed called", logger.String("otpID", id))
	return s.stg.UpdateStatusToConfirmed(ctx, id)
}

func (s *otpService) GetByIDAndEmail(ctx context.Context, id, email string) (bool, error) {
	s.log.Info("OtpService.GetByIDAndEmail called", logger.String("otpID", id), logger.String("email", email))
	return s.stg.GetByIDAndEmail(ctx, id, email)
}

func generateCode() string {
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}
