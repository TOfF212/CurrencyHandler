package models

import (
	"errors"
)

var (
	ErrorCurrencyNotFound        = errors.New("currency is not found ")
	ErrorAmountRequired          = errors.New("amount of currency is required ")
	ErrorAmountInvalid           = errors.New("amount of currency is invalid ")
	ErrorInitialCurrencyRequired = errors.New("initial currency is required ")
	ErrorFinalCurrencyRequired   = errors.New("final currency is required ")
)

type Currency struct {
	ID       uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Currency string  `gorm:"type:varchar(3);not null" json:"currency"`
	Rate     float64 `gorm:"type:decimal(10,4);not null" json:"rate"`
}

type CurrencyRequest struct {
	Amount       float64 `json:"amount"`
	FromCurrency string  `json:"from"`
	ToCurrency   string  `json:"to"`
}

type CurrencyResponse struct {
	Amount          float64 `json:"amount_from"`
	ConvertedAmount float64 `json:"amount_to"`
	Rate            float64 `json:"rate"`
	FromCurrency    string  `json:"from"`
	ToCurrency      string  `json:"to"`
}
