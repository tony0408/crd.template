package exchange

import (
	"sync"

	"../coin"
	"../market"
	"../pair"
	"../user"
)

type Exchange interface {
	GetName() ExchangeName
	GetTradingWebURL(pair *pair.Pair) string

	// FixSymbol()
	GetCode(symbol string) string //get own standard code  eg:  BTC|ABC
	GetSymbol(code string) string //get exchange symbol

	InitCoins()
	SetCoins() error
	GetCoins() []*coin.Coin
	InitPairs()
	SetPairs() error
	GetPairs() []*pair.Pair
	GetPairCode(pair *pair.Pair) string
	HasPair(*pair.Pair) bool

	GetMaker(pair *pair.Pair) (maker *market.Maker, err error)
	UpdateMaker(pair *pair.Pair, maker *market.Maker) error

	GetConstrainFetchMethod(pair *pair.Pair) *ConstrainFetchMethod

	GetLotSize(pair *pair.Pair) float64     //stepSize    for  quantity
	GetPriceFilter(pair *pair.Pair) float64 //tickSize    for  price

	GetFee(pair *pair.Pair) float64   //the exchange fee for the pair
	GetTxFee(coin *coin.Coin) float64 //the tx fee for withdraw the coin

	CanWithdraw(coin *coin.Coin) bool // is enable withdraw
	CanDeposit(coin *coin.Coin) bool  // is enable deposit

	Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool

	LimitSell(pair *pair.Pair, quantity, rate float64) (*market.Order, error)
	LimitBuy(pair *pair.Pair, quantity, rate float64) (*market.Order, error)

	OrderStatus(order *market.Order) error
	CancelOrder(order *market.Order) error
	CancelAllOrder() error //TODO need to impl cancel all order for exchanges
	ListOrders() (*[]market.Order, error)

	GetBalance(coin *coin.Coin) float64
	UpdateAllBalances()

	OrderBook(p *pair.Pair) (*market.Maker, error)

	UpdatePairConstrain()
	UpdateCoinConstrain()

	UpdateAllBalancesByUser(u *user.User)
}

type ExchangeManager struct{}

var instance *ExchangeManager
var once sync.Once

var exMap = make(map[ExchangeName]Exchange)
var exList = make([]Exchange, 0)
var supportList = make([]ExchangeName, 0)

func CreateExchangeManager() *ExchangeManager {
	once.Do(func() {
		instance = &ExchangeManager{}
		instance.init()
	})
	return instance
}

func (e *ExchangeManager) init() {
	e.initExchangeNames()
}

func (e *ExchangeManager) Add(exchange Exchange) {
	key := exchange.GetName()
	exMap[key] = exchange
	exList = append(exList, exchange)
}

func (e *ExchangeManager) GetSupportExchanges() []ExchangeName {
	return supportList
}

func (e *ExchangeManager) Get(name ExchangeName) Exchange {
	return exMap[name]
}

func (e *ExchangeManager) GetStr(name string) Exchange {
	for _, v := range e.GetSupportExchanges() {
		// log.Printf(" string(%v) == %v", string(v), name)
		if string(v) == name {
			return e.Get(v)
		}
	}
	return nil
}

func (e *ExchangeManager) Quantity() int {
	return len(exList)
}

func (e *ExchangeManager) GetById(i int) Exchange {
	return exList[i]
}

func (e *ExchangeManager) SubsetPairs(e1, e2 Exchange) []*pair.Pair {
	var pairs []*pair.Pair
	ep1 := e1.GetPairs()

	// log.Printf("ep1 [0]:%v ", ep1[0])

	for _, p := range ep1 {
		// log.Printf("%s:%v  --- %s:%v", e1.GetName(), p, e2.GetName(), e2.HasPair(p))

		if e2.HasPair(p) {
			// log.Printf("%s:%v  --- %s:%v", e1.GetName(), p, e2.GetName(), e2.HasPair(p))
			// if p.Base.Code == "USDT" {
			// 	log.Printf("%v - %v", 1, p)
			// }
			pairs = append(pairs, p)
		}
	}

	return pairs
}
