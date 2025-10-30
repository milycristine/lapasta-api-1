package auth

import (
	"crypto/sha256"
	"encoding/json"
	"lapasta/internal/models"
	"net/http"
)

func LoginHandler(s AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest models.Login
		response := models.ResponseDefaultModel{
			IsSuccess: true,
		}

		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			response.IsSuccess = false
			response.ErrorMessage = "Erro ao decodificar a requisição de login"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		hashed := sha256.Sum256([]byte(loginRequest.Password))

		_, err := s.Login(loginRequest.Username, hashed)
		if err != nil {
			response.IsSuccess = false
			response.ErrorMessage = "Usuário ou senha inválidas"
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		user, _ := s.GetUsuario(loginRequest.Username)
		response.IsSuccess = true
		response.Data = user
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
