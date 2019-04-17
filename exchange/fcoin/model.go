package fcoin

import "encoding/json"

//Get Struct by Exchange
//Convert Sample Json to Go Struct

type JsonResponse struct {
	Success bool            `json:"Success"`
	Message interface{}     `json:"Message"`
	Data    json.RawMessage `json:"data"`
	Error   interface{}     `json:"Error"`
}

type CoinsData struct {
	Status int      `json:"status"`
	Data   []string `json:"data"`
}

type PairsData struct {
	Status int `json:"status"`
	Data   []struct {
		Name          string `json:"name"`
		BaseCurrency  string `json:"base_currency"`
		QuoteCurrency string `json:"quote_currency"`
		PriceDecimal  int    `json:"price_decimal"`
		AmountDecimal int    `json:"amount_decimal"`
	} `json:"data"`
}

type OrderBook struct {
	Status int `json:"status"`
	Data   struct {
		Bids []float64 `json:"bids"`
		Asks []float64 `json:"asks"`
		Ts   int64     `json:"ts"`
		Seq  int       `json:"seq"`
		Type string    `json:"type"`
	} `json:"data"`
}

type fcoin struct {
	Status int      `json:"status"`
	Data   []string `json:"data"`
}

type AccountBalances struct {
	Status int `json:"status"`
	Data   []struct {
		Currency  string `json:"currency"`
		Available string `json:"available"`
		Frozen    string `json:"frozen"`
		Balance   string `json:"balance"`
	} `json:"data"`
}

type TradeHistory struct {
	Status int `json:"status"`
	Data   []struct {
		Price        string `json:"price"`
		FillFees     string `json:"fill_fees"`
		FilledAmount string `json:"filled_amount"`
		Side         string `json:"side"`
		Type         string `json:"type"`
		CreatedAt    int    `json:"created_at"`
	} `json:"data"`
}
