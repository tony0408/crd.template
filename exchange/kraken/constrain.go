package kraken

import (
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
func (e *Kraken) UpdatePairConstrain() {//map[*pair.Pair]*exchange.PairConstrain {
	pairData := GetKrakenPair()
	pairConstrainMap := make(map[*pair.Pair]*exchange.PairConstrain)
	//If Exchange doesn't provide constrain info, Leave kraken
	//Modify according to type and structure
	for _, symbol := range pairData {
		pairConstrain := &exchange.PairConstrain{}

		base := coin.GetCoin(e.GetCode(symbol.Quote))
		target := coin.GetCoin(e.GetCode(symbol.Base))

		pairConstrain.Pair = pair.GetPair(base, target)

		pairConstrain.LotSize = math.Pow10(symbol.LotDecimals * -1)
		pairConstrain.TickSize = math.Pow10(symbol.PairDecimals * -1)
		pairConstrainMap[pairConstrain.Pair] = pairConstrain
		//l, err := json.Marshal(pairConstrain)
		// if err != nil {
		// 	log.Printf("Kraken UpdatePairConstrain Marshal err: %s\n", err)
		// }
		// if pairConstrain.Pair.Name != "" {
		// 	key := fmt.Sprintf("%s-Constrain-%s", exchange.KRAKEN, pairConstrain.Pair.Name)
		// 	err = e.GetMakerDB().Set(key, string(l))
		// 	if err != nil {
		// 		log.Printf("Kraken UpdatePairConstrain Set DB err: %s\n", err)
		// 	}
		// }
	}
	//return pairConstrainMap
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
func (e *Kraken) UpdateCoinConstrain() {//map[*coin.Coin]*exchange.CoinConstrain {
	//return nil
}

/***************************************************/
var symbolMap = make(map[string]string)

/*Standard Coin Code
Coin has same code but it is different currency
Fix the coin code to bitontop standard*/
func (e *Kraken) FixSymbol() { //key: exchange specific    val： bitontop standard
	symbolMap["XXBT"] = "BTC"
	symbolMap["XDAO"] = "DAO"
	symbolMap["XETC"] = "ETC"
	symbolMap["XETH"] = "ETH"
	symbolMap["XICN"] = "ICN"
	symbolMap["XLTC"] = "LTC"
	symbolMap["XMLN"] = "MLN"
	symbolMap["XNMC"] = "NMC"
	symbolMap["XREP"] = "REP"
	symbolMap["XXDG"] = "XDG"
	symbolMap["XXLM"] = "XLM"
	symbolMap["XXMR"] = "XMR"
	symbolMap["XXRP"] = "XRP"
	symbolMap["XXVN"] = "XVN"
	symbolMap["XZEC"] = "ZEC"
	symbolMap["ZCAD"] = "CAD"
	symbolMap["ZEUR"] = "EUR"
	symbolMap["ZGBP"] = "GBP"
	symbolMap["ZJPY"] = "JPY"
	symbolMap["ZKRW"] = "KRW"
	symbolMap["ZUSD"] = "USD"
}

/*Get Exchange Standard Code*/
func (e *Kraken) GetSymbol(code string) string {
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
func (e *Kraken) GetCode(symbol string) string {
	symbol = strings.ToUpper(symbol)
	if val, ok := symbolMap[symbol]; ok {
		return val
	}
	return symbol
}
