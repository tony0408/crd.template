package fcoin

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

type Fcoin struct {
	Name         string `bson:"name"`
	Website      string `bson:"website"`
	RedisManager *db.RedisManager
	RedisServer  string
	RedisDB      int
	API_KEY      string
	API_SECRET   string
	WalletStatus []exchange.Wallet_Stat
}

var pairList = make([]*pair.Pair, 0) //the last num is the number of pairs on this exchange
var coinList = make([]*coin.Coin, 0)
var balanceMap cmap.ConcurrentMap

// var balanceMap = make(map[*coin.Coin]float64)

var instance *Fcoin
var once sync.Once

/***************************************************/
/*Create New Exchange
Add Exchange Name(Capital Letter) to meta.go
Name: Exchange Name
Website: Exchange Website URL
MakerDB: Exchange Redis Server & Number(Import from Config)
Execute Coins & Pairs Initial
API_KEY: Import from Config
API_SECRET: Import from Config
WalletStatus: If API doesn't provide Wallet Status, import data from Postgres*/
func CreateFcoin(config *exchange.Config) *Fcoin {
	once.Do(func() {
		instance = &Fcoin{}
		instance.Name = "Fcoin"
		instance.Website = "https://www.fcoin.com/"

		instance.RedisManager = db.CreateRedisManager()
		instance.RedisServer = config.RedisServer
		instance.RedisDB = config.RedisDB

		instance.API_KEY = config.API_KEY
		instance.API_SECRET = config.API_SECRET

		instance.WalletStatus = config.WalletStatus

		if balanceMap == nil {
			balanceMap = cmap.New()
		}

		instance.FixSymbol()
		instance.InitCoins()
		instance.InitPairs()
	})
	return instance
}

func (e *Fcoin) GetMakerDB() *db.Redis {
	key := string(exchange.FCOIN)
	d := e.RedisManager.Get(key)
	if d == nil {
		d = db.CreateRedis()
		d.Init(instance.RedisServer, instance.RedisDB)
		e.RedisManager.Add(key, d)
	}
	return d
}

/*Initial the Pairs of Exchange
Step 1: Change Instance Name (e *<exchange Instance Name>)
Step 2: Get API Data
Step 3: Get Each Symbol
Step 4: Identify Base & Target
Step 5: Get Coin Standard Code ex. e.GetCode(base)
Step 6: Get Coin
Step 7: Add Pair to Exchange Pairs Arrary*/
func (e *Fcoin) InitPairs() {

	pairData := GetFcoinPair()
	if pairData != nil {
		//		fmt.Printf("pair not nil: %v ---", pairData.Data)
		for _, symbol := range pairData.Data {
			//Modify according to type and structure
			base := coin.GetCoin(e.GetCode(symbol.QuoteCurrency))
			target := coin.GetCoin(e.GetCode(symbol.BaseCurrency))
			if base != nil && target != nil {
				pair := pair.GetPair(base, target)
				pairList = append(pairList, pair)
			}
		}
	}

}

/*Initial the Coins of Exchange
Step 1: Change Instance Name (e *<exchange Instance Name>)
Step 2: Get API Data
Step 3: Get Each Coin
Step 4: Check the coin (Use Standard Code ex. e.GetCode(coin)) exists or not
Step 5: if the coin doesn't exist in coinmap, Add the coin in coinmap
	- Code: General Short Code
	!--Fill below if API provide the following information--!
	- Name: Coin Full Name
	- Website: Coin Official Website
	- Explorer: Coin Block Explorer
	- Health: the health of the chain
	- Blockheigh: the heigh of the chain
	- Blocktime: the time of the block created
	- Blocklast: the last block of the chain*/
func (e *Fcoin) InitCoins() {
	coinInfo := GetFcoinCoin()

	if coinInfo != nil {
		for _, data := range coinInfo.Data {
			//Modify according to type and structure
			c := coin.GetCoin(e.GetCode(data))
			if c == nil {
				c = &coin.Coin{}
				c.Code = e.GetCode(data)
				//				c.Name = data.FullName
				coin.AddCoin(c)
			}
			coinList = append(coinList, c)
		}
	}
}

/***************************************************/
/*Upload updated Maker to Redis
Step 1: Change Instance Name (e *<exchange Instance Name>)
Step 2: Change Exchange Name exchange.<Capital Letter Exchange Name>*/
func (e *Fcoin) UpdateMaker(pair *pair.Pair, maker *market.Maker) error {
	m, err := json.Marshal(maker)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s-%s", exchange.FCOIN, pair.Name)
	return e.GetMakerDB().Set(key, string(m))
}

/*Get Maker from Redis
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Change Exchange Name    exchange.<Capital Letter Exchange Name>
Step 3: Change Error Exchange Name    <exchange Name> does not have the pair*/
func (e *Fcoin) GetMaker(pair *pair.Pair) (maker *market.Maker, err error) {
	key := fmt.Sprintf("%s-%s", exchange.FCOIN, pair.Name)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Fcoin does not have the pair : %v", pair.Name))
	}
	if str, ok := val.(string); ok {

		if err := json.Unmarshal([]byte(str), &maker); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Fcoin GetMaker Key: %v can't convert to string: %v", key, val))
	}
	return maker, err
}

