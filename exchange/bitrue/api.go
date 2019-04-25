package bitrue

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	//"time"

	"../../coin"
	"../../exchange"
	"../../market"
	"../../pair"
	"../../user"
)

/*The Base Endpoint URL*/
const (
	API_URL string = "https://www.bitrue.com"
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
func (e *Bitrue) OrderBook(p *pair.Pair) (*market.Maker, error) {
	orderBook := BitrueOrderBook{}
	symbol := e.GetPairCode(p)

	strRequestUrl := "/api/v1/depth"
	strUrl := API_URL + strRequestUrl
	maker := &market.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)
	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["limit"] = "0"
	jsonBitrueOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	err := json.Unmarshal([]byte(jsonBitrueOrderbook), &orderBook)
	if err != nil {
		return nil, fmt.Errorf("Bitrue OrderBook json Unmarshal error:%v", err)
	}

	//Convert Exchange Struct to Maker
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
func GetBitrueCoin() BitruePair {
	coinsInfo := BitruePair{}

	strRequestUrl := "/api/v1/exchangeInfo"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	json.Unmarshal([]byte(jsonCurrencyReturn), &coinsInfo)

	return coinsInfo
}

/*Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
// func GetBitruePair() PairsData {
// 	pairsInfo := PairsData{}

// 	strRequestUrl := "Symbol API PATH"
// 	strUrl := API_URL + strRequestUrl

// 	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
// 	json.Unmarshal([]byte(jsonSymbolsReturn), &pairsInfo)

// 	return pairsInfo
// }

/*************** Private API ***************/
func (e *Bitrue) UpdateAllBalances() {
	e.UpdateAllBalancesByUser(nil)
}

/*Get Exchange Account All Coins Balance  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Call ApiKey Function (Depend on API request)
Step 5: Get Coin Availamount and store in balanceMap*/
func (e *Bitrue) UpdateAllBalancesByUser(u *user.User) {
	var uInstance *Bitrue
	if u != nil {
		uInstance = &Bitrue{}
		uInstance.API_KEY = u.API_KEY
		uInstance.API_SECRET = u.API_SECRET
	} else {
		uInstance = e
	}

	if uInstance.API_KEY == "" || uInstance.API_SECRET == "" {
		log.Printf("Bitrue API Key or Secret Key are nil.")
		return
	}

	accountBalance := AccountBalances{}
	strRequest := "/api/v1/account"

	jsonBalanceReturn := uInstance.ApiKeyRequest("GET", make(map[string]string), strRequest)

	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("Bitrue Get Balance Json Unmarshal Err: %v %v", err, jsonBalanceReturn)
		return
	} else {
		for _, data := range accountBalance.Balances {
			c := coin.GetCoin(e.GetCode(data.Asset))
			if c != nil {
				available, err := strconv.ParseFloat(data.Free, 64)
				if err == nil {
					balanceMap.Set(c.Code, available)
				}
			}
		}
	}

	//TODO: GetBalance
}

/*Withdraw the coin to another address
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Call ApiKey Function (Depend on API request)
Step 5: Check the success of withdraw*/
func (e *Bitrue) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	return false
}

/*Get the Status of a Singal Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (Status reference ../market/market.go)*/
func (e *Bitrue) OrderStatus(order *market.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("Bitrue API Key or Secret Key are nil.")
	}

	orderStatus := TradeHistory{}
	strRequest := "/api/v1/order"

	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)

	mapParams := make(map[string]string)
	mapParams["method"] = "GET"
	mapParams["symbol"] = strings.ToUpper(e.GetPairCode(order.Pair))
	mapParams["orderId"] = order.OrderID
	mapParams["timestamp"] = timestamp

	jsonOrderStatus := e.ApiKeyRequest("GET", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("Bitrue OrderStatus Unmarshal Err: %v %v", err, jsonOrderStatus)
	} else {
		if strconv.Itoa(orderStatus.OrderID) == order.OrderID {
			switch orderStatus.Status {
			case "NEW":
				order.Status = market.New
			case "PARTIALLY_FILLED":
				order.Status = market.Partial
			case "REJECTED":
				order.Status = market.Rejected
			case "PENDING_CANCEL":
				order.Status = market.Canceling
			case "CANCELED":
				order.Status = market.Canceled
			case "FILLED":
				order.Status = market.Filled
			case "EXPIRED":
				order.Status = market.Expired
			default:
				order.Status = market.Other
			}
		}
	}

	return nil
}

func (e *Bitrue) ListOrders() (*[]market.Order, error) {
	return nil, nil
}

