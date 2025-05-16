package service

import (
	"context"
	"errors"

	"asset-service/internal/repo"
)

type AssetService struct{ assets repo.AssetRepo }

func NewAssetService(a repo.AssetRepo) *AssetService { return &AssetService{assets: a} }

func (s *AssetService) Upload(ctx context.Context, uid int64, name string, data []byte) error {
	if name == "" {
		return errors.New("asset name required")
	}
	return s.assets.Upsert(ctx, uid, name, data)
}

func (s *AssetService) Download(ctx context.Context, uid int64, name string) ([]byte, error) {
	return s.assets.Get(ctx, uid, name)
}

func (s *AssetService) List(ctx context.Context, uid int64) ([]string, error) {
	return s.assets.List(ctx, uid)
}

func (s *AssetService) Delete(ctx context.Context, uid int64, name string) error {
	return s.assets.Delete(ctx, uid, name)
}

func (s *AssetService) Exists(ctx context.Context, name string) (bool, error) {
	return s.assets.Exists(ctx, name)
}