/***************************************************/
func (e *Fcoin) SetCoins() error {
	return nil
}

func (e *Fcoin) GetCoins() []*coin.Coin {
	return coinList
}

func (e *Fcoin) SetPairs() error {
	return nil
}

/*Get Exchange All Pairs
Step 1: Change Instance Name    (e *<exchange Instance Name>)*/
func (e *Fcoin) GetPairs() []*pair.Pair {
	return pairList
}

/*Get Pair Code base on Exchange
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Change Format of Code   ex. ADABTC in Binance, eos_btc in TradeSatoshi*/
func (e *Fcoin) GetPairCode(pair *pair.Pair) string {
	//Modify according to Exchange Request
	code := fmt.Sprintf("%s%s", strings.ToUpper(e.GetSymbol(pair.Target.Code)), strings.ToUpper(e.GetSymbol(pair.Base.Code)))
	return code
}

/*Check the exchange has the pair
Step 1: Change Instance Name    (e *<exchange Instance Name>)*/
func (e *Fcoin) HasPair(pair *pair.Pair) bool {
	m, err := e.GetMaker(pair)
	if err == nil && m != nil && m.Bids != nil {
		return true
	}
	return false
}

/*************** pairs on the exchanges ***************/
/*Get Exchange Name
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Change Exchange Name    exchange.<Capital Letter Exchange Name>*/
func (e *Fcoin) GetName() exchange.ExchangeName {
	return exchange.FCOIN
}

/*Get Exchange Taker Fee
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Change Return base on the taker fee that exchange provides*/
func (e *Fcoin) GetFee(pair *pair.Pair) float64 { // Taker fee for each coin
	return 0.002 //Taker Fee: 0.2%
}

/*Get Pair LotSize(Quantity)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2:
	Condition 1: API provides this information  --Refer Binance Code
		key: Constrain key in Redis ex. key := fmt.Sprintf("%s-Constrain-%s", exchange.<Capital Letter Exchange Name>, pair.Name)
		val: Get Redis Json Data ex. val, err := e.GetMakerDB().Get(key)
		constrain: Json Data Unmarshal to Struct
		return constrain.lotSize
	Condition 2: API doesn't provides this information
		return Minimum value*/
func (e *Fcoin) GetLotSize(pair *pair.Pair) float64 {
	key := fmt.Sprintf("%s-Constrain-%s", exchange.FCOIN, pair.Name)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		log.Printf("Fcoin GetLotSize Key: %v Err: %s\n", key, err)
		return 0.00000001
	}
	constrain := exchange.PairConstrain{}
	if str, ok := val.(string); ok {
		if err := json.Unmarshal([]byte(str), &constrain); err != nil {
			log.Printf("Fcoin GetLotSize Key: %v Unmarshal Err: %s\n", key, err)
			return 0.00000001
		}
	} else {
		log.Printf("Fcoin GetLotSize Key: %v can't convert to string: %v", key, val)
		return 0.00000001
	}
	return constrain.LotSize
}

/*Get Pair PriceFilter(Price)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2:
	Condition 1: API provides this information  --Refer Binance Code
		key: Constrain key in Redis ex. key := fmt.Sprintf("%s-Constrain-%s", exchange.<Capital Letter Exchange Name>, pair.Name)
		val: Get Redis Json Data ex. val, err := e.GetMakerDB().Get(key)
		constrain: Json Data Unmarshal to Struct
		return constrain.tickSize
	Condition 2: API doesn't provides this information
		return Minimum value*/
func (e *Fcoin) GetPriceFilter(pair *pair.Pair) float64 { // tickSize for price
	key := fmt.Sprintf("%s-Constrain-%s", exchange.FCOIN, pair.Name)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		log.Printf("Fcoin GetPriceFilter Key: %v Err: %s\n", key, err)
		return 0.00000001
	}
	constrain := exchange.PairConstrain{}
	if str, ok := val.(string); ok {
		if err := json.Unmarshal([]byte(str), &constrain); err != nil {
			log.Printf("Fcoin GetPriceFilter Key: %v Unmarshal Err: %s\n", key, err)
			return 0.00000001
		}
	} else {
		log.Printf("Fcoin GetPriceFilter Key: %v can't convert to string: %v", key, val)
		return 0.00000001
	}
	return constrain.TickSize
	//	return 0.00000001
}

func (e *Fcoin) GetConstrainFetchMethod(pair *pair.Pair) *exchange.ConstrainFetchMethod {
	constrainFetchMethod := &exchange.ConstrainFetchMethod{}
	constrainFetchMethod.Fee = true
	constrainFetchMethod.LotSize = true
	constrainFetchMethod.TickSize = true
	constrainFetchMethod.TxFee = false
	constrainFetchMethod.Withdraw = false
	constrainFetchMethod.Deposit = false
	constrainFetchMethod.Confirmation = false
	return constrainFetchMethod
}

