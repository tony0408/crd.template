package exchange

import (
	"../coin"
	"../pair"
)

type Config struct {
	RedisServer  string
	RedisDB      int
	Account_ID   string
	API_KEY      string
	API_SECRET   string
	WalletStatus []Wallet_Stat
}

type PairConstrain struct {
	Pair     *pair.Pair //the code on excahnge with the same chain, eg: BCH, BCC on different exchange, but they are the same chain
	LotSize  float64    // the decimal place for this coin on exchange for the pairs, eg:  BTC: 0.00001    NEO:1   LTC: 0.001 ETH:0.01
	TickSize float64
	Issue    string //the issue for the pair if have any problem
}

type CoinConstrain struct {
	Coin         *coin.Coin
	TxFee        float64 // the withdraw fee for this exchange
	Withdraw     bool
	Deposit      bool
	Confirmation int
	Issue        string //the issue for the chain if have any problem
}

type ConstrainFetchMethod struct {
	Fee          bool
	LotSize      bool
	TickSize     bool
	TxFee        bool
	Withdraw     bool
	Deposit      bool
	Confirmation bool
}

type Wallet_Stat struct {
	Currency string
	Withdraw bool
	Deposit  bool
	TxFee    float64
}
