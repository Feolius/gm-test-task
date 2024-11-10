package ratescollector

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiurl = "https://latest.currency-api.pages.dev/v1/currencies/rub.json"

type RatesCollector struct {
	db *sql.DB
}

type apiResponse struct {
	Date string                 `json:"date"`
	Rub  map[string]interface{} `json:"rub"`
}

func (rc *RatesCollector) Collect(ctx context.Context) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, apiurl, nil)
	if err != nil {
		return fmt.Errorf("could not create api request to collect rates: %w", err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("could not fetch rates: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("could not fetch rates: response status %s, response msg: %s", response.Status, body)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("could not read rates response: %w", err)
	}
	var responseData apiResponse
	if err = json.Unmarshal(body, &responseData); err != nil {
		return fmt.Errorf("could not parse rates response: %w", err)
	}
	return nil
}

func NewRatesCollector(db *sql.DB) *RatesCollector {
	return &RatesCollector{db: db}
}