/*************** coins on the exchanges ***************/
/*Get Coin Balance
Step 1: Change Instance Name    (e *<exchange Instance Name>)*/
func (e *Fcoin) GetBalance(coin *coin.Coin) float64 {
	if tmp, ok := balanceMap.Get(coin.Code); ok {
		return tmp.(float64)
	} else {
		return 0.0
	}
}

/*Get Coin Withdraw Fee
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2:
	Condition 1: API provides this information  --Refer Binance Code
		key: Constrain key in Redis ex. key := fmt.Sprintf("%s-Constrain-%s", exchange.<Capital Letter Exchange Name>, coin.Code)
		val: Get Redis Json Data ex. val, err := e.GetMakerDB().Get(key)
		constrain: Json Data Unmarshal to Struct
		return constrain.TxFee
	Condition 2: API doesn't provides this information
		return Minimum value*/
func (e *Fcoin) GetTxFee(coin *coin.Coin) float64 { // Withdraw Fee
	key := fmt.Sprintf("%s-Constrain-%s", exchange.FCOIN, coin.Code)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		log.Printf("Fcoin GetTxFee Key: %v Err: %s\n", key, err)
		return 100.001
	}
	constrain := exchange.CoinConstrain{}
	if str, ok := val.(string); ok {
		if err := json.Unmarshal([]byte(str), &constrain); err != nil {
			log.Printf("Fcoin GetTxFee Key: %v Unmarshal Err: %s\n", key, err)
			return 100.001
		}
	} else {
		log.Printf("Fcoin GetConfirmaGetTxFeetion Key: %v can't convert to string: %v", key, val)
		return 100.001
	}
	return constrain.TxFee
}

/*Get Coin Confirmation
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2:
	Condition 1: API provides this information  --Refer Binance Code
		key: Constrain key in Redis ex. key := fmt.Sprintf("%s-Constrain-%s", exchange.<Capital Letter Exchange Name>, coin.Code)
		val: Get Redis Json Data ex. val, err := e.GetMakerDB().Get(key)
		constrain: Json Data Unmarshal to Struct
		return constrain.Confirmation
	Condition 2: API doesn't provides this information
		return 0*/
func (e *Fcoin) GetConfirmation(coin *coin.Coin) int { // deposit confirmations
	key := fmt.Sprintf("%s-Constrain-%s", exchange.CRYPTOPIA, coin.Code)
	val, err := e.GetMakerDB().Get(key)
	if err != nil {
		log.Printf("Fcoin GetConfirmation Key: %v Err: %s\n", key, err)
		return 1001
	}
	constrain := exchange.CoinConstrain{}
	if str, ok := val.(string); ok {
		if err := json.Unmarshal([]byte(str), &constrain); err != nil {
			log.Printf("Fcoin GetConfirmation Key: %v Unmarshal Err: %s\n", key, err)
			return 1001
		}
	} else {
		log.Printf("Fcoin GetConfirmation Key: %v can't convert to string: %v", key, val)
		return 1001
	}
	return constrain.Confirmation
}

/*Check Coin Withdraw Enable
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2:
	Condition 1: API provides this information  --Refer Binance Code
		key: Constrain key in Redis ex. key := fmt.Sprintf("%s-Constrain-%s", exchange.<Capital Letter Exchange Name>, coin.Code)
		val: Get Redis Json Data ex. val, err := e.GetMakerDB().Get(key)
		constrain: Json Data Unmarshal to Struct
		return constrain.Withdraw
	Condition 2: API doesn't provides this information
		Manually write to Postgres
		When Initial Exchange, read postgres data to be constrain
		-- Detail Ask Chun --*/
func (e *Fcoin) CanWithdraw(coin *coin.Coin) bool { // does withdraw enable
	return true
}

/*Check Coin Deposit Enable
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2:
	Condition 1: API provides this information  --Refer Binance Code
		key: Constrain key in Redis ex. key := fmt.Sprintf("%s-Constrain-%s", exchange.<Capital Letter Exchange Name>, coin.Code)
		val: Get Redis Json Data ex. val, err := e.GetMakerDB().Get(key)
		constrain: Json Data Unmarshal to Struct
		return constrain.Deposit
	Condition 2: API doesn't provides this information
		Manually write to Postgres
		When Initial Exchange, read postgres data to be constrain
		-- Detail Ask Chun --*/
func (e *Fcoin) CanDeposit(coin *coin.Coin) bool { // does deposit enable
	return true
}

/*Get trading website URL
Step 1: Find the website's Exchange page, copy it's URL
Step 2: Change the pair's syntax to match the URL syntax
*/
func (e *Fcoin) GetTradingWebURL(pair *pair.Pair) string {
	return fmt.Sprintf("https://www.cryptopia.co.nz/Exchange/?market=%s_%s", strings.ToUpper(pair.Target.Code), strings.ToUpper(pair.Base.Code))
}
