package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenInvalid = errors.New("token is invalid or expired")
)

func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		panic("JWT_SECRET_KEY is not set in environment")
	}
	return []byte(secret)
}

func GenerateJWT(claims map[string]interface{}) (string, string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	aClaims := accessToken.Claims.(jwt.MapClaims)
	rClaims := refreshToken.Claims.(jwt.MapClaims)

	for k, v := range claims {
		aClaims[k] = v
		rClaims[k] = v
	}

	aClaims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	aClaims["iat"] = time.Now().Unix()

	rClaims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	rClaims["iat"] = time.Now().Unix()

	accessStr, err := accessToken.SignedString(getSecretKey())
	if err != nil {
		return "", "", err
	}

	refreshStr, err := refreshToken.SignedString(getSecretKey())
	if err != nil {
		return "", "", err
	}

	return accessStr, refreshStr, nil
}

func ParseToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getSecretKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result := make(map[string]interface{})
		for k, v := range claims {
			result[k] = v
		}
		return result, nil
	}

	return nil, ErrTokenInvalid
}

func ExtractClaims(tokenString string) (map[string]interface{}, error) {
	parsedClaims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if userID, ok := parsedClaims["user_id"]; ok {
		result["user_id"] = stringify(userID)
	}
	if role, ok := parsedClaims["user_role"]; ok {
		result["user_role"] = stringify(role)
	}
	return result, nil
}

func stringify(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
