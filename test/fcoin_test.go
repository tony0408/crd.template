package test

import (
	"log"
	"testing"

	"../coin"
	"../exchange"
	"../exchange/fcoin"
	"../pair"
	"github.com/davecgh/go-spew/spew"
)

/********************API********************/
func Test_Fcoin_Balance(t *testing.T) {
	e := initFcoin()
	e.UpdateAllBalances()

	for k, v := range e.GetPairs() { // pairs from binance
		if v != nil {
			base := e.GetBalance(v.Base)
			target := e.GetBalance(v.Target)
			log.Printf("%d  v:%v  #  %v:%v  %v:%v", k, v.Name, v.Base.Code, base, v.Target.Code, target)
		}
	}
}

func Test_Fcoin_Withdraw(t *testing.T) {
	e := initFcoin()
	c := coin.GetCoin("BTC")
	amount := 0.0
	addr := "Address"
	tag := ""
	if e.Withdraw(c, amount, addr, tag) {
		log.Printf("Fcoin %s Withdraw Successful!", c.Code)
	}
}

func Test_Fcoin_Trade(t *testing.T) {
	e := initFcoin()
	p := pair.GetPair(coin.GetCoin("BTC"), coin.GetCoin("ETH"))
	log.Printf("=========in test===%+v=========", p.Target)
	rate := 0.00071300
	quantity := 1.0

	order, err := e.LimitBuy(p, quantity, rate)
	if err == nil {
		log.Printf("Fcoin Limit Buy: %v", order)
	} else {
		log.Printf("Fcoin Limit Buy Err: %s", err)
	}

	err = e.OrderStatus(order)
	if err == nil {
		log.Printf("Fcoin Order Status: %v", order)
	} else {
		log.Printf("Fcoin Order Status Err: %s", err)
	}

	err = e.CancelOrder(order)
	if err == nil {
		log.Printf("Fcoin Cancel Order: %v", order)
	} else {
		log.Printf("Fcoin Cancel Err: %s", err)
	}
}

func Test_Fcoin_OrderBook(t *testing.T) {
	e := initFcoin()

	for _, pair := range e.GetPairs() { // pairs from binance
		if pair != nil {
			orderbook, err := e.OrderBook(pair)
			if err == nil {
				log.Printf("%s: %+v", pair.Name, orderbook)
				//log.Printf("%s: ", pair.Name)
			}
		}

	}
}

/********************General********************/
func Test_Fcoin_ConstrainFetch(t *testing.T) {
	e := initFcoin()

	p := pair.GetPair(coin.GetCoin("BTC"), coin.GetCoin("BCH"))

	status := e.GetConstrainFetchMethod(p)
	// "Binance ConstrainFetchMethod: %v",
	spew.Dump(status)
}

func Test_Fcoin_Constrain(t *testing.T) {
	e := initFcoin()

	pair := pair.GetPairByKey("BTC|ETH")
	coinName := coin.GetCoin(pair.Target.Code)
	log.Printf("Taker Fee: %.8f", e.GetTxFee(coinName))
	log.Printf("Withdraw Fee: %.8f", e.GetFee(pair))
	log.Printf("Lot Size: %.8f", e.GetLotSize(pair))
	log.Printf("Price Filter: %.8f", e.GetPriceFilter(pair))
	log.Printf("Withdraw: %v", e.CanWithdraw(coinName))
	log.Printf("Deposit: %v", e.CanDeposit(coinName))
}

func Test_Fcoin_GetMaker(t *testing.T) {
	e := initFcoin()

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

func initFcoin() exchange.Exchange {
	pair.Init()
	config := &exchange.Config{}
	config.RedisServer = "RedisAddr:Port"
	config.RedisDB = 0
	config.API_KEY = "7eb8721a8fb341718035abc2335d8545"
	config.API_SECRET = "523d79a231ed44a981e1056151861b09"
	ex := fcoin.CreateFcoin(config)
	log.Printf("Initial [ %v ]", ex.GetName())
	config = nil
	return ex
}
