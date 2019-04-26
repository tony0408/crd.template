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

/* type AccountBalances map[string]*AccountDetails
type AccountDetails float64
type AccountDetail2 struct {
	name float64
} */

type AccountBalances struct {
	ADA  float64 `json:"ADA,string"`
	BCH  float64 `json:"BCH,string"`
	DASH float64 `json:"DASH,string"`
	EOS  float64 `json:"EOS,string"`
	GNO  float64 `json:"GNO,string"`
	QTUM float64 `json:"QTUM,string"`
	KFEE float64 `json:"KFEE,string"`
	USDT float64 `json:"USDT,string"`
	XDAO float64 `json:"XDAO,string"`
	XETC float64 `json:"XETC,string"`
	XETH float64 `json:"XETH,string"`
	XICN float64 `json:"XICN,string"`
	XLTC float64 `json:"XLTC,string"`
	XMLN float64 `json:"XMLN,string"`
	XNMC float64 `json:"XNMC,string"`
	XREP float64 `json:"XREP,string"`
	XXBT float64 `json:"XXBT,string"`
	XXDG float64 `json:"XXDG,string"`
	XXLM float64 `json:"XXLM,string"`
	XXMR float64 `json:"XXMR,string"`
	XXRP float64 `json:"XXRP,string"`
	XTZ  float64 `json:"XTZ,string"`
	XXVN float64 `json:"XXVN,string"`
	XZEC float64 `json:"XZEC,string"`
	ZCAD float64 `json:"ZCAD,string"`
	ZEUR float64 `json:"ZEUR,string"`
	ZGBP float64 `json:"ZGBP,string"`
	ZJPY float64 `json:"ZJPY,string"`
	ZKRW float64 `json:"ZKRW,string"`
	ZUSD float64 `json:"ZUSD,string"`
}

type Order struct {
	TransactionID  string           `json:"-"`
	ReferenceID    string           `json:"refid"`
	UserRef        int              `json:"userref"`
	Status         string           `json:"status"`
	OpenTime       float64          `json:"opentm"`
	StartTime      float64          `json:"starttm"`
	ExpireTime     float64          `json:"expiretm"`
	Description    string 			`json:"descr"`
	Volume         string           `json:"vol"`
	VolumeExecuted float64          `json:"vol_exec,string"`
	Cost           float64          `json:"cost,string"`
	Fee            float64          `json:"fee,string"`
	Price          float64          `json:"price,string"`
	StopPrice      float64          `json:"stopprice.string"`
	LimitPrice     float64          `json:"limitprice,string"`
	Misc           string           `json:"misc"`
	OrderFlags     string           `json:"oflags"`
	CloseTime      float64          `json:"closetm"`
	Reason         string           `json:"reason"`
}

type AddOrderResponse struct {
	Description    string 			`json:"descr"`
	TransactionIds []string         `json:"txid"`
}
