package cryptopia

import (
	"strings"

	"../../coin"
	"../../exchange"
	"../../pair"
)

func (e *Cryptopia) UpdatePairConstrain() {
	pairData := GetCryptopiaPair()
	pairConstrainMap := make(map[*pair.Pair]*exchange.PairConstrain)
	//If Exchange doesn't provide constrain info, Leave blank
	//Modify according to type and structure
	for _, symbol := range *pairData {
		pairConstrain := &exchange.PairConstrain{}

		base := coin.GetCoin(e.GetCode(symbol.BaseSymbol))
		target := coin.GetCoin(e.GetCode(symbol.Symbol))

		pairConstrain.Pair = pair.GetPair(base, target)

		pairConstrain.LotSize = symbol.MinimumTrade
		pairConstrain.TickSize = symbol.MinimumPrice

		pairConstrainMap[pairConstrain.Pair] = pairConstrain

	}
	// return pairConstrainMap
}

func (e *Cryptopia) UpdateCoinConstrain() {
	coinInfo := GetCryptopiaCoin()
	coinConstrainMap := make(map[*coin.Coin]*exchange.CoinConstrain)
	//If Exchange doesn't provide constrain info, Leave cryptopia
	//Modify according to type and structure
	for _, data := range *coinInfo {
		coinConstrain := &exchange.CoinConstrain{}
		coinConstrain.Coin = coin.GetCoin(e.GetCode(data.Symbol))
		coinConstrain.TxFee = data.WithdrawFee
		if data.Status == "OK" {
			coinConstrain.Withdraw = true
			coinConstrain.Deposit = true
		} else {
			coinConstrain.Withdraw = false
			coinConstrain.Deposit = false
		}
		coinConstrain.Confirmation = data.DepositConfirmations

		coinConstrainMap[coinConstrain.Coin] = coinConstrain

	}
	// return coinConstrainMap
}

/***************************************************/
var symbolMap = make(map[string]string)

func (e *Cryptopia) FixSymbol() { //key: exchange specific    val： bitontop standard
	symbolMap["ACC"] = "ACC1"
	symbolMap["BCS"] = "BCS1"
	symbolMap["BITS"] = "BITS1"
	symbolMap["CAP"] = "CAP1"
	symbolMap["CAT"] = "CAT1"
	symbolMap["CMT"] = "CMT1"
	symbolMap["HAV"] = "HAV1"
	symbolMap["HC"] = "HC1"
	symbolMap["IQ"] = "IQ1"
	symbolMap["LDC"] = "LDC1"
	symbolMap["QBT"] = "QBT1"
	symbolMap["VCC"] = "VCC1"
}

func (e *Cryptopia) GetSymbol(code string) string { // get exchange standard
	code = strings.ToUpper(code)
	for k, v := range symbolMap {
		if code == v {
			return k
		}
	}
	// log.Printf("GetSymbol error!")
	return code
}

func (e *Cryptopia) GetCode(symbol string) string { // get bitontop standard
	symbol = strings.ToUpper(symbol)
	if val, ok := symbolMap[symbol]; ok {
		return val
	}
	return symbol
}
