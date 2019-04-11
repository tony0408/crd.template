package blank

//Get Struct by Exchange
//Convert Sample Json to Go Struct

type PairsData []struct {
	ID            string `json:"id"`
	BaseCurrency  string `json:"baseCurrency"`
	QuoteCurrency string `json:"quoteCurrency"`
	LotSize       string `json:"lotSize"`
	TickSize      string `json:"tickSize"`
}

type CoinsData []struct {
	ID                  string `json:"id"`
	FullName            string `json:"fullName"`
	Crypto              bool   `json:"crypto"`
	DepositStatus       bool   `json:"depositStatus"`
	DepositConfirmation int    `json:"depositConfirmation"`
	WithdrawStatus      bool   `json:"withdrawStatus"`
	WithdrawFee         string `json:"withdrawFee"`
}

type OrderBook struct {
	Bids      [][]interface{} `json:"bids"`
	Asks      [][]interface{} `json:"asks"`
	Timestamp float64         `json:"timestamp"`
	// Message      type         `json:"msg"`
}
