package exchange

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Rate struct {
	log    *logrus.Entry
	xrHost string
	apiKey string
}

type Resp struct {
	Success bool `json:"success"`
	Query   struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	} `json:"query"`
	Info struct {
		Timestamp int     `json:"timestamp"`
		Rate      float64 `json:"rate"`
	} `json:"info"`
	Historical string  `json:"historical"`
	Date       string  `json:"date"`
	Result     float64 `json:"result"`
}

func NewExchangeRate(log *logrus.Entry, xrHost string, apiKey string) *Rate {
	return &Rate{
		log:    log.WithField("transport", "exchange"),
		xrHost: xrHost,
		apiKey: apiKey,
	}
}

func (e *Rate) Conversion(currency string, amount float64) (float64, error) {
	amountStr := fmt.Sprintf("%v", amount)

	url := e.xrHost + e.apiKey + currency + "&from=usd&amount=" + amountStr

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return 0, err
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("exchange api internal server error: %w", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return 0, fmt.Errorf("currency not found")
	case http.StatusForbidden:
		return 0, fmt.Errorf("invalid amount")
	default:
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return 0, fmt.Errorf("unexpected error")
		}
		return 0, fmt.Errorf("unexpected status code", res.StatusCode, string(body))
	}
	var result Resp

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("err decoding response: %w", err)
	}

	return result.Result, nil

}
