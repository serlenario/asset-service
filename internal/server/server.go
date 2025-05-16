package server

import (
	"net/http"

	"asset-service/internal/config"
	"asset-service/internal/handlers"
	"asset-service/internal/middleware"
	"asset-service/internal/repo"
	"asset-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(pool *pgxpool.Pool, cfg *config.Config) *http.Server {
	userRepo := repo.NewPG(pool)
	sessRepo := repo.NewPG(pool)
	assetRepo := repo.NewPG(pool)

	authSvc := service.NewAuthService(userRepo, sessRepo)
	assetSvc := service.NewAssetService(assetRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth", handlers.AuthHandler(authSvc))

	authMW := middleware.Auth(authSvc)
	mux.Handle("/api/assets/", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.UploadAssetHandler(assetSvc).ServeHTTP(w, r)
		case http.MethodGet:
			handlers.DownloadAssetHandler(assetSvc).ServeHTTP(w, r)
		case http.MethodDelete:
			handlers.DeleteAssetHandler(assetSvc).ServeHTTP(w, r)
		default:
			http.Error(w, `{"error":"bad method"}`, http.StatusBadRequest)
		}
	})))
	mux.Handle("/api/assets", authMW(http.HandlerFunc(handlers.ListAssetsHandler(assetSvc))))

	return &http.Server{Addr: cfg.Server.Address, Handler: mux}
}
