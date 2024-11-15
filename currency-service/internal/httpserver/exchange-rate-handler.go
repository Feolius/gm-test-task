package httpserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type currencyRateLoader interface {
	load(ctx context.Context, date, currency string) (float64, error)
}

type exchangeRateHandler struct {
	loader currencyRateLoader
}

type errResponse struct {
	Message string `json:"message"`
}

type successResponse struct {
	Rate float64 `json:"rate"`
}

func (h *exchangeRateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	query := request.URL.Query()
	date := query.Get("date")
	currency := query.Get("currency")

	if date == "" {
		writer.WriteHeader(http.StatusBadRequest)
		responseErr := json.NewEncoder(writer).Encode(errResponse{Message: "Missing query parameter: date"})
		if responseErr != nil {
			log.Printf("failed to write response: %v", responseErr)
		}
		return
	}
	if currency == "" {
		writer.WriteHeader(http.StatusBadRequest)
		responseErr := json.NewEncoder(writer).Encode(errResponse{Message: "Missing query parameter: currency"})
		if responseErr != nil {
			log.Printf("failed to write response: %v", responseErr)
		}
		return
	}

	res, err := h.loader.load(request.Context(), date, currency)
	if err != nil {
		log.Printf("failed to load currency rate for date %s and currency %s: %v", date, currency, err)
		writer.WriteHeader(http.StatusInternalServerError)
		responseErr := json.NewEncoder(writer).Encode(errResponse{Message: "Internal Server Error"})
		if responseErr != nil {
			log.Printf("failed to write response: %v", err)
		}
		return
	}
	if res == 0 {
		writer.WriteHeader(http.StatusNotFound)
		responseErr := json.NewEncoder(writer).Encode(errResponse{Message: "Rate not found"})
		if responseErr != nil {
			log.Printf("failed to write response: %v", err)
		}
		return
	}
	responseErr := json.NewEncoder(writer).Encode(successResponse{Rate: res})
	if responseErr != nil {
		log.Printf("failed to write response: %v", err)
	}
}
