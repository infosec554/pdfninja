package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
	"test/pkg/jwt"
)

// SendOTP godoc
// @Router       /otp/send [POST]
// @Security     ApiKeyAuth
// @Summary      Send OTP to email
// @Tags         otp
// @Accept       json
// @Produce      json
// @Param        request body models.SendOtpRequest true "email info"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) SendOTP(c *gin.Context) {
	var req models.SendOtpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	otpID, err := h.services.Otp().SendOtp(ctx, req.Email)
	if err != nil {
		handleResponse(c, h.log, "failed to send OTP", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "OTP sent successfully", http.StatusOK, gin.H{"otp_id": otpID})
}

// ConfirmOTP godoc
// @Router       /otp/confirm [POST]
// @Security     ApiKeyAuth
// @Summary      Confirm OTP code
// @Tags         otp
// @Accept       json
// @Produce      json
// @Param        request body models.ConfirmOtpRequest true "otp info"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  models.Response
// @Failure      401  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) ConfirmOTP(c *gin.Context) {
	var req models.ConfirmOtpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	email, savedCode, expiresAt, err := h.services.Otp().GetUnconfirmedByID(ctx, req.OtpID)
	if err != nil {
		handleResponse(c, h.log, "OTP not found", http.StatusUnauthorized, err.Error())
		return
	}

	if time.Now().After(expiresAt) {
		handleResponse(c, h.log, "OTP expired", http.StatusUnauthorized, nil)
		return
	}

	if req.Code != savedCode {
		handleResponse(c, h.log, "Incorrect OTP", http.StatusUnauthorized, nil)
		return
	}

	if err := h.services.Otp().UpdateStatusToConfirmed(ctx, req.OtpID); err != nil {
		handleResponse(c, h.log, "Failed to confirm OTP", http.StatusInternalServerError, err.Error())
		return
	}

	// ✅ JWT tokenni TO‘G‘RI email bilan generate qilamiz
	claims := map[string]interface{}{
		"otp_id": req.OtpID,
		"email":  email,
	}

	token, _, err := jwt.GenerateJWT(claims)
	if err != nil {
		handleResponse(c, h.log, "Failed to generate token", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "OTP confirmed", http.StatusOK, gin.H{"otp_confirmation_token": token})
}
