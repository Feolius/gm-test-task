package httpserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type currencyRateLoader interface {
	load(ctx context.Context, date, currency string) (float64, error)
}

type rateResponse struct {
	Rate float64 `json:"rate"`
}

type exchangeRateHandler struct {
	loader currencyRateLoader
}

func (h *exchangeRateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	query := request.URL.Query()
	date := query.Get("date")
	currency := query.Get("currency")

	var missedParams []string
	if date == "" {
		missedParams = append(missedParams, "date")
	}
	if currency == "" {
		missedParams = append(missedParams, "currency")
	}
	if len(missedParams) > 0 {
		writer.WriteHeader(http.StatusBadRequest)
		responseErr := json.NewEncoder(writer).Encode(errResponse{Message: "Missing query parameters: " +
			strings.Join(missedParams, ", ")})
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
	responseErr := json.NewEncoder(writer).Encode(rateResponse{Rate: res})
	if responseErr != nil {
		log.Printf("failed to write response: %v", err)
	}
}
