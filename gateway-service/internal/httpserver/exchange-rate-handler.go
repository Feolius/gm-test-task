package httpserver

import (
	"log"
	"net/http"
)

type requestProxy interface {
	handle(writer http.ResponseWriter, request *http.Request) error
}

type proxyHandler struct {
	proxy requestProxy
}

func (h *proxyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := h.proxy.handle(writer, request)
	if err != nil {
		log.Printf("failed to proxy request: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
