package services

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetExchangeRate() (map[string]float64, error) {
	resp, err := http.Get("https://v6.exchangerate-api.com/v6/453ea34d02c88e15836b7835/latest/USD")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	rowRates := result["conversion_rates"].(map[string]interface{})
	rates := make(map[string]float64)

	for currencyCode, value := range rowRates {
		switch v := value.(type) {
		case float64:
			rates[currencyCode] = v

		default:
			log.Printf("unsupported value type for currency %s: %T", currencyCode, value)
		}
	}
	return rates, err
}
