package test

import (
	"log"
	"testing"

	"../coin"
	"../exchange"
	"../exchange/coineal"
	"../pair"
	"github.com/davecgh/go-spew/spew"
)

/********************API********************/
func Test_Coineal_Balance(t *testing.T) {
	e := initCoineal()
	e.UpdateAllBalances()

	for k, v := range e.GetPairs() { // pairs from binance
		if v != nil {
			base := e.GetBalance(v.Base)
			target := e.GetBalance(v.Target)
			log.Printf("%d  v:%v  #  %v:%v  %v:%v", k, v.Name, v.Base.Code, base, v.Target.Code, target)
		}
	}
}

func Test_Coineal_Withdraw(t *testing.T) {
	e := initCoineal()
	c := coin.GetCoin("BTC")
	amount := 0.0
	addr := "Address"
	tag := ""
	if e.Withdraw(c, amount, addr, tag) {
		log.Printf("Coineal %s Withdraw Successful!", c.Code)
	}
}

func Test_Coineal_Trade(t *testing.T) {
	e := initCoineal()
	p := pair.GetPair(coin.GetCoin("BTC"), coin.GetCoin("ZCL"))
	rate := 0.000001
	quantity := 1.0

	order, err := e.LimitBuy(p, quantity, rate)
	if err == nil {
		log.Printf("Coineal Limit Buy: %v", order)
	} else {
		log.Printf("Coineal Limit Buy Err: %s", err)
	}

	err = e.OrderStatus(order)
	if err == nil {
		log.Printf("Coineal Order Status: %v", order)
	} else {
		log.Printf("Coineal Order Status Err: %s", err)
	}

	err = e.CancelOrder(order)
	if err == nil {
		log.Printf("Coineal Cancel Order: %v", order)
	} else {
		log.Printf("Coineal Cancel Err: %s", err)
	}
}

func Test_Coineal_OrderBook(t *testing.T) {
	e := initCoineal()

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
func Test_Coineal_ConstrainFetch(t *testing.T) {
	e := initCoineal()

	p := pair.GetPair(coin.GetCoin("BTC"), coin.GetCoin("BCH"))

	status := e.GetConstrainFetchMethod(p)
	// "Binance ConstrainFetchMethod: %v",
	spew.Dump(status)
}

func Test_Coineal_Constrain(t *testing.T) {
	e := initCoineal()

	pair := pair.GetPairByKey("BTC|ETH")
	coinName := coin.GetCoin(pair.Target.Code)
	log.Printf("Taker Fee: %.8f", e.GetTxFee(coinName))
	log.Printf("Withdraw Fee: %.8f", e.GetFee(pair))
	log.Printf("Lot Size: %.8f", e.GetLotSize(pair))
	log.Printf("Price Filter: %.8f", e.GetPriceFilter(pair))
	log.Printf("Withdraw: %v", e.CanWithdraw(coinName))
	log.Printf("Deposit: %v", e.CanDeposit(coinName))
}

func Test_Coineal_GetMaker(t *testing.T) {
	e := initCoineal()

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
func initCoineal() exchange.Exchange {
	pair.Init()
	config := &exchange.Config{}
	config.RedisServer = "RedisAddr:Port"
	config.RedisDB = 0
	config.API_KEY = ""
	config.API_SECRET = ""
	ex := coineal.CreateCoineal(config)
	log.Printf("Initial [ %v ]", ex.GetName())
	config = nil
	return ex
}
