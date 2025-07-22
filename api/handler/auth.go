package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"test/pkg/jwt"
)

// AuthorizerMiddleware tekshiradi JWT tokenni va contextga user_id va user_role qo‘shadi
func (h Handler) AuthorizerMiddleware(c *gin.Context) {
	// 1. Headerdan "Authorization"ni olish
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// 2. "Bearer TOKEN" formatini tekshirish
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format. Use 'Bearer {token}'"})
		return
	}

	tokenStr := parts[1]

	// 3. Tokenni JWT orqali parse qilish
	claims, err := jwt.ExtractClaims(tokenStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// 4. user_id mavjudligini tekshirish
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	// 5. Contextga user_id qo‘shamiz
	c.Set("user_id", userID)

	// 6. Agar user_role mavjud bo‘lsa, uni ham qo‘shamiz (ixtiyoriy)
	if role, ok := claims["user_role"].(string); ok {
		c.Set("user_role", role)
	}

	// 7. Davom ettiramiz
	c.Next()
}
