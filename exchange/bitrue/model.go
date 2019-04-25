package bitrue

import "encoding/json"

type JsonResponse struct {
	Status int             `json:"code"` //success: 200
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
}

type BitrueOrderBook struct {
	LastUpdateID int64           `json:"lastUpdateId"`
	Bids         [][]interface{} `json:"bids"`
	Asks         [][]interface{} `json:"asks"`
}
type BitruePair struct {
	Timezone   string `json:"timezone"`
	ServerTime int64  `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []struct {
		Symbol             string   `json:"symbol"`
		Status             string   `json:"status"`
		BaseAsset          string   `json:"baseAsset"`
		BaseAssetPrecision int      `json:"baseAssetPrecision"`
		QuoteAsset         string   `json:"quoteAsset"`
		QuotePrecision     int      `json:"quotePrecision"`
		OrderTypes         []string `json:"orderTypes"`
		IcebergAllowed     bool     `json:"icebergAllowed"`
		Filters            []struct {
			FilterType  string `json:"filterType"`
			MinPrice    string `json:"minPrice,omitempty"`
			MaxPrice    string `json:"maxPrice,omitempty"`
			PriceScale  int    `json:"priceScale,omitempty"`
			MinQty      string `json:"minQty,omitempty"`
			MaxQty      string `json:"maxQty,omitempty"`
			VolumeScale int    `json:"volumeScale,omitempty"`
		} `json:"filters"`
	} `json:"symbols"`
}

type PlaceOrder struct {
	Symbol        string `json:"symbol"`
	OrderID       int    `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
	TransactTime  int64  `json:"transactTime"`
}

type TradeHistory struct {
	Symbol              string `json:"symbol"`
	OrderID             int    `json:"orderId"`
	ClientOrderID       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StopPrice           string `json:"stopPrice"`
	IcebergQty          string `json:"icebergQty"`
	Time                int64  `json:"time"`
	UpdateTime          int64  `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
}

type AccountBalances struct {
	MakerCommission  int  `json:"makerCommission"`
	TakerCommission  int  `json:"takerCommission"`
	BuyerCommission  int  `json:"buyerCommission"`
	SellerCommission int  `json:"sellerCommission"`
	CanTrade         bool `json:"canTrade"`
	CanWithdraw      bool `json:"canWithdraw"`
	CanDeposit       bool `json:"canDeposit"`
	UpdateTime       int  `json:"updateTime"`
	Balances         []struct {
		Asset  string `json:"asset"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
	} `json:"balances"`
}

type CancelOrder struct {
	Symbol            string `json:"symbol"`
	OrigClientOrderID string `json:"origClientOrderId"`
	OrderID           int    `json:"orderId"`
	ClientOrderID     string `json:"clientOrderId"`
}
