package httpserver

import (
	"context"
	"gateway-service/internal/authclient"
	"gateway-service/internal/config"
	"gateway-service/internal/repository"
	"net"
	"net/http"
	"time"
)

const ReadHeaderTimeoutInSeconds = 2

func NewHttpServer(ctx context.Context, cfg *config.Config) *http.Server {
	userRepository := repository.NewInMemoryUserRepository()
	userSearcher := &inMemoryUserSearcher{userRepository}
	authClient := authclient.NewAuthClient(userSearcher, cfg)
	authHandler := authMiddleware(&authClientAuthenticator{authClient})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	mux.Handle("POST /login", loginHandler{cfg, &authTokenProvider{authClient}})
	mux.Handle("GET /exchange-rate", authHandler(&exchangeRateHandler{}))
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ReadHeaderTimeout: ReadHeaderTimeoutInSeconds * time.Second,
	}

	return httpServer
}
