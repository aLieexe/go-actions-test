package auth

import (
	"context"
	"errors"
	"fmt"
	"jwt-golang/internal/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtProvider struct {
	accessTokenSecret  []byte
	refreshTokenSecret []byte
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
}
type CustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewJwtProvider(accessTokenSecret string, refreshTokenSecret string, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *JwtProvider {
	return &JwtProvider{
		accessTokenSecret:  []byte(accessTokenSecret),
		accessTokenTTL:     accessTokenTTL * time.Second,
		refreshTokenSecret: []byte(refreshTokenSecret),
		refreshTokenTTL:    refreshTokenTTL * time.Second,
	}
}

func (p *JwtProvider) GenerateAccessToken(ctx context.Context, user *models.Users) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "jwt-golang-auth",
			Subject:   user.Id,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(p.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(p.accessTokenSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return tokenString, nil
}

func (p *JwtProvider) GenerateRefreshToken(ctx context.Context, user *models.Users) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "jwt-golang-auth",
			Subject:   user.Id,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(p.refreshTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(p.refreshTokenSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return tokenString, nil
}

func (p *JwtProvider) validateAccessToken(authorizationHeader string) error {
	if authorizationHeader == "" {
		return errors.New("authorization header is required")
	}

	tokenString, err := p.extractTokenFromHeader(authorizationHeader)
	if err != nil {
		return err
	}
	if tokenString == "" {
		return errors.New("invalid authorization header format")
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.accessTokenSecret, nil
	})

	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return errors.New("token is not valid")
	}

	return nil
}

func (p *JwtProvider) validateRefreshToken(refreshToken string) error {
	token, err := jwt.ParseWithClaims(refreshToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.refreshTokenSecret, nil
	})

	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return errors.New("token is not valid")
	}

	return nil
}

func (p *JwtProvider) extractTokenFromHeader(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("failed to parse")
	}
	return parts[1], nil
}

func (p *JwtProvider) extractClaims(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return p.accessTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

func (p *JwtProvider) extractUserIdFromRefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return p.refreshTokenSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", errors.New("invalid token claims")
}
