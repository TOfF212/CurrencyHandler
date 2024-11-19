package main


import (
    "encoding/json"
    "fmt"
    "net/http"
		"errors"
    
)


type CurrencyRequest  struct {
	Amount float64  `json:"amount"`
  InitialCurrency string `json:"from"`
  FinalCurrency string `json:"to"`
}


var(
	errAmountRequired = errors.New("Amount of currency is required ")
	errAmountInvalid = errors.New("Amount of currency is invalid ")
	errInitialCurrencyRequired = errors.New("Initial currency is required ")
	errFinalCurrencyRequired = errors.New("final currency is required ")
)


func CurrencyTransferHandle(w http.ResponseWriter, r *http.Request) {
    var currRequest = CurrencyRequest{}
    if err := json.NewDecoder(r.Body).Decode(&currRequest); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
		switch {
		case currRequest.Amount == 0:
			err := errAmountRequired
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case currRequest.Amount < 0:
			err := errAmountInvalid
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case currRequest.InitialCurrency == "":
			err := errInitialCurrencyRequired
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case currRequest.FinalCurrency == "":
			err := errFinalCurrencyRequired
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rate, err := getExchangeRate(currRequest.InitialCurrency, currRequest.FinalCurrency)
    if err != nil {
        http.Error(w, "Failed to get exchange rate", http.StatusInternalServerError)
        return
    }
		convertedAmount := currRequest.Amount * rate
    response := map[string]interface{}{
        "converted_amount": convertedAmount,
        "from_currency":    currRequest.InitialCurrency,
        "to_currency":      currRequest.FinalCurrency,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func getExchangeRate(from, to string) (float64, error) {
	resp, err := http.Get("https://v6.exchangerate-api.com/v6/453ea34d02c88e15836b7835/latest/" + from)
	if err != nil {
			return 0, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return 0, err
	}	

	rates := result["conversion_rates"].(map[string]interface{})
	rate, ok := rates[to]
	if !ok {
			return 0, fmt.Errorf("exchange rate for %s to %s not found", from, to)
	}

	return rate.(float64), nil
}


func main() {

    http.HandleFunc("/convert", CurrencyTransferHandle)

    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Printf("Server error: %s\n", err.Error())
    }
}
