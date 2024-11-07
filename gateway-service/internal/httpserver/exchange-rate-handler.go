package httpserver

import (
	"net/http"
)

type exchangeRateHandler struct {
}

func (h *exchangeRateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
