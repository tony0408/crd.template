package coineal

import "encoding/json"

type CoinealOrderBook struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Tick struct {
			Asks [][]interface{} `json:"asks"`
			Bids [][]interface{} `json:"bids"`
			Time interface{}     `json:"time"`
		} `json:"tick"`
	} `json:"data"`
}
type CoinealCoins struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Symbol          string `json:"symbol"`
		CountCoin       string `json:"count_coin"`
		AmountPrecision int    `json:"amount_precision"`
		BaseCoin        string `json:"base_coin"`
		PricePrecision  int    `json:"price_precision"`
	} `json:"data"`
}
type JsonResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}
type CoinealBalance struct {
	TotalAsset string `json:"total_asset"`
	CoinList   []struct {
		Normal      interface{} `json:"normal"`
		BtcValuatin interface{} `json:"btcValuatin"`
		Locked      interface{} `json:"locked"`
		Coin        interface{} `json:"coin"`
	} `json:"coin_list"`
}
type CoinealOrder struct {
	OrderID int `json:"order_id"`
}
type CoinealOrderStatus struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TradeList []interface{} `json:"trade_list"`
		OrderInfo struct {
			Side         string        `json:"side"`
			TotalPrice   string        `json:"total_price"`
			CreatedAt    int64         `json:"created_at"`
			AvgPrice     string        `json:"avg_price"`
			CountCoin    string        `json:"countCoin"`
			Source       int           `json:"source"`
			Type         int           `json:"type"`
			SideMsg      string        `json:"side_msg"`
			Volume       string        `json:"volume"`
			Price        string        `json:"price"`
			SourceMsg    string        `json:"source_msg"`
			StatusMsg    string        `json:"status_msg"`
			DealVolume   string        `json:"deal_volume"`
			ID           int           `json:"id"`
			RemainVolume string        `json:"remain_volume"`
			BaseCoin     string        `json:"baseCoin"`
			TradeList    []interface{} `json:"tradeList"`
			Status       int           `json:"status"`
		} `json:"order_info"`
	} `json:"data"`
}
