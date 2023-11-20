package models

import "time"

const RatePerApiCall float32 = 0.1
const Currency = "INR"

type Bill struct {
	ID        string    `json:"id"`
	APIUsage  int       `json:"api_usage"`
	BillValue float32   `json:"bill_value"`
	Sattled   bool      `json:"sattled"`
	UserID    string    `json:"user_id"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func (bill *Bill) CalculateBillValue(usage int) {
	bill.APIUsage = usage
	bill.BillValue = RatePerApiCall * float32(usage)
	if bill.Currency == "" {
		bill.Currency = Currency
	}
}
