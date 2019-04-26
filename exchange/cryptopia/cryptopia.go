package cryptopia

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	cmap "github.com/orcaman/concurrent-map"

	"../../coin"
	"../../db"
	"../../exchange"
	"../../market"
	"../../pair"
)

type Cryptopia struct {
	Name         string `bson:"name"`
	Website      string `bson:"website"`
	RedisManager *db.RedisManager
	RedisServer  string
	RedisDB      int
	API_KEY      string
	API_SECRET   string
}

var pairList = make([]*pair.Pair, 0) //the last num is the number of pairs on this exchange
var coinList = make([]*coin.Coin, 0)
var balanceMap cmap.ConcurrentMap

// var balanceMap = make(map[*coin.Coin]float64)

var instance *Cryptopia
var once sync.Once

/***************************************************/
func CreateCryptopia(config *exchange.Config) *Cryptopia {
	once.Do(func() {
		instance = &Cryptopia{}
		instance.Name = "Cryptopia"
		instance.Website = "https://www.cryptopia.co.nz/"

		instance.RedisManager = db.CreateRedisManager()
		instance.RedisServer = config.RedisServer
		instance.RedisDB = config.RedisDB

		instance.API_KEY = config.API_KEY
		instance.API_SECRET = config.API_SECRET

		if balanceMap == nil {
			balanceMap = cmap.New()
		}

		instance.FixSymbol()
		instance.InitCoins()
		instance.InitPairs()
	})
	return instance
}

func (e *Cryptopia) GetMakerDB() *db.Redis {
	key := string(exchange.CRYPTOPIA)
	d := e.RedisManager.Get(key)
	if d == nil {
		d = db.CreateRedis()
		d.Init(instance.RedisServer, instance.RedisDB)
		e.RedisManager.Add(key, d)
	}
	return d
}

func (e *Cryptopia) InitPairs() {
	pairData := GetCryptopiaPair()

	for _, symbol := range *pairData {
		//Modify according to type and structure
		base := coin.GetCoin(e.GetCode(symbol.BaseSymbol))
		target := coin.GetCoin(e.GetCode(symbol.Symbol))
		if base != nil && target != nil {
			pair := pair.GetPair(base, target)
			pairList = append(pairList, pair)
		}
	}
}

func (e *Cryptopia) InitCoins() {
	coinInfo := GetCryptopiaCoin()

	for _, data := range *coinInfo {
		//Modify according to type and structure
		c := coin.GetCoin(e.GetCode(data.Symbol))
		if c == nil || c.Name == "" {
			c = &coin.Coin{}
			c.Code = e.GetCode(data.Symbol)
			c.Name = data.Name
			coin.AddCoin(c)
		}
		coinList = append(coinList, c)
	}
}

/***************************************************/
func (e *Cryptopia) UpdateMaker(pair *pair.Pair, maker *market.Maker) error {
	m, err := json.Marshal(maker)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s-%s", exchange.CRYPTOPIA, pair.Name)
	return e.GetMakerDB().Set(key, string(m))
}

func (e *Cryptopia) GetMaker(pair *pair.Pair) (maker *market.Maker, err error) {
	key := fmt.Sprintf("%s-%s", exchange.CRYPTOPIA, pair.Name)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cryptopia does not have the pair : %v", pair.Name))
	}
	if str, ok := val.(string); ok {

		if err := json.Unmarshal([]byte(str), &maker); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Cryptopia GetMaker Key: %v can't convert to string: %v", key, val))
	}
	return maker, err
}

/***************************************************/
func (e *Cryptopia) SetCoins() error {
	return nil
}

func (e *Cryptopia) GetCoins() []*coin.Coin {
	return coinList
}

func (e *Cryptopia) SetPairs() error {
	return nil
}

func (e *Cryptopia) GetPairs() []*pair.Pair {
	return pairList
}

/*Get Exchange A Pair
Step 1: Change Instance Name    (e *<exchange Instance Name>)*/
func (e *Cryptopia) GetPair(key string) *pair.Pair {
	return nil
}

func (e *Cryptopia) GetPairCode(pair *pair.Pair) string {
	//Modify according to Exchange Request
	code := fmt.Sprintf("%s_%s", strings.ToUpper(e.GetSymbol(pair.Target.Code)), strings.ToUpper(e.GetSymbol(pair.Base.Code)))
	return code
}

func (e *Cryptopia) HasPair(pair *pair.Pair) bool {
	m, err := e.GetMaker(pair)
	if err == nil && m != nil && m.Bids != nil {
		return true
	}
	return false
}

/*************** pairs on the exchanges ***************/
func (e *Cryptopia) GetName() exchange.ExchangeName {
	return exchange.CRYPTOPIA
}

func (e *Cryptopia) GetFee(pair *pair.Pair) float64 { // Taker fee for each coin
	return 0.002 //Taker Fee: 0.2%
}

