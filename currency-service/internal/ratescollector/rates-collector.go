package ratescollector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const apiurl = "https://latest.currency-api.pages.dev/v1/currencies/rub.json"

type CurrencyRate struct {
	Date     string
	Currency string
	Rate     float64
}

type CurrencyRateSaver interface {
	Save(ctx context.Context, cr *CurrencyRate) error
}

type RatesCollector struct {
	saver CurrencyRateSaver
}

type apiResponse struct {
	Date string             `json:"date"`
	Rub  map[string]float64 `json:"rub"`
}

func (rc *RatesCollector) Collect(ctx context.Context) error {
	log.Printf("rates collector job started at %v\n", time.Now())
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

	var insertErrors []error
	for currency, rate := range responseData.Rub {
		currencyRate := &CurrencyRate{
			Date:     responseData.Date,
			Currency: currency,
			Rate:     rate,
		}
		insertError := rc.saver.Save(ctx, currencyRate)
		if insertError != nil {
			insertErrors = append(insertErrors, insertError)
		}
	}
	if len(insertErrors) > 0 {
		return errors.Join(insertErrors...)
	}

	log.Printf("rates collector job finished at %v\n", time.Now())
	return nil
}

func NewRatesCollector(saver CurrencyRateSaver) *RatesCollector {
	return &RatesCollector{saver}
}
