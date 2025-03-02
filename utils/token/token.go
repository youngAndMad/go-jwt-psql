package token

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userId uint) (string, error) {
	tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(getJwtSecret()))
}

func getJwtSecret() string {
	return os.Getenv("JWT_SECRET")
}

func ValidateToken(c *gin.Context) error {
	tokenString := ExtractTokenFromHttpRequest(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getJwtSecret()), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractTokenFromHttpRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func parseJwt(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getJwtSecret()), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func ExtractUserIdFromJwt(c *gin.Context) (uint, error) {
	tokenString := ExtractTokenFromHttpRequest(c)

	tokenClaims, err := parseJwt(tokenString)

	if err != nil {
		return 0, err
	}

	uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", tokenClaims["userId"]), 10, 32)

	if err != nil {
		return 0, err
	}

	return uint(uid), nil
}
