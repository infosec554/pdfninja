package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"convertpdfgo/api/models"
	"convertpdfgo/pkg/jwt"
	"convertpdfgo/pkg/security"
)

// SignUp godoc
// @Summary      Register a new user
// @Description  Register a new user (name, email, password)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body models.SignupRequest true "Signup data"
// @Success      201 {object} models.SignupResponse
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /signup [post]
func (h Handler) SignUp(c *gin.Context) {
	var req models.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Parolni hash qilish
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		handleResponse(c, h.log, "failed to hash password", http.StatusInternalServerError, err.Error())
		return
	}
	req.Password = hashedPassword

	// Foydalanuvchini yaratish
	userID, err := h.services.User().Create(ctx, req)
	if err != nil {
		handleResponse(c, h.log, "failed to create user", http.StatusInternalServerError, err.Error())
		return
	}

	// UserID va xohlasangiz token qaytarishingiz mumkin
	handleResponse(c, h.log, "user created successfully", http.StatusCreated, models.SignupResponse{
		ID: userID,
	})
}

// Login godoc
// @Summary      User login
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

	claims := map[string]interface{}{
		"user_id":   user.ID,
		"user_role": user.Role, // User yoki admin
	}

	accessToken, refreshToken, err := jwt.GenerateJWT(claims)
	if err != nil {
		handleResponse(c, h.log, "failed to generate token", http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.LoginResponse{
		ID:           user.ID,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	handleResponse(c, h.log, "login successful", http.StatusOK, resp)
}

// GetMyProfile godoc
// @Summary      Get my profile
// @Description  Get user profile (JWT token required)
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200 {object} models.User
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /me [get]
// @Security ApiKeyAuth
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
		ID:           userID,
		Role:         userRole,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	handleResponse(c, h.log, "tokens refreshed", http.StatusOK, resp)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password (user must send old and new password)
// @Tags user
// @Accept json
// @Produce json
// @Param change_password body models.ChangePasswordRequest true "Change password"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /change-password [post]
// @Security ApiKeyAuth
func (h Handler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handleResponse(c, h.log, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.services.User().ChangePassword(ctx, userID.(string), req.OldPassword, req.NewPassword)
	if err != nil {
		handleResponse(c, h.log, err.Error(), http.StatusBadRequest, nil)
		return
	}

	handleResponse(c, h.log, "password changed successfully", http.StatusOK, nil)
}

// @Summary      Google orqali login yoki registratsiya
// @Description  Google OAuth code orqali login yoki ro‘yxatdan o‘tish (JWT tokenlar qaytaradi)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body models.GoogleAuthRequest true "Google authorization code"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /auth/google [post]
func (h Handler) GoogleAuth(c *gin.Context) {
	var req models.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	googleUser, err := h.services.Google().ExchangeCodeForUser(c.Request.Context(), req.Code)
	if err != nil {
		handleResponse(c, h.log, "Google login failed", http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := h.services.User().GoogleAuth(c.Request.Context(), googleUser.Email, googleUser.Name, googleUser.GoogleID)
	if err != nil {
		handleResponse(c, h.log, "failed to create/login user", http.StatusInternalServerError, err.Error())
		return
	}

	token, refresh, err := jwt.GenerateJWT(map[string]interface{}{
		"user_id":   userID,
		"user_role": "user",
	})
	if err != nil {
		handleResponse(c, h.log, "failed to generate token", http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.LoginResponse{
		ID:           userID,
		AccessToken:  token,
		RefreshToken: refresh,
		Role:         "user",
	}
	handleResponse(c, h.log, "login via google", http.StatusOK, resp)
}

// @Summary GitHub orqali login yoki registratsiya
// @Description GitHub OAuth code orqali login yoki ro‘yxatdan o‘tish (JWT tokenlar qaytaradi)
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.GithubAuthRequest true "GitHub authorization code"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth/github [post]
func (h Handler) GithubAuth(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}
	githubUser, err := h.services.Github().ExchangeCodeForUser(c.Request.Context(), req.Code)
	if err != nil {
		handleResponse(c, h.log, "GitHub login failed", http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := h.services.User().GithubAuth(c.Request.Context(), githubUser.Email, githubUser.Name, githubUser.GithubID)
	if err != nil {
		handleResponse(c, h.log, "failed to create/login user", http.StatusInternalServerError, err.Error())
		return
	}
	token, refresh, err := jwt.GenerateJWT(map[string]interface{}{
		"user_id":   userID,
		"user_role": "user",
	})
	resp := models.LoginResponse{
		ID:           userID,
		AccessToken:  token,
		RefreshToken: refresh,
		Role:         "user",
	}
	handleResponse(c, h.log, "login via github", http.StatusOK, resp)
}

// FacebookAuth godoc
// @Summary      Facebook orqali login yoki registratsiya
// @Description  Facebook OAuth code orqali login yoki ro‘yxatdan o‘tish (JWT tokenlar qaytaradi)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body models.FacebookAuthRequest true "Facebook authorization code"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} models.Response
// @Failure      401 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /auth/facebook [post]
func (h Handler) FacebookAuth(c *gin.Context) {
	var req models.FacebookAuthRequest // { "code": "FACEBOOK_RETURNED_CODE" }

	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "invalid request", http.StatusBadRequest, err.Error())
		return
	}

	// 1. Facebook’dan user infoni olish
	fbUser, err := h.services.Facebook().ExchangeCodeForUser(c.Request.Context(), req.Code)
	if err != nil {
		handleResponse(c, h.log, "Facebook login failed", http.StatusUnauthorized, err.Error())
		return
	}

	// 2. User bazada bormi, bo‘lmasa register qilamiz (GoogleAuth usuliga o‘xshab)
	userID, err := h.services.User().FacebookAuth(c.Request.Context(), fbUser.Email, fbUser.Name, fbUser.FacebookID)
	if err != nil {
		handleResponse(c, h.log, "failed to create/login user", http.StatusInternalServerError, err.Error())
		return
	}

	// 3. JWT tokenlar generatsiya qilinadi
	token, refresh, err := jwt.GenerateJWT(map[string]interface{}{
		"user_id":   userID,
		"user_role": "user",
	})
	if err != nil {
		handleResponse(c, h.log, "failed to generate token", http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.LoginResponse{
		ID:           userID,
		AccessToken:  token,
		RefreshToken: refresh,
		Role:         "user",
	}
	handleResponse(c, h.log, "login via facebook", http.StatusOK, resp)
}

// Logout godoc
// @Summary      Logout (chiqish)
// @Description  JWT tokenlarni va sessionni bekor qiladi
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body models.LogoutRequest false "Logout request (refresh_token optional)"
// @Success      200 {object} models.Response
// @Failure      401 {object} models.Response
// @Router       /logout [post]
// @Security     ApiKeyAuth
func (h Handler) Logout(c *gin.Context) {
	accessToken := ExtractBearerToken(c)
	var req models.LogoutRequest
	_ = c.ShouldBindJSON(&req)

	// Contextni uzating!
	ctx := c.Request.Context()

	if accessToken != "" {
		_ = h.services.Redis().BlacklistToken(ctx, accessToken)
	}
	if req.RefreshToken != "" {
		_ = h.services.Redis().BlacklistToken(ctx, req.RefreshToken)
	}

	// Cookie ni tozalash (agar front uchun kerak bo‘lsa)
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	handleResponse(c, h.log, "Logged out successfully", http.StatusOK, nil)
}

// Helper: Bearer tokenni olish
func ExtractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
