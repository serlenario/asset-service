package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"asset-service/internal/service"
)

func UploadAssetHandler(svc *service.AssetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"bad method"}`, http.StatusBadRequest)
			return
		}
		uid := r.Context().Value("uid").(int64)
		name := r.URL.Path[len("/api/assets/"):]
		data, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		if err := svc.Upload(r.Context(), uid, name, data); err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func DownloadAssetHandler(svc *service.AssetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"bad method"}`, http.StatusBadRequest)
			return
		}
		uid := r.Context().Value("uid").(int64)
		name := r.URL.Path[len("/api/assets/"):]
		data, err := svc.Download(r.Context(), uid, name)
		if err != nil {
			if exists, _ := svc.Exists(r.Context(), name); exists {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			} else {
				http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			}
			return
		}
		w.Write(data)
	}
}

func ListAssetsHandler(svc *service.AssetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"bad method"}`, http.StatusBadRequest)
			return
		}
		uid := r.Context().Value("uid").(int64)
		list, err := svc.List(r.Context(), uid)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	}
}

func DeleteAssetHandler(svc *service.AssetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, `{"error":"bad method"}`, http.StatusBadRequest)
			return
		}
		uid := r.Context().Value("uid").(int64)
		name := r.URL.Path[len("/api/assets/"):]
		if err := svc.Delete(r.Context(), uid, name); err != nil {
			if exists, _ := svc.Exists(r.Context(), name); exists {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			} else {
				http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
