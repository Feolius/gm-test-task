package httpserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type exchangeHistoryItem struct {
	Date string  `json:"date"`
	Rate float64 `json:"rate"`
}

type exchangeHistoryLoader interface {
	getHistory(ctx context.Context, currency, startDate, endDate string) ([]exchangeHistoryItem, error)
}

type exchangeHistoryHandler struct {
	loader exchangeHistoryLoader
}

func (h *exchangeHistoryHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	query := request.URL.Query()
	currency := query.Get("currency")
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")
	encoder := json.NewEncoder(writer)

	var missedParams []string
	if currency == "" {
		missedParams = append(missedParams, "currency")
	}
	if startDate == "" {
		missedParams = append(missedParams, "start_date")
	}
	if endDate == "" {
		missedParams = append(missedParams, "end_date")
	}
	if len(missedParams) > 0 {
		writer.WriteHeader(http.StatusBadRequest)
		responseErr := encoder.Encode(errResponse{Message: "Missed query parameters: " +
			strings.Join(missedParams, ", ")})
		if responseErr != nil {
			log.Printf("failed to write response: %v", responseErr)
		}
		return
	}

	historyItems, err := h.loader.getHistory(context.Background(), currency, startDate, endDate)
	if err != nil {
		log.Printf("failed to load exhange history for currency %s, start date %s and end date %s: %v",
			currency, startDate, endDate, err)
		writer.WriteHeader(http.StatusInternalServerError)
		responseErr := encoder.Encode(errResponse{Message: "Internal Server Error"})
		if responseErr != nil {
			log.Printf("failed to write response: %v", err)
		}
		return
	}
	err = encoder.Encode(historyItems)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to write response: %v", err)
	}
}
