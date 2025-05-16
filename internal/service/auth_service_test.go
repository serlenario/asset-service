package service_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"asset-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type userRepoMock struct{ mock.Mock }

func (m *userRepoMock) GetByLogin(ctx context.Context, login string) (int64, string, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(int64), args.String(1), args.Error(2)
}

type sessionRepoMock struct{ mock.Mock }

func (m *sessionRepoMock) DeleteByUser(ctx context.Context, uid int64) error {
	return m.Called(ctx, uid).Error(0)
}
func (m *sessionRepoMock) Create(ctx context.Context, uid int64, ip net.IP, expiresAt time.Time) (string, error) {
	args := m.Called(ctx, uid, ip, expiresAt)
	return args.String(0), args.Error(1)
}
func (m *sessionRepoMock) Validate(ctx context.Context, token string) (int64, time.Time, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(int64), args.Get(1).(time.Time), args.Error(2)
}

func TestAuthenticate_Success(t *testing.T) {
	prev := service.ComparePassword
	defer func() { service.ComparePassword = prev }()
	service.ComparePassword = func(hashedPassword, password []byte) error {
		return nil
	}

	ctx := context.Background()
	uMock := new(userRepoMock)
	sMock := new(sessionRepoMock)
	uMock.On("GetByLogin", ctx, "alice").
		Return(int64(1), "$2a$10$someirrelevanthash", nil)
	sMock.On("DeleteByUser", ctx, int64(1)).Return(nil)
	sMock.On("Create", ctx, int64(1), mock.AnythingOfType("net.IP"), mock.AnythingOfType("time.Time")).
		Return("session123", nil)

	svc := service.NewAuthService(uMock, sMock)
	token, err := svc.Authenticate(ctx, "alice", "password", "127.0.0.1")

	assert.NoError(t, err)
	assert.Equal(t, "session123", token)
}

func TestAuthenticate_Fail(t *testing.T) {
	ctx := context.Background()
	uMock := new(userRepoMock)
	sMock := new(sessionRepoMock)
	uMock.On("GetByLogin", ctx, "bob").Return(int64(0), "", errors.New("not found"))

	svc := service.NewAuthService(uMock, sMock)
	_, err := svc.Authenticate(ctx, "bob", "pass", "127.0.0.1")

	assert.Error(t, err)
}

func TestValidate_Success(t *testing.T) {
	ctx := context.Background()
	sMock := new(sessionRepoMock)
	exp := time.Now().Add(1 * time.Hour)
	sMock.On("Validate", ctx, "token123").Return(int64(2), exp, nil)

	svc := service.NewAuthService(nil, sMock)
	uid, err := svc.Validate(ctx, "token123")

	assert.NoError(t, err)
	assert.Equal(t, int64(2), uid)
}

func TestValidate_Expired(t *testing.T) {
	ctx := context.Background()
	sMock := new(sessionRepoMock)
	past := time.Now().Add(-1 * time.Hour)
	sMock.On("Validate", ctx, "tokenOld").Return(int64(3), past, nil)

	svc := service.NewAuthService(nil, sMock)
	_, err := svc.Validate(ctx, "tokenOld")

	assert.Error(t, err)
}
