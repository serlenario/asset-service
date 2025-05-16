package service_test

import (
	"context"
	
	"github.com/stretchr/testify/mock"
)

type assetRepoMock struct{ mock.Mock }

func (m *assetRepoMock) Upsert(ctx context.Context, uid int64, name string, data []byte) error {
	return m.Called(ctx, uid, name, data).Error(0)
}
func (m *assetRepoMock) Get(ctx context.Context, uid int64, name string) ([]byte, error) {
	args := m.Called(ctx, uid, name)
	return args.Get(0).([]byte), args.Error(1)
}
func (m *assetRepoMock) List(ctx context.Context, uid int64) ([]string, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]string), args.Error(1)
}
func (m *assetRepoMock) Delete(ctx context.Context, uid int64, name string) error {
	return m.Called(ctx, uid, name).Error(0)
}
func (m *assetRepoMock) Exists(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}
