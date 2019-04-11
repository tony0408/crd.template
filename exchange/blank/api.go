package blank

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"../../coin"
	"../../exchange"
	"../../market"
	"../../pair"
	"../../user"
)

/*The Base Endpoint URL*/
const (
	API_URL string = "API Connection URL"
)

/*API Base Knowledge
Path: API function. Usually after the base endpoint URL
Method:
	Get - Call a URL, API return a response
	Post - Call a URL & send a request, API return a response
Public API:
	It doesn't need authorization/signature , can be called by browser to get response.
	using exchange.HttpGetRequest/exchange.HttpPostRequest
Private API:
	Authorization/Signature is requried. The signature request should look at Exchange API Document.
	using ApiKeyGet/ApiKeyPost
Response:
	Response is a json structure.
	Copy the json to https://transform.now.sh/json-to-go/ convert to go Struct.
	Add the go Struct to model.go

ex. Get /api/v1/depth
Get - Method
/api/v1/depth - Path*/

/*************** Public API ***************/
/*Get Pair Market Depth
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Get Exchange Pair Code ex. symbol := e.GetPairCode(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct*/
func (e *Blank) OrderBook(p *pair.Pair) (*market.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetPairCode(p)

	strRequestUrl := fmt.Sprintf("Orderbook API PATH/%s", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &market.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonBlankOrderbook := exchange.HttpGetRequest(strUrl, nil)
	err := json.Unmarshal([]byte(jsonBlankOrderbook), &orderBook)
	if err != nil {
		return nil, fmt.Errorf("Blank OrderBook json Unmarshal error:%v", err)
	}

	//Convert Exchange Struct to Maker
	maker.Timestamp = orderBook.Timestamp
	for _, bid := range orderBook.Bids {
		var buydata market.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid[0].(string), 64)
		if err != nil {
			return nil, err
		}
		buydata.Quantity, err = strconv.ParseFloat(bid[1].(string), 64)
		if err != nil {
			return nil, err
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata market.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask[0].(string), 64)
		if err != nil {
			return nil, err
		}
		selldata.Quantity, err = strconv.ParseFloat(ask[1].(string), 64)
		if err != nil {
			return nil, err
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*Get Coins Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func GetBlankCoin() *CoinsData {
	coinsInfo := &CoinsData{}

	strRequestUrl := "Currency API PATH"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	json.Unmarshal([]byte(jsonCurrencyReturn), &coinsInfo)

	return coinsInfo
}

/*Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func GetBlankPair() *PairsData {
	pairsInfo := &PairsData{}

	strRequestUrl := "Symbol API PATH"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	json.Unmarshal([]byte(jsonSymbolsReturn), &pairsInfo)

	return pairsInfo
}

/*************** Private API ***************/
func (e *Blank) UpdateAllBalances() {
	e.UpdateAllBalancesByUser(nil)
}

/*Get Exchange Account All Coins Balance  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Call ApiKey Function (Depend on API request)
Step 5: Get Coin Availamount and store in balanceMap*/
func (e *Blank) UpdateAllBalancesByUser(u *user.User) {
	var uInstance *Blank
	if u != nil {
		uInstance = &Blank{}
		uInstance.API_KEY = u.API_KEY
		uInstance.API_SECRET = u.API_SECRET
		// uInstance.MakerDB = e.MakerDB
	} else {
		uInstance = e
	}

	if uInstance.API_KEY == "" || uInstance.API_SECRET == "" {
		log.Printf("Blank API Key or Secret Key are nil.")
		return
	}

	//TODO: GetBalance
}

/*Withdraw the coin to another address
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Call ApiKey Function (Depend on API request)
Step 5: Check the success of withdraw*/
func (e *Blank) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	return false
}

/*Get the Status of a Singal Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (Status reference ../market/market.go)*/
func (e *Blank) OrderStatus(order *market.Order) error {
	return nil
}

func (e *Blank) ListOrders() (*[]market.Order, error) {
	return nil, nil
}

/*Cancel an Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (order.Status = market.Canceling)*/
func (e *Blank) CancelOrder(order *market.Order) error {
	return nil
}

/*Cancel All Order
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (order.Status = market.Canceling)*/
func (e *Blank) CancelAllOrder() error {
	return nil
}

/*Place a limit Sell Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Blank) LimitSell(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	return nil, nil
}

/*Place a limit Buy Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Blank) LimitBuy(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	return nil, nil
}

/*************** Signature Http Request ***************/
/*Method: GET and Signature is required  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Blank) ApiKeyGet(mapParams map[string]string, strRequestPath string) string {
	strMethod := "GET"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	mapParams["AccessKeyId"] = e.API_KEY
	mapParams["Timestamp"] = timestamp

	mapParams["Signature"] = ComputeHmac256(strMethod, e.API_SECRET)

	strUrl := API_URL + strRequestPath
	return exchange.HttpGetRequest(strUrl, mapParams)
}

/*Method: POST and Signature is required  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Blank) ApiKeyPost(mapParams map[string]string, strRequestPath string) string {
	strMethod := "POST"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	//Signature Request Params
	mapParams2Sign := make(map[string]string)
	mapParams2Sign["AccessKeyId"] = e.API_KEY
	mapParams2Sign["Timestamp"] = timestamp

	mapParams2Sign["Signature"] = ComputeHmac256(strMethod, e.API_SECRET)

	strUrl := API_URL + strRequestPath

	return exchange.HttpPostRequest(strUrl, mapParams)
}

//Signature加密
func ComputeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
