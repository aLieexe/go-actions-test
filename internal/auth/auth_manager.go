package auth

import (
	"context"
	"fmt"
	"jwt-golang/internal/models"
)

type AuthManager struct {
	Provider AuthProvider
}

func NewAuthManager(provider AuthProvider) *AuthManager {
	return &AuthManager{
		Provider: provider,
	}
}

type AuthProvider interface {
	GenerateAccessToken(ctx context.Context, user *models.Users) (string, error)
	GenerateRefreshToken(ctx context.Context, user *models.Users) (string, error)
	validateAccessToken(token string) error
	validateRefreshToken(token string) error
	extractTokenFromHeader(authHeader string) (string, error)
	extractClaims(tokenString string) (*CustomClaims, error)
	extractUserIdFromRefreshToken(tokenString string) (string, error)
}

func (m *AuthManager) GenerateTokens(ctx context.Context, user *models.Users) ([]string, error) {
	accessToken, err := m.Provider.GenerateAccessToken(ctx, user)

	if err != nil {
		return nil, err
	}

	refreshToken, err := m.Provider.GenerateRefreshToken(ctx, user)

	if err != nil {
		return nil, err
	}

	tokens := make([]string, 0)
	tokens = append(tokens, accessToken)
	tokens = append(tokens, refreshToken)

	return tokens, nil
}

func (m *AuthManager) ValidateAccessToken(authorizationHeader string) error {
	return m.Provider.validateAccessToken(authorizationHeader)
}

func (m *AuthManager) ValidateRefreshToken(authorizationHeader string) error {
	return m.Provider.validateRefreshToken(authorizationHeader)
}

// validate refreshToken, generate than i think can call the generate token?
// idk
func (m *AuthManager) GetUserIdFromRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	err := m.Provider.validateRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh token is invalid")
	}

	userId, err := m.Provider.extractUserIdFromRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (m *AuthManager) ExtractTokenFromHeader(authorizationHeader string) (string, error) {
	return m.Provider.extractTokenFromHeader(authorizationHeader)
}

func (m *AuthManager) ExtractClaims(tokenString string) (*CustomClaims, error) {
	return m.Provider.extractClaims(tokenString)
}