func (e *Cryptopia) GetLotSize(pair *pair.Pair) float64 { // stepSize for quantity
	key := fmt.Sprintf("%s-Constrain-%s", exchange.CRYPTOPIA, pair.Name)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		log.Printf("Cryptopia GetLotSize Key: %v Err: %s\n", key, err)
		return 0.00000001
	}
	constrain := exchange.PairConstrain{}
	if str, ok := val.(string); ok {
		if err := json.Unmarshal([]byte(str), &constrain); err != nil {
			log.Printf("Cryptopia GetLotSize Key: %v Unmarshal Err: %s\n", key, err)
			return 0.00000001
		}
	} else {
		log.Printf("Cryptopia GetLotSize Key: %v can't convert to string: %v", key, val)
		return 0.00000001
	}
	return constrain.LotSize //return 0.00100000
}
func (e *Cryptopia) GetPriceFilter(pair *pair.Pair) float64 { // tickSize for price
	key := fmt.Sprintf("%s-Constrain-%s", exchange.CRYPTOPIA, pair.Name)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		log.Printf("Cryptopia GetPriceFilter Key: %v Err: %s\n", key, err)
		return 0.00000001
	}
	constrain := exchange.PairConstrain{}
	if str, ok := val.(string); ok {
		if err := json.Unmarshal([]byte(str), &constrain); err != nil {
			log.Printf("Cryptopia GetPriceFilter Key: %v Unmarshal Err: %s\n", key, err)
			return 0.00000001
		}
	} else {
		log.Printf("Cryptopia GetPriceFilter Key: %v can't convert to string: %v", key, val)
		return 0.00000001
	}
	return constrain.TickSize
}

func (e *Cryptopia) GetConstrainFetchMethod(pair *pair.Pair) *exchange.ConstrainFetchMethod {
	constrainFetchMethod := &exchange.ConstrainFetchMethod{}
	constrainFetchMethod.Fee = true
	constrainFetchMethod.LotSize = true
	constrainFetchMethod.TickSize = true
	constrainFetchMethod.TxFee = true
	constrainFetchMethod.Withdraw = true
	constrainFetchMethod.Deposit = true
	constrainFetchMethod.Confirmation = true
	return constrainFetchMethod
}

/*************** coins on the exchanges ***************/
func (e *Cryptopia) GetBalance(coin *coin.Coin) float64 {
	if tmp, ok := balanceMap.Get(coin.Code); ok {
		return tmp.(float64)
	} else {
		return 0.0
	}
}

func (e *Cryptopia) GetTxFee(coin *coin.Coin) float64 { // Withdraw Fee
	key := fmt.Sprintf("%s-Constrain-%s", exchange.CRYPTOPIA, coin.Code)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		log.Printf("Cryptopia GetTxFee Key: %v Err: %s\n", key, err)
		return 100.001
	}
	constrain := exchange.CoinConstrain{}
	if str, ok := val.(string); ok {
		if err := json.Unmarshal([]byte(str), &constrain); err != nil {
			log.Printf("Cryptopia GetTxFee Key: %v Unmarshal Err: %s\n", key, err)
			return 100.001
		}
	} else {
		log.Printf("Cryptopia GetConfirmaGetTxFeetion Key: %v can't convert to string: %v", key, val)
		return 100.001
	}
	return constrain.TxFee
}

func (e *Cryptopia) GetConfirmation(coin *coin.Coin) int { // deposit confirmations
	key := fmt.Sprintf("%s-Constrain-%s", exchange.CRYPTOPIA, coin.Code)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		log.Printf("Cryptopia GetConfirmation Key: %v Err: %s\n", key, err)
		return 1001
	}
	constrain := exchange.CoinConstrain{}
	if str, ok := val.(string); ok {
		if err := json.Unmarshal([]byte(str), &constrain); err != nil {
			log.Printf("Cryptopia GetConfirmation Key: %v Unmarshal Err: %s\n", key, err)
			return 1001
		}
	} else {
		log.Printf("Cryptopia GetConfirmation Key: %v can't convert to string: %v", key, val)
		return 1001
	}
	return constrain.Confirmation
}

func (e *Cryptopia) CanWithdraw(coin *coin.Coin) bool { // does withdraw enable
	key := fmt.Sprintf("%s-Constrain-%s", exchange.CRYPTOPIA, coin.Code)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		log.Printf("Cryptopia CanWithdraw Key: %v Err: %s\n", key, err)
		return true
	}
	constrain := exchange.CoinConstrain{}
	if str, ok := val.(string); ok {
		if err := json.Unmarshal([]byte(str), &constrain); err != nil {
			log.Printf("Cryptopia CanWithdraw Key: %v Unmarshal Err: %s\n", key, err)
			return true
		}
	} else {
		log.Printf("Cryptopia CanWithdraw Key: %v can't convert to string: %v", key, val)
		return true
	}
	return constrain.Withdraw
}
func (e *Cryptopia) CanDeposit(coin *coin.Coin) bool { // does deposit enable
	key := fmt.Sprintf("%s-Constrain-%s", exchange.CRYPTOPIA, coin.Code)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		log.Printf("Cryptopia CanDeposit Key: %v Err: %s\n", key, err)
		return true
	}
	constrain := exchange.CoinConstrain{}
	if str, ok := val.(string); ok {
		if err := json.Unmarshal([]byte(str), &constrain); err != nil {
			log.Printf("Cryptopia CanDeposit Key: %v Unmarshal Err: %s\n", key, err)
			return true
		}
	} else {
		log.Printf("Cryptopia CanDeposit Key: %v can't convert to string: %v", key, val)
		return true
	}
	return constrain.Deposit
}
func (e *Cryptopia) GetTradingWebURL(pair *pair.Pair) string {
	return fmt.Sprintf("https://www.cryptopia.co.nz/Exchange/?market=%s_%s", strings.ToUpper(pair.Target.Code), strings.ToUpper(pair.Base.Code))

}
