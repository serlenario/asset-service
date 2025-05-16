package handlers

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	"asset-service/internal/service"
)

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func AuthHandler(svc *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"bad method"}`, http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		var req authRequest
		if json.Unmarshal(body, &req) != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		token, err := svc.Authenticate(r.Context(), req.Login, req.Password, ip)
		if err != nil {
			http.Error(w, `{"error":"invalid login/password"}`, http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