/*Cancel an Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (order.Status = market.Canceling)*/
func (e *Bitrue) CancelOrder(order *market.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("Bitrue API Key or Secret Key are nil.")
	}

	//jsonResponse := JsonResponse{}
	strRequest := "/api/v1/order"
	cancelOrder := CancelOrder{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strings.ToUpper(e.GetPairCode(order.Pair))

	jsonCancelOrder := e.ApiKeyRequest("DELETE", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("Bitrue CancelOrder Unmarshal Err: %v %v", err, jsonCancelOrder)
	} else if strconv.Itoa(cancelOrder.OrderID) != order.OrderID {
		return fmt.Errorf("Bitrue CancelOrder failed:%+v Message:%v", cancelOrder)
	}

	order.Status = market.Canceling

	return nil
}

/*Cancel All Order
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (order.Status = market.Canceling)*/
func (e *Bitrue) CancelAllOrder() error {
	return nil
}

/*Place a limit Sell Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Bitrue) LimitSell(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("Bitrue API Key or Secret Key are nil.")
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api/v1/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = strings.ToUpper(e.GetPairCode(pair))
	mapParams["side"] = "SELL"
	mapParams["type"] = "LIMIT"
	mapParams["price"] = fmt.Sprint(rate)
	mapParams["quantity"] = fmt.Sprint(quantity)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("Bitrue LimitSell Unmarshal Err: %v %v", err, jsonPlaceReturn)
	} else {
		order := &market.Order{
			OrderID:      strconv.Itoa(placeOrder.OrderID),
			Pair:         pair,
			Rate:         rate,
			Quantity:     quantity,
			Side:         "Sell",
			Status:       market.New,
			JsonResponse: jsonPlaceReturn,
		}
		return order, nil
	}

	return nil, nil
}

/*Place a limit Buy Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Bitrue) LimitBuy(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {

	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("Bitrue API Key or Secret Key are nil.")
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api/v1/order"

	mapParams := make(map[string]string)
	mapParams["method"] = "POST"
	mapParams["symbol"] = strings.ToUpper(e.GetPairCode(pair))
	mapParams["side"] = "BUY"
	mapParams["type"] = "LIMIT"
	mapParams["price"] = fmt.Sprint(rate)
	mapParams["quantity"] = fmt.Sprint(quantity)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("Bitrue LimitBuy Unmarshal Err: %v %v", err, jsonPlaceReturn)
	} else {
		order := &market.Order{
			OrderID:      strconv.Itoa(placeOrder.OrderID),
			Pair:         pair,
			Rate:         rate,
			Quantity:     quantity,
			Side:         "Buy",
			Status:       market.New,
			JsonResponse: jsonPlaceReturn,
		}
		return order, nil
	}

	return nil, nil
}

/*************** Signature Http Request ***************/
/*Method: GET and Signature is required  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bitrue) ApiKeyGet(mapParams map[string]string, strRequestPath string) string {
	strMethod := mapParams["method"]
	delete(mapParams, "method")

	strUrl := API_URL + strRequestPath

	var strParams string
	if nil != mapParams {
		strParams = Map2UrlQuery(mapParams)
	}

	signature := ComputeHmac256(strParams, e.API_SECRET)
	signMessage := strUrl + "?" + strParams + "&signature=" + signature

	httpClient := &http.Client{}
	request, err := http.NewRequest(strMethod, signMessage, nil) //, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("X-MBX-APIKEY", e.API_KEY)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
	//return exchange.HttpGetRequest(strUrl, mapParams)
}

/*Method: POST and Signature is required  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bitrue) ApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string) string {
	mapParams["timestamp"] = fmt.Sprintf("%.0d", time.Now().UnixNano()/1e6)

	strUrl := API_URL + strRequestPath

	var strParams string
	if nil != mapParams {
		strParams = Map2UrlQuery(mapParams)
	}

	signature := ComputeHmac256(strParams, e.API_SECRET)
	signMessage := strUrl + "?" + strParams + "&signature=" + signature

	httpClient := &http.Client{}
	request, err := http.NewRequest(strMethod, signMessage, nil) //, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}

	request.Header.Add("X-MBX-APIKEY", e.API_KEY)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)

	//return exchange.HttpPostRequest(strUrl, mapParams)
}

//Signature加密
func ComputeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return hex.EncodeToString(h.Sum(nil))
	//return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func Map2UrlQuery(mapParams map[string]string) string {
	var strParams string
	keySort := []string{}
	for key := range mapParams {
		keySort = append(keySort, key)
	}

	sort.Strings(keySort)
	for _, key := range keySort {
		strParams += (key + "=" + mapParams[key] + "&")
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}
