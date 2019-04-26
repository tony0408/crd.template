package test

import (
	"log"
	"testing"

	"../coin"
	"../exchange"
	"../exchange/kraken"
	"../pair"
	"github.com/davecgh/go-spew/spew"
)

/********************API********************/
func Test_Kraken_Balance(t *testing.T) {
	e := initKraken()
	e.UpdateAllBalances()

	for k, v := range e.GetPairs() { // pairs from binance
		if v != nil {
			base := e.GetBalance(v.Base)
			target := e.GetBalance(v.Target)
			log.Printf("%d  v:%v  #  %v:%v  %v:%v", k, v.Name, v.Base.Code, base, v.Target.Code, target)
		}
	}
}

func Test_Kraken_Withdraw(t *testing.T) {
	e := initKraken()
	c := coin.GetCoin("BTC")
	amount := 0.0
	addr := "Address"
	tag := ""
	if e.Withdraw(c, amount, addr, tag) {
		log.Printf("Kraken %s Withdraw Successful!", c.Code)
	}
}

func Test_Kraken_Trade(t *testing.T) {
	e := initKraken()
	p := pair.GetPair(coin.GetCoin("BTC"), coin.GetCoin("ETH"))
	rate := 0.00071
	quantity := 1.0

	order, err := e.LimitBuy(p, quantity, rate)
	if err == nil {
		log.Printf("Kraken Limit Buy: %v", order)
	} else {
		log.Printf("Kraken Limit Buy Err: %s", err)
	}

	err = e.OrderStatus(order)
	if err == nil {
		log.Printf("Kraken Order Status: %v", order)
	} else {
		log.Printf("Kraken Order Status Err: %s", err)
	}

	err = e.CancelOrder(order)
	if err == nil {
		log.Printf("Kraken Cancel Order: %v", order)
	} else {
		log.Printf("Kraken Cancel Err: %s", err)
	}
}

func Test_Kraken_OrderBook(t *testing.T) {
	e := initKraken()

	for _, pair := range e.GetPairs() { // pairs from binance
		if pair != nil {
			orderbook, err := e.OrderBook(pair)
			if err == nil {
				log.Printf("%s: %+v", pair.Name, orderbook)
			}
		}

	}
}

/********************General********************/
func Test_Kraken_ConstrainFetch(t *testing.T) {
	e := initKraken()

	p := pair.GetPair(coin.GetCoin("BTC"), coin.GetCoin("BCH"))

	status := e.GetConstrainFetchMethod(p)
	// "Binance ConstrainFetchMethod: %v",
	spew.Dump(status)
}

func Test_Kraken_Constrain(t *testing.T) {
	e := initKraken()

	pair := pair.GetPairByKey("BTC|ETH")
	coinName := coin.GetCoin(pair.Target.Code)
	log.Printf("Taker Fee: %.8f", e.GetTxFee(coinName))
	log.Printf("Withdraw Fee: %.8f", e.GetFee(pair))
	log.Printf("Lot Size: %.8f", e.GetLotSize(pair))
	log.Printf("Price Filter: %.8f", e.GetPriceFilter(pair))
	log.Printf("Withdraw: %v", e.CanWithdraw(coinName))
	log.Printf("Deposit: %v", e.CanDeposit(coinName))
}

func Test_Kraken_GetMaker(t *testing.T) {
	e := initKraken()

	pair := pair.GetPairByKey("BTC|ETH")
	maker, _ := e.GetMaker(pair)
	log.Printf("Pair Code: %s", e.GetPairCode(pair))
	log.Printf("Maker: %v", maker)
}

/* Modify Config
Redis Server: xx.xx.xx.xx:xxxx
Redis DB: DB Number
API Key: Exchange API Key
API Secret: Exchange API Secret Key */
func initKraken() exchange.Exchange {
	pair.Init()
	config := &exchange.Config{}
	config.RedisServer = "RedisAddr:Port"
	config.RedisDB = 0
	config.API_KEY = "qgUnBhOAS3k7WNOVqP8MFcwEe5zaOmWjVxSx5j22NUQn1IFWMk0dASdv"

	config.API_SECRET = "JK/s7nUqNTbOiOBSa2NW1G1EP00t7Qd2iq4BMrQ4y1xS02xx2mwJ96l5w2ExB3x6E4H8DtIMLUFun5c8+pLSQA=="
	ex := kraken.CreateKraken(config)
	log.Printf("Initial [ %v ]", ex.GetName())
	config = nil
	return ex
}
