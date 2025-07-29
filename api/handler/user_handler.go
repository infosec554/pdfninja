package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test/api/models"
	"test/pkg/jwt"
	"test/pkg/security"
)

// SignUp godoc
// @Summary      Create new user
// @Description  Register a new user after OTP confirmation
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body models.CreateUser true "Signup data"
// @Success      201 {object} map[string]string
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /signup [post]
// @Security     ApiKeyAuth
func (h Handler) SignUp(c *gin.Context) {
	var req models.CreateUser

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	// 1. Parse OTP token (JWT)
	claims, err := jwt.ParseToken(req.OtpToken)
	if err != nil {
		handleResponse(c, h.log, "invalid otp token", http.StatusUnauthorized, err.Error())
		return
	}

	email, emailOk := claims["email"].(string)
	otpID, otpOk := claims["otp_id"].(string)
	if !emailOk || !otpOk || email != req.Email {
		handleResponse(c, h.log, "email mismatch or token malformed", http.StatusUnauthorized, nil)
		return
	}

	// 2. Redis orqali OTP kodni tekshiramiz
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Parolni hash qilamiz
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		handleResponse(c, h.log, "failed to hash password", http.StatusInternalServerError, err.Error())
		return
	}
	req.Password = hashedPassword

	// 4. Foydalanuvchini yaratamiz
	userID, err := h.services.User().Create(ctx, req)
	if err != nil {
		handleResponse(c, h.log, "failed to create user", http.StatusInternalServerError, err.Error())
		return
	}

	// 5. OTP ni tasdiqlangan deb belgilaymiz
	_ = h.services.Otp().UpdateStatusToConfirmed(ctx, otpID)

	handleResponse(c, h.log, "user created successfully", http.StatusCreated, gin.H{"user_id": userID})
}

// Login godoc
// @Summary      Login user
// @Description  User login with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.LoginRequest true "Login credentials"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /login [post]
// @Security     ApiKeyAuth
func (h Handler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid login request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.services.User().GetForLoginByEmail(ctx, req.Email)
	if err != nil {
		handleResponse(c, h.log, "user not found", http.StatusUnauthorized, err.Error())
		return
	}

	if err := security.CompareHashAndPassword(user.Password, req.Password); err != nil {
		handleResponse(c, h.log, "invalid credentials", http.StatusUnauthorized, "email or password is incorrect")
		return
	}

	// ✅ TO‘G‘RI CLAIMS
	claims := map[string]interface{}{
		"user_id":   user.ID,
		"user_role": "user", // yoki user.Role agar mavjud bo‘lsa
	}

	accessToken, refreshToken, err := jwt.GenerateJWT(claims)
	if err != nil {
		handleResponse(c, h.log, "failed to generate token", http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	handleResponse(c, h.log, "login successful", http.StatusOK, resp)
}

// GetMyProfile godoc
// @Summary      Get my profile
// @Description  Foydalanuvchining o‘z profilini olish (JWT token orqali)
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200 {object} models.User
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /me [get]
// @Security     ApiKeyAuth
func (h *Handler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	user, err := h.services.User().GetByID(c.Request.Context(), userID.(string))
	if err != nil {
		handleResponse(c, h.log, "failed to get user", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "user profile", http.StatusOK, user)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Return a new access and refresh token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh body models.RefreshTokenRequest true "Refresh token"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /refresh-token [post]
// @Security     ApiKeyAuth
func (h Handler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "refresh_token is required", http.StatusBadRequest, err.Error())
		return
	}

	claims, err := jwt.ParseToken(req.RefreshToken)
	if err != nil {
		handleResponse(c, h.log, "invalid refresh token", http.StatusUnauthorized, err.Error())
		return
	}

	userID, ok1 := claims["user_id"].(string)
	userRole, ok2 := claims["user_role"].(string)

	if !ok1 || !ok2 || userID == "" {
		handleResponse(c, h.log, "invalid claims in refresh token", http.StatusUnauthorized, nil)
		return
	}

	accessToken, newRefreshToken, err := jwt.GenerateJWT(map[string]interface{}{
		"user_id":   userID,
		"user_role": userRole,
	})
	if err != nil {
		handleResponse(c, h.log, "failed to generate new tokens", http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	handleResponse(c, h.log, "tokens refreshed", http.StatusOK, resp)
}
