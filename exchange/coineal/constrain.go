package coineal

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"

	"../../coin"
	"../../exchange"
	"../../pair"
)

/*Update Pairs Constrain  --If API provide those information
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Change Exchange Name    exchange.<Capital Letter Exchange Name>
Step 3: Get Pairs Data from API
Step 4: Get Each Symbol
Step 5: Identify Base & Target and Get Pair
Step 6: Add LotSize  - float64
Step 7: Add TickSize  - float64*/
func (e *Coineal) UpdatePairConstrain() {
	pairData := GetCoinealPair()

	//If Exchange doesn't provide constrain info, Leave coineal
	//Modify according to type and structure

	for _, symbol := range pairData.Data {
		pairConstrain := &exchange.PairConstrain{}
		target := coin.GetCoin(e.GetCode(symbol.BaseCoin))
		base := coin.GetCoin(e.GetCode(symbol.CountCoin))
		pairConstrain.Pair = pair.GetPair(base, target)
		pairConstrain.LotSize = math.Pow10(-1 * symbol.AmountPrecision)
		pairConstrain.TickSize = math.Pow10(-1 * symbol.PricePrecision)

		l, err := json.Marshal(pairConstrain)
		if err != nil {
			log.Printf("Coineal UpdatePairConstrain Marshal err: %s\n", err)
		}
		if pairConstrain.Pair.Name != "" {
			key := fmt.Sprintf("%s-Constrain-%s", exchange.COINEAL, pairConstrain.Pair.Name)
			err = e.GetMakerDB().Set(key, string(l))
			if err != nil {
				log.Printf("Coineal UpdatePairConstrain Set DB err: %s\n", err)
			}
		}
	}
}

/*Update Coins Constrain  --If API provide those information
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Change Exchange Name    exchange.<Capital Letter Exchange Name>
Step 3: Get Coins Data from API
Step 4: Get Each Coin
Step 5: Get the coin (Use Standard Code ex. e.GetCode(coin))
Step 6: Add TxFee - float64
Step 7: Add Withdraw Status - Bool
Step 7: Add Deposite Status - Bool
Step 7: Add Confirmation - Int*/
func (e *Coineal) UpdateCoinConstrain() {
	// coinInfo := GetCoinealCoins()

	// //If Exchange doesn't provide constrain info, Leave coineal
	// //Modify according to type and structure
	// for _, data := range coinInfo {
	// 	coinConstrain := &exchange.CoinConstrain{}
	// 	coinConstrain.Coin = coin.GetCoin(e.GetCode(data.ID))
	// 	coinConstrain.TxFee, _ = strconv.ParseFloat(data.WithdrawFee, 64)
	// 	coinConstrain.Withdraw = data.WithdrawStatus
	// 	coinConstrain.Deposit = data.DepositStatus
	// 	coinConstrain.Confirmation = data.DepositConfirmation
	// 	l, err := json.Marshal(coinConstrain)
	// 	if err != nil {
	// 		log.Printf("Coineal UpdateCoinConstrain Marshal err: %s\n", err)
	// 	}
	// 	if coinConstrain.Coin != nil {
	// 		key := fmt.Sprintf("%s-Constrain-%s", exchange.COINEAL, coinConstrain.Coin.Code)
	// 		err = e.GetMakerDB().Set(key, string(l))
	// 		if err != nil {
	// 			log.Printf("Coineal UpdateCoinConstrain Set DB err: %s\n", err)
	// 		}
	// 	}
	// }
}

/***************************************************/
var symbolMap = make(map[string]string)

/*Standard Coin Code
Coin has same code but it is different currency
Fix the coin code to bitontop standard*/
func (e *Coineal) FixSymbol() { //key: exchange specific    val： bitontop standard
	symbolMap["-"] = ""
}

/*Get Exchange Standard Code*/
func (e *Coineal) GetSymbol(code string) string {
	code = strings.ToUpper(code)
	for k, v := range symbolMap {
		if code == v {
			return k
		}
	}
	// log.Printf("GetSymbol error!")
	return code
}

/*Get Bitontop Standard Code*/
func (e *Coineal) GetCode(symbol string) string {
	symbol = strings.ToUpper(symbol)
	if val, ok := symbolMap[symbol]; ok {
		return val
	}
	return symbol
}
