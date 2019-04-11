package cryptopia

import "encoding/json"

//Get Struct by Exchange
//Convert Sample Json to Go Struct

type JsonResponse struct {
	Success bool            `json:"Success"`
	Message interface{}     `json:"Message"`
	Data    json.RawMessage `json:"data"`
	Error   interface{}     `json:"Error"`
}

type PairsData []struct {
	ID               int         `json:"Id"`
	Label            string      `json:"Label"`
	Currency         string      `json:"Currency"`
	Symbol           string      `json:"Symbol"`
	BaseCurrency     string      `json:"BaseCurrency"`
	BaseSymbol       string      `json:"BaseSymbol"`
	Status           string      `json:"Status"`
	StatusMessage    interface{} `json:"StatusMessage"`
	TradeFee         float64     `json:"TradeFee"`
	MinimumTrade     float64     `json:"MinimumTrade"`
	MaximumTrade     float64     `json:"MaximumTrade"`
	MinimumBaseTrade float64     `json:"MinimumBaseTrade"`
	MaximumBaseTrade float64     `json:"MaximumBaseTrade"`
	MinimumPrice     float64     `json:"MinimumPrice"`
	MaximumPrice     float64     `json:"MaximumPrice"`
}

type CoinsData []struct {
	ID                   int         `json:"Id"`
	Name                 string      `json:"Name"`
	Symbol               string      `json:"Symbol"`
	Algorithm            string      `json:"Algorithm"`
	WithdrawFee          float64     `json:"WithdrawFee"`
	MinWithdraw          float64     `json:"MinWithdraw"`
	MaxWithdraw          float64     `json:"MaxWithdraw"`
	MinBaseTrade         float64     `json:"MinBaseTrade"`
	IsTipEnabled         bool        `json:"IsTipEnabled"`
	MinTip               float64     `json:"MinTip"`
	DepositConfirmations int         `json:"DepositConfirmations"`
	Status               string      `json:"Status"`
	StatusMessage        interface{} `json:"StatusMessage"`
	ListingStatus        string      `json:"ListingStatus"`
}

type OrderBook struct {
	Buy []struct {
		TradePairID int     `json:"TradePairId"`
		Label       string  `json:"Label"`
		Price       float64 `json:"Price"`
		Volume      float64 `json:"Volume"`
		Total       float64 `json:"Total"`
	} `json:"Buy"`
	Sell []struct {
		TradePairID int     `json:"TradePairId"`
		Label       string  `json:"Label"`
		Price       float64 `json:"Price"`
		Volume      float64 `json:"Volume"`
		Total       float64 `json:"Total"`
	} `json:"Sell"`
}

type AccountBalances []struct {
	CurrencyID      int         `json:"CurrencyId"`
	Symbol          string      `json:"Symbol"`
	Total           float64     `json:"Total"`
	Available       float64     `json:"Available"`
	Unconfirmed     float64     `json:"Unconfirmed"`
	HeldForTrades   float64     `json:"HeldForTrades"`
	PendingWithdraw float64     `json:"PendingWithdraw"`
	Address         interface{} `json:"Address"`
	Status          string      `json:"Status"`
	StatusMessage   interface{} `json:"StatusMessage"`
	BaseAddress     string      `json:"BaseAddress"`
}

type PlaceOrder struct {
	OrderID      int64   `json:"OrderId"`
	FilledOrders []int64 `json:"FilledOrders"`
}

type CryptopiaResponse struct {
	Success bool        `json:"Success"`
	Error   interface{} `json:"Error"`
	Data    []int       `json:"Data"`
}

type TradeHistory []struct {
	OrderID     int     `json:"OrderId"`
	TradePairID int     `json:"TradePairId"`
	Market      string  `json:"Market"`
	Type        string  `json:"Type"`
	Rate        float64 `json:"Rate"`
	Amount      float64 `json:"Amount"`
	Total       float64 `json:"Total"`
	Remaining   float64 `json:"Remaining"`
	TimeStamp   string  `json:"TimeStamp"`
}
