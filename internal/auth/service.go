package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

const sessionDuration = 24 * time.Hour

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, name, email, password string) (User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}
	return s.repo.CreateUser(ctx, name, email, hash)
}

func (s *Service) Login(ctx context.Context, email, password string) (Session, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == ErrUserNotFound {
			return Session{}, ErrInvalidCredentials
		}
		return Session{}, err
	}

	if !CheckPassword(password, user.PasswordHash) {
		return Session{}, ErrInvalidCredentials
	}

	token, err := generateToken()
	if err != nil {
		return Session{}, err
	}

	expiresAt := time.Now().Add(sessionDuration)
	if err := s.repo.CreateSession(ctx, token, user.ID, expiresAt); err != nil {
		return Session{}, err
	}

	return Session{Token: token, UserID: user.ID, ExpiresAt: expiresAt}, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	return s.repo.DeleteSession(ctx, token)
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
