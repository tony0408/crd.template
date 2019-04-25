package bitrue

import (
	"strings"
)

/*Update Pairs Constrain  --If API provide those information
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Change Exchange Name    exchange.<Capital Letter Exchange Name>
Step 3: Get Pairs Data from API
Step 4: Get Each Symbol
Step 5: Identify Base & Target and Get Pair
Step 6: Add LotSize  - float64
Step 7: Add TickSize  - float64*/
func (e *Bitrue) UpdatePairConstrain() {
	// pairData := GetBitrueCoin()

	// //If Exchange doesn't provide constrain info, Leave bitrue
	// //Modify according to type and structure
	// for _, symbol := range pairData.Symbols {
	// 	pairConstrain := &exchange.PairConstrain{}

	// 	base := coin.GetCoin(e.GetCode(symbol.QuoteAsset))
	// 	target := coin.GetCoin(e.GetCode(symbol.BaseAsset))

	// 	pairConstrain.Pair = pair.GetPair(base, target)

	// 	lotsize, err := strconv.ParseFloat(symbol.LotSize, 64)
	// 	if err != nil {
	// 		log.Printf("Bitrue Lot_Size Err: %s\n", err)
	// 	}
	// 	pairConstrain.LotSize = lotsize

	// 	ticksize, err := strconv.ParseFloat(symbol.TickSize, 64)
	// 	if err != nil {
	// 		log.Printf("Bitrue Tick_Size Err: %s\n", err)
	// 	}
	// 	pairConstrain.TickSize = ticksize

	// 	l, err := json.Marshal(pairConstrain)
	// 	if err != nil {
	// 		log.Printf("Bitrue UpdatePairConstrain Marshal err: %s\n", err)
	// 	}
	// 	if pairConstrain.Pair.Name != "" {
	// 		key := fmt.Sprintf("%s-Constrain-%s", exchange.BITRUE, pairConstrain.Pair.Name)
	// 		err = e.GetMakerDB().Set(key, string(l))
	// 		if err != nil {
	// 			log.Printf("Bitrue UpdatePairConstrain Set DB err: %s\n", err)
	// 		}
	// 	}
	// }
	return
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
func (e *Bitrue) UpdateCoinConstrain() {
	// coinInfo := GetBitrueCoin()

	// //If Exchange doesn't provide constrain info, Leave bitrue
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
	// 		log.Printf("Bitrue UpdateCoinConstrain Marshal err: %s\n", err)
	// 	}
	// 	if coinConstrain.Coin != nil {
	// 		key := fmt.Sprintf("%s-Constrain-%s", exchange.BITRUE, coinConstrain.Coin.Code)
	// 		err = e.GetMakerDB().Set(key, string(l))
	// 		if err != nil {
	// 			log.Printf("Bitrue UpdateCoinConstrain Set DB err: %s\n", err)
	// 		}
	// 	}
	// }
	return
}

/***************************************************/
var symbolMap = make(map[string]string)

/*Standard Coin Code
Coin has same code but it is different currency
Fix the coin code to bitontop standard*/
func (e *Bitrue) FixSymbol() { //key: exchange specific    val： bitontop standard
	symbolMap["BCHSV"] = "BSV"
}

/*Get Exchange Standard Code*/
func (e *Bitrue) GetSymbol(code string) string {
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
func (e *Bitrue) GetCode(symbol string) string {
	symbol = strings.ToUpper(symbol)
	if val, ok := symbolMap[symbol]; ok {
		return val
	}
	return symbol
}
