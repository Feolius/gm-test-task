package httpserver

import (
	"context"
	auth "gateway-service/internal/authenticator"
	"gateway-service/internal/config"
	"gateway-service/internal/userstorage"
	"net"
	"net/http"
)

func NewHttpServer(ctx context.Context, cfg *config.Config) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	mux.Handle("POST /login", loginHandler{cfg, auth.NewAuthenticator(userstorage.NewInMemoryStorage(), cfg)})
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	return httpServer
}
