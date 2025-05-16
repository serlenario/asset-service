package service

import (
	"context"
	"errors"
	"net"
	"time"

	"asset-service/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

var ComparePassword = bcrypt.CompareHashAndPassword

var ErrInvalidCredentials = errors.New("invalid login/password")

type AuthService struct {
	users    repo.UserRepo
	sessions repo.SessionRepo
}

func NewAuthService(u repo.UserRepo, s repo.SessionRepo) *AuthService {
	return &AuthService{users: u, sessions: s}
}

func (s *AuthService) Authenticate(ctx context.Context, login, password, ipStr string) (string, error) {
	uid, hash, err := s.users.GetByLogin(ctx, login)
	if err != nil || ComparePassword([]byte(hash), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}
	_ = s.sessions.DeleteByUser(ctx, uid)
	exp := time.Now().Add(24 * time.Hour)
	return s.sessions.Create(ctx, uid, net.ParseIP(ipStr), exp)
}

func (s *AuthService) Validate(ctx context.Context, token string) (int64, error) {
	uid, exp, err := s.sessions.Validate(ctx, token)
	if err != nil || time.Now().After(exp) {
		return 0, errors.New("unauthorized")
	}
	return uid, nil
}
