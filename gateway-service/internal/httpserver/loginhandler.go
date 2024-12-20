package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"gateway-service/internal/config"
	"log"
	"net/http"
)

type tokenProvider interface {
	getToken(ctx context.Context, username, password string) (string, error)
}

type loginHandler struct {
	cfg *config.Config
	tokenProvider
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (l *loginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var req loginRequest
	err := decoder.Decode(&req)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")

	token, err := l.tokenProvider.getToken(request.Context(), req.Username, req.Password)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		var invalidCredentialsErr *invalidCredentialsError
		if errors.As(err, &invalidCredentialsErr) {
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(struct {
				Error string `json:"error"`
			}{Error: "invalid credentials"})
			return
		}
		log.Printf("authentication error: %s", err.Error())
		return
	}

	res := loginResponse{
		Token: token,
	}
	err = json.NewEncoder(writer).Encode(res)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("login response encoding error: %s", err.Error())
	}
}
