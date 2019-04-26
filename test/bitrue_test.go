package test

import (
	"log"
	"testing"

	"../coin"
	"../exchange"
	"../exchange/bitrue"
	"../pair"
	"github.com/davecgh/go-spew/spew"
)

/********************API********************/
func Test_Bitrue_Balance(t *testing.T) {
	e := initBitrue()
	e.UpdateAllBalances()

	for k, v := range e.GetPairs() { // pairs from binance
		if v != nil {
			base := e.GetBalance(v.Base)
			target := e.GetBalance(v.Target)
			log.Printf("%d  v:%v  #  %v:%v  %v:%v", k, v.Name, v.Base.Code, base, v.Target.Code, target)
		}
	}
}

func Test_Bitrue_Withdraw(t *testing.T) {
	e := initBitrue()
	c := coin.GetCoin("BTC")
	amount := 0.0
	addr := "Address"
	tag := ""
	if e.Withdraw(c, amount, addr, tag) {
		log.Printf("Bitrue %s Withdraw Successful!", c.Code)
	}
}

func Test_Bitrue_Trade(t *testing.T) {
	e := initBitrue()
	p := pair.GetPair(coin.GetCoin("BTC"), coin.GetCoin("ETH"))
	rate := 0.071300
	quantity := 1.0

	order, err := e.LimitBuy(p, quantity, rate)
	if err == nil {
		log.Printf("Bitrue Limit Buy: %v", order)
	} else {
		log.Printf("Bitrue Limit Buy Err: %s", err)
	}

	err = e.OrderStatus(order)
	if err == nil {
		log.Printf("Bitrue Order Status: %v", order)
	} else {
		log.Printf("Bitrue Order Status Err: %s", err)
	}

	err = e.CancelOrder(order)
	if err == nil {
		log.Printf("Bitrue Cancel Order: %v", order)
	} else {
		log.Printf("Bitrue Cancel Err: %s", err)
	}
}

func Test_Bitrue_OrderBook(t *testing.T) {
	e := initBitrue()

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
func Test_Bitrue_ConstrainFetch(t *testing.T) {
	e := initBitrue()

	p := pair.GetPair(coin.GetCoin("BTC"), coin.GetCoin("BCH"))

	status := e.GetConstrainFetchMethod(p)
	// "Binance ConstrainFetchMethod: %v",
	spew.Dump(status)
}

func Test_Bitrue_Constrain(t *testing.T) {
	e := initBitrue()

	pair := pair.GetPairByKey("BTC|ETH")
	coinName := coin.GetCoin(pair.Target.Code)
	log.Printf("Taker Fee: %.8f", e.GetTxFee(coinName))
	log.Printf("Withdraw Fee: %.8f", e.GetFee(pair))
	log.Printf("Lot Size: %.8f", e.GetLotSize(pair))
	log.Printf("Price Filter: %.8f", e.GetPriceFilter(pair))
	log.Printf("Withdraw: %v", e.CanWithdraw(coinName))
	log.Printf("Deposit: %v", e.CanDeposit(coinName))
}

func Test_Bitrue_GetMaker(t *testing.T) {
	e := initBitrue()

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
func initBitrue() exchange.Exchange {
	pair.Init()
	config := &exchange.Config{}
	config.RedisServer = "RedisAddr:Port"
	config.RedisDB = 0
	config.API_KEY = "aa4fd929d856095d8e5da4a9e9e67d1e6445a66c46f73feb689227066210ffc9"
	config.API_SECRET = "eae894d2c61c8ecdda593d9d637f5edb4b79da8ab369831741222a2a9288cc09"
	ex := bitrue.CreateBitrue(config)
	log.Printf("Initial [ %v ]", ex.GetName())
	config = nil
	return ex
}
