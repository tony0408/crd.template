package test

import (
	"log"
	"testing"
	"time"

	"../coin"
	"../exchange"
	"../exchange/cryptopia"
	"../pair"
	"github.com/davecgh/go-spew/spew"
)

/********************API********************/
func Test_Cryptopia_Balance(t *testing.T) {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	time.LoadLocation("America/Vancouver")

	e := initCryptopia()
	e.UpdateAllBalances()

	for k, v := range e.GetPairs() { // pairs from binance
		if v != nil {
			base := e.GetBalance(v.Base)
			target := e.GetBalance(v.Target)
			log.Printf("%d  v:%v  #  %v:%v  %v:%v", k, v.Name, v.Base.Code, base, v.Target.Code, target)
		}
	}
}

func Test_Cryptopia_Withdraw(t *testing.T) {
	e := initCryptopia()
	c := coin.GetCoin("BTC")
	amount := 0.0
	addr := "Address"
	tag := ""
	if e.Withdraw(c, amount, addr, tag) {
		log.Printf("Cryptopia %s Withdraw Successful!", c.Code)
	}
}

func Test_Cryptopia_Trade(t *testing.T) {
	e := initCryptopia()
	p := pair.GetPair(coin.GetCoin("BTC"), coin.GetCoin("ETH"))
	rate := 0.5
	quantity := 1.0

	// Place Order in 2 Second
	order, err := e.LimitBuy(p, quantity, rate)
	if err == nil {
		log.Printf("Cryptopia Limit Buy: %v", order)
	} else {
		log.Printf("Cryptopia Limit Buy Err: %s", err)
	}

	err = e.OrderStatus(order)
	if err == nil {
		log.Printf("Cryptopia Order Status: %v", order)
	} else {
		log.Printf("Cryptopia Order Status Err: %s", err)
	}

	err = e.CancelOrder(order)
	if err == nil {
		log.Printf("Cryptopia Cancel Order: %v", order)
	} else {
		log.Printf("Cryptopia Cancel Err: %s", err)
	}
}

func Test_Cryptopia_OrderBook(t *testing.T) {
	e := initCryptopia()

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
func Test_Cryptopia_ConstrainFetch(t *testing.T) {
	e := initCryptopia()

	p := pair.GetPair(coin.GetCoin("BTC"), coin.GetCoin("BCH"))

	status := e.GetConstrainFetchMethod(p)
	// "Binance ConstrainFetchMethod: %v",
	spew.Dump(status)
}

func Test_Cryptopia_Constrain(t *testing.T) {
	e := initCryptopia()

	pair := pair.GetPairByKey("BTC|ETH")
	coinName := coin.GetCoin(pair.Target.Code)
	log.Printf("Taker Fee: %.8f", e.GetTxFee(coinName))
	log.Printf("Withdraw Fee: %.8f", e.GetFee(pair))
	log.Printf("Lot Size: %.8f", e.GetLotSize(pair))
	log.Printf("Price Filter: %.8f", e.GetPriceFilter(pair))
	log.Printf("Withdraw: %v", e.CanWithdraw(coinName))
	log.Printf("Deposit: %v", e.CanDeposit(coinName))
}

func Test_Cryptopia_GetMaker(t *testing.T) {
	e := initCryptopia()

	pair := pair.GetPairByKey("BTC|ETH")
	maker, _ := e.GetMaker(pair)
	log.Printf("Pair Code: %s", e.GetPairCode(pair))
	log.Printf("Maker: %v", maker)
}

func initCryptopia() exchange.Exchange {
	pair.Init()
	config := &exchange.Config{}
	config.RedisServer = ""
	config.RedisDB = 9
	config.API_KEY = ""
	config.API_SECRET = ""
	ex := cryptopia.CreateCryptopia(config)
	log.Printf("Initial [ %v ]", ex.GetName())
	config = nil
	return ex
}
