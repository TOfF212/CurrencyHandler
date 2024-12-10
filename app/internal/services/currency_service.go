package services

import (
	"encoding/json"
	"log"
	"myproject/internal/config"
	"myproject/internal/database"
	"myproject/internal/models"
	"myproject/internal/redis"
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

func GateRate(currCode string) (float64, error) {
	var rdb redis.RedisDataBase
	cfg := config.LoadConfig()
	db := database.DataBasePostgres{URL: cfg.DatabaseURL}
	rdb.Init()

	rate, err := rdb.GateRate(currCode)
	if err == models.ErrorCurrencyNotFound {
		curr, err := db.GetCurrency(currCode)
		if err != nil {
			return 0, err
		}
		rdb.SetCurrency(curr)
		rate = curr.Rate
	} else if err != nil {
		return 0, err
	}
	return rate, nil
}
