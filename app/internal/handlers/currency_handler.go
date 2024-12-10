package handlers

import (
	"encoding/json"
	"myproject/internal/models"
	"myproject/internal/services"
	"net/http"
)

func CheckRequest(w http.ResponseWriter, currRequest models.CurrencyRequest) error {
	switch {
	default:
		return nil
	case currRequest.Amount == 0:
		err := models.ErrorAmountRequired
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	case currRequest.Amount < 0:
		err := models.ErrorAmountInvalid
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	case currRequest.FromCurrency == "":
		err := models.ErrorInitialCurrencyRequired
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	case currRequest.ToCurrency == "":
		err := models.ErrorFinalCurrencyRequired
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
}
func CurrencyTransferHandle(w http.ResponseWriter, r *http.Request) {
	var currRequest = models.CurrencyRequest{}
	if err := json.NewDecoder(r.Body).Decode(&currRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := CheckRequest(w, currRequest); err != nil {
		return
	}

	rateTo, err := services.GateRate(currRequest.ToCurrency)
	if err == models.ErrorCurrencyNotFound {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if err != nil {
		http.Error(w, "Failed to get exchange rate", http.StatusInternalServerError)
		return
	}

	rateFrom, err := services.GateRate(currRequest.FromCurrency)
	if err == models.ErrorCurrencyNotFound {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
