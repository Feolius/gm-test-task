package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	auth "gateway-service/internal/authenticator"
	"gateway-service/internal/config"
	"net/http"
)

type authenticator interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type loginHandler struct {
	cfg           *config.Config
	authenticator authenticator
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (l loginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var req loginRequest
	err := decoder.Decode(&req)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := l.authenticator.Login(request.Context(), req.Username, req.Password)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		if errors.Is(err, auth.InvalidCredentialsError) {
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(struct {
				Error string `json:"error"`
			}{Error: "invalid credentials"})
			return
		}
		fmt.Printf("authentication error: %s", err.Error())
		return
	}
	res := loginResponse{
		Token: token,
	}
	err = json.NewEncoder(writer).Encode(res)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("login response encoding error: %s", err.Error())
	}
}
