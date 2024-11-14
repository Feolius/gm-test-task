package httpserver

import (
	"database/sql"
	"net/http"
)

type exchangeRateHandler struct {
	db *sql.DB
}

func (e *exchangeRateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// TODO implement me
	panic("implement me")
}
