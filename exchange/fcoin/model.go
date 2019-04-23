package fcoin

import "encoding/json"

//Get Struct by Exchange
//Convert Sample Json to Go Struct

type JsonResponse struct {
	Status  int             `json:"status"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"msg"`
}

type CoinsData []struct {
	Data string `json:"data"`
}

type PairsData []struct {
	Name          string `json:"name"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	PriceDecimal  int    `json:"price_decimal"`
	AmountDecimal int    `json:"amount_decimal"`
}
type OrderBook struct {
	Bids []float64 `json:"bids"`
	Asks []float64 `json:"asks"`
	Ts   int64     `json:"ts"`
	Seq  int       `json:"seq"`
	Type string    `json:"type"`
}

type fcoin struct {
	Status int      `json:"status"`
	Data   []string `json:"data"`
}

type AccountBalances []struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Balance   string `json:"balance"`
}

type TradeHistory []struct {
	ID            string `json:"id"`
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	Price         string `json:"price"`
	Amount        string `json:"amount"`
	State         string `json:"state"`
	ExecutedValue string `json:"executed_value"`
	FillFees      string `json:"fill_fees"`
	FilledAmount  string `json:"filled_amount"`
	CreatedAt     int    `json:"created_at"`
	Source        string `json:"source"`
}
