package auth

import (
	"context"
	"labflow/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type TokenUtil struct {
	SecretKey string
	Redis     *redis.Client
}

func NewJwtUtils(secretKey string, redisClient *redis.Client) *TokenUtil {
	return &TokenUtil{
		SecretKey: secretKey,
		Redis:     redisClient,
	}
}

func (t TokenUtil) GenerateJWT(ctx context.Context, auth *models.Auth) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     auth.ID,
		"expire": time.Now().Add(time.Hour * 24 * 30).UnixMilli(),
	})

	tokenString, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	_, err = t.Redis.Set(ctx, tokenString, auth.ID, time.Hour*24*30).Result()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t TokenUtil) ParseJWT(ctx context.Context, tokenString string) (*models.Auth, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	expire := claims["expire"].(float64)
	if int64(expire) < time.Now().UnixMilli() {
		return nil, jwt.ErrSignatureInvalid
	}

	result, err := t.Redis.Exists(ctx, tokenString).Result()
	if err != nil {
		return nil, err
	}

	if result == 0 {
		return nil, jwt.ErrTokenNotValidYet
	}

	id := claims["id"].(float64)
	auth := &models.Auth{
		ID: int(id),
	}

	return auth, nil
}

func (t TokenUtil) InvalidateJWT(ctx context.Context, tokenString string) error {
	_, err := t.Redis.Del(ctx, tokenString).Result()
	return err
}

func (t TokenUtil) GenerateResetToken(ctx context.Context, email string) (string, error) {
	claims := jwt.MapClaims{
		"email":   email,
		"purpose": "reset_password",
		"expire":  time.Now().Add(time.Minute * 5).UnixMilli(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.SecretKey))
}

func (t TokenUtil) ParseResetToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	purpose, _ := claims["purpose"].(string)
	if purpose != "reset_password" {
		return "", jwt.ErrSignatureInvalid
	}

	expireFloat, ok := claims["expire"].(float64)
	if !ok || int64(expireFloat) < time.Now().UnixMilli() {
		return "", jwt.ErrSignatureInvalid
	}

	email, _ := claims["email"].(string)
	return email, nil
}
