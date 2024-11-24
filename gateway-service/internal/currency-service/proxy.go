package currency_service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const ExchangeRatePath = "/exchange-rate"
const ExchangeHistoryPath = "/exchange-history"

type CurrencyServiceProxy struct {
	url *url.URL
}

func (p *CurrencyServiceProxy) ExchangeHistory(writer http.ResponseWriter, request *http.Request) error {
	return p.proxyToPath(ExchangeHistoryPath, writer, request)
}

func (p *CurrencyServiceProxy) ExchangeRate(writer http.ResponseWriter, request *http.Request) error {
	return p.proxyToPath(ExchangeRatePath, writer, request)
}

func (p *CurrencyServiceProxy) proxyToPath(path string, writer http.ResponseWriter, request *http.Request) error {
	reqUrl := *p.url
	reqUrl.Path += path
	reqUrl.RawQuery = request.URL.RawQuery

	req, err := http.NewRequest(request.Method, reqUrl.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create exhange-rate request: %w", err)
	}
	req = req.WithContext(request.Context())
	req.Body = request.Body
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute exhange-rate request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read exhange-rate response: %w", err)
	}

	copyHeader(writer.Header(), res.Header)
	writer.WriteHeader(res.StatusCode)
	_, err = io.Copy(writer, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to write exhange-rate service response: %w", err)
	}
	return nil
}

func NewCurrencyProxy(url *url.URL) *CurrencyServiceProxy {
	return &CurrencyServiceProxy{url}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
