package kraken

import "encoding/json"

//Get Struct by Exchange
//Convert Sample Json to Go Struct

type ResponseReturn struct {
	Error  []interface{}   `json:"error"`
	Result json.RawMessage `json:"result"`
}

type CoinData struct {
	Aclass          string `json:"aclass"`
	Altname         string `json:"altname"`
	Decimals        int    `json:"decimals"`
	DisplayDecimals int    `json:"display_decimals"`
}

type PairData struct {
	Altname           string        `json:"altname"`
	AclassBase        string        `json:"aclass_base"`
	Base              string        `json:"base"`
	AclassQuote       string        `json:"aclass_quote"`
	Quote             string        `json:"quote"`
	Lot               string        `json:"lot"`
	PairDecimals      int           `json:"pair_decimals"`
	LotDecimals       int           `json:"lot_decimals"`
	LotMultiplier     int           `json:"lot_multiplier"`
	LeverageBuy       []interface{} `json:"leverage_buy"`
	LeverageSell      []interface{} `json:"leverage_sell"`
	Fees              [][]float64   `json:"fees"`
	FeesMaker         [][]float64   `json:"fees_maker"`
	FeeVolumeCurrency string        `json:"fee_volume_currency"`
	MarginCall        int           `json:"margin_call"`
	MarginStop        int           `json:"margin_stop"`
}

type OrderBook struct {
	Bids [][]interface{} `json:"bids"`
	Asks [][]interface{} `json:"asks"`
}
