package handlers

import (
	"encoding/json"
	"myproject/internal/config"
	"myproject/internal/database"
	"myproject/internal/models"
	"myproject/internal/redis"
	"net/http"
)

func CurrencyTransferHandle(w http.ResponseWriter, r *http.Request) {
	var rdb redis.RedisDataBase
	cfg := config.LoadConfig()
	db := database.DataBasePostgres{URL: cfg.DatabaseURL}
	rdb.Init()

	var currRequest = models.CurrencyRequest{}
	if err := json.NewDecoder(r.Body).Decode(&currRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch {
	case currRequest.Amount == 0:
		err := models.ErrorAmountRequired
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case currRequest.Amount < 0:
		err := models.ErrorAmountInvalid
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case currRequest.FromCurrency == "":
		err := models.ErrorInitialCurrencyRequired
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case currRequest.ToCurrency == "":
		err := models.ErrorFinalCurrencyRequired
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rateTo, err := rdb.GateRate(currRequest.ToCurrency)
	if err == models.ErrorCurrencyNotFound {
		curr, err := db.GetCurrency(currRequest.ToCurrency)
		if err == models.ErrorCurrencyNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err != nil {
			http.Error(w, "Failed to get exchange rate", http.StatusInternalServerError)
			return
		}
		rdb.SetCurrency(curr)
		rateTo = curr.Rate
	} else if err != nil {
		http.Error(w, "Failed to get exchange rate", http.StatusInternalServerError)
		return
	}

	rateFrom, err := rdb.GateRate(currRequest.FromCurrency)
	if err == models.ErrorCurrencyNotFound {
		curr, err := db.GetCurrency(currRequest.FromCurrency)
		if err == models.ErrorCurrencyNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err != nil {
			http.Error(w, "Failed to get exchange rate", http.StatusInternalServerError)
			return
		}
		rdb.SetCurrency(curr)
		rateFrom = curr.Rate
	} else if err != nil {
		http.Error(w, "Failed to get exchange rate", http.StatusInternalServerError)
		return
	}

	rate := rateTo / rateFrom
	convertedAmount := currRequest.Amount * rate

	response := models.CurrencyResponse{
		Amount:          currRequest.Amount,
		ConvertedAmount: convertedAmount,
		Rate:            rate,
		FromCurrency:    currRequest.FromCurrency,
		ToCurrency:      currRequest.ToCurrency,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
