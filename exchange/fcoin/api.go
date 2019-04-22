package fcoin

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	//	"io/ioutil"
	"strings"

	"../../coin"
	"../../exchange"
	"../../market"
	"../../pair"
	"../../user"
)

/*The Base Endpoint URL*/
const (
	API_URL string = "https://api.fcoin.com/v2/"
)

/*API Base Knowledge
Path: API function. Usually after the base endpoint URL
Method:
	Get - Call a URL, API return a response
	Post - Call a URL & send a request, API return a response
Public API:
	It doesn't need authorization/signature , can be called by browser to get 0response.
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
func (e *Fcoin) OrderBook(p *pair.Pair) (*market.Maker, error) {

	jsonResponse := JsonResponse{}
	orderBook := OrderBook{}
	symbol := strings.ToLower(e.GetPairCode(p))

	strRequestUrl := fmt.Sprintf("market/depth/L150/%s", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &market.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonFcoinOrderbook := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonFcoinOrderbook), &jsonResponse); err != nil {
		log.Printf("Fcoin Get Coin json Unmarshal error: %v %v", err, jsonFcoinOrderbook)
		return nil, nil
	} else if jsonResponse.Status != 0 {
		log.Printf("Fcoin Get Coin failed:%v Message:%v", jsonResponse.Status, jsonResponse.Message)
		return nil, nil
	}

	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		log.Printf("Fcoin Get Coin Json Unmarshal Err: %v %s", err, jsonResponse.Data)
		return nil, nil
	}

	//Convert Exchange Struct to Maker
	maker.Timestamp = float64(orderBook.Ts)
	var buyRate float64
	for i, bid := range orderBook.Bids {
		var buydata market.Order

		//Modify according to type and structure
		if i%2 == 0 {
			buyRate = bid
		}

		if i%2 == 1 {
			buydata.Rate = buyRate
			buydata.Quantity = bid
			maker.Bids = append(maker.Bids, buydata)
		}
	}
	var sellRate float64
	for i, ask := range orderBook.Asks {
		var selldata market.Order

		//Modify according to type and structure
		if i%2 == 0 {
			sellRate = ask
		}

		if i%2 == 1 {
			selldata.Rate = sellRate
			selldata.Quantity = ask
			maker.Asks = append(maker.Asks, selldata)
		}

	}

	return maker, nil
}

/*Get Coins Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func GetFcoinCoin() []string {

	jsonResponse := JsonResponse{}
	coinsData := []string{}

	strRequestUrl := "public/currencies"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		log.Printf("Fcoin Get Coin json Unmarshal error: %v %v", err, jsonCurrencyReturn)
		return nil
	} else if jsonResponse.Status != 0 {
		log.Printf("Fcoin Get Coin failed, jsonResponse error:%v Message:%v", jsonResponse.Status, jsonResponse.Message)
		return nil
	}

	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		log.Printf("Fcoin Get Coin Json Unmarshal Err: %v %s", err, jsonResponse.Data)
		return nil
	}

	//	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	//	json.Unmarshal([]byte(jsonCurrencyReturn), &coinsInfo)

	return coinsData
}

/*Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func GetFcoinPair() *PairsData {

	jsonResponse := JsonResponse{}
	pairsData := &PairsData{}

	strRequestUrl := "public/symbols"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		log.Printf("Fcoin Get Coin json Unmarshal error: %v %v", err, jsonCurrencyReturn)
		return nil
	} else if jsonResponse.Status != 0 {
		log.Printf("Fcoin Get Coin failed:%v Message:%v", jsonResponse.Status, jsonResponse.Message)
		return nil
	}

	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		log.Printf("Fcoin Get Coin Json Unmarshal Err: %v %s", err, jsonResponse.Data)
		return nil
	}

	//	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	//	json.Unmarshal([]byte(jsonSymbolsReturn), &pairsInfo)

	return pairsData
}

/*************** Private API ***************/
func (e *Fcoin) UpdateAllBalances() {
	e.UpdateAllBalancesByUser(nil)
}

/*Get Exchange Account All Coins Balance  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Call ApiKey Function (Depend on API request)
Step 5: Get Coin Availamount and store in balanceMap*/
func (e *Fcoin) UpdateAllBalancesByUser(u *user.User) {
	var uInstance *Fcoin
	if u != nil {
		uInstance = &Fcoin{}
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

	jsonResponse := JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "accounts/balance"

	jsonBalanceReturn := uInstance.ApiKeyGet(nil, strRequest)
	log.Printf("jsonBalanceReturn: %v", jsonBalanceReturn)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("Fcoin Get Balance Json Unmarshal Err: %v %v", err, jsonBalanceReturn)
		return
	} else if jsonResponse.Status != 0 {
		log.Printf("Fcoin Get Balance Err: %v %v", jsonResponse.Error, jsonResponse.Message)
		return
	}

	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("Fcoin Get Balance Data Unmarshal Err: %v %v", err, jsonResponse.Data)
		return
	} else {
		for _, data := range accountBalance {
			c := coin.GetCoin(e.GetCode(data.Currency))
			if c != nil {
				available, err := strconv.ParseFloat(data.Available, 64)
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
func (e *Fcoin) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("Fcoin API Key or Secret Key are nil.")
		return false
	}

	jsonResponse := JsonResponse{}
	strRequest := "broker/otc/assets/transfer/out"

	mapParams := make(map[string]string)
	mapParams["Currency"] = e.GetSymbol(coin.Code)
	mapParams["Amount"] = fmt.Sprintf("%f", quantity)
	//mapParams["Address"] = addr
	//mapParams["PaymentId"] = coin.Code

	jsonSubmitWithdraw := e.ApiKeyPost(mapParams, strRequest)
	log.Printf("jsonSubmitWithdraw: %+v", jsonSubmitWithdraw)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		log.Printf("Fcoin Withdraw Json Unmarshal failed: %v %v", err, jsonSubmitWithdraw)
		return false
	} else if jsonResponse.Status != 0 {
		log.Printf("Fcoin Withdraw failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
		return false
	}

	return false
}

/*Get the Status of a Singal Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (Status reference ../market/market.go)*/
func (e *Fcoin) OrderStatus(order *market.Order) error {
	//log.Printf("=========OrderStatus order===%+v=========", order) // ============================================
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("Fcoin API Key or Secret Key are nil.")
	}

	jsonResponse := JsonResponse{}
	orderStatus := TradeHistory{}
	strRequest := fmt.Sprintf("orders/%s", order.OrderID)

	mapParams := make(map[string]string)
	//mapParams["Market"] = fmt.Sprintf("%s/%s", e.GetSymbol(order.Pair.Target.Code), e.GetSymbol(order.Pair.Base.Code))

	jsonOrderStatus := e.ApiKeyGet(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("Fcoin OrderStatus Unmarshal Err: %v %v", err, jsonOrderStatus)
	} else if jsonResponse.Status != 0 {
		return fmt.Errorf("Fcoin Get OrderStatus failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
	}

	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("Fcoin Get OrderStatus Data Unmarshal Err: %v %v", err, jsonResponse.Data)
	} else {
		if orderStatus.ID == order.OrderID {
			if orderStatus.Amount == "0" {
				order.Status = market.New
			} else if orderStatus.FilledAmount == orderStatus.Amount {
				order.Status = market.Filled
			} else {
				order.Status = market.Partial
			}
		}

		/*for _, list := range orderStatus {
			orderIDStr := fmt.Sprintf("%d", list.OrderID)
			if orderIDStr == order.OrderID {
				if list.Remaining == 0 {
					order.Status = market.Filled
				} else if list.Remaining == list.Amount {
					order.Status = market.New
				} else {
					order.Status = market.Partial
				}
			}
		}*/
	}

	return nil
}

func (e *Fcoin) ListOrders() (*[]market.Order, error) {
	return nil, nil
}

/*Cancel All Order
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (order.Status = market.Canceling)*/
func (e *Fcoin) CancelAllOrder() error {
	return nil
}

/*Cancel an Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (order.Status = market.Canceling)*/
func (e *Fcoin) CancelOrder(order *market.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("Fcoin API Key or Secret Key are nil.")
	}

	jsonResponse := JsonResponse{}
	strRequest := fmt.Sprintf("orders/%s/submit-cancel", order.OrderID)

	mapParams := make(map[string]string)

	jsonCancelOrder := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("Fcoin CancelOrder Unmarshal Err: %v %v", err, jsonCancelOrder)
	} else if jsonResponse.Status != 0 {
		return fmt.Errorf("Fcoin CancelOrder failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
	}

	order.Status = market.Canceling

	return nil
}

/*Place a limit Sell Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Fcoin) LimitSell(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("Fcoin API Key or Secret Key are nil.")
	}

	jsonResponse := JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "orders"

	mapParams := make(map[string]string)
	mapParams["symbol"] = strings.ToLower(e.GetPairCode(pair))
	mapParams["type"] = "limit"
	mapParams["side"] = "sell"
	mapParams["price"] = fmt.Sprint(rate)
	mapParams["amount"] = fmt.Sprint(quantity)

	jsonPlaceReturn := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("Fcoin LimitSell Unmarshal Err: %v %v", err, jsonPlaceReturn)
	} else if jsonResponse.Status != 0 {
		return nil, fmt.Errorf("Fcoin LimitSell failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
	}

	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("Fcoin LimitSell Data Unmarshal Err: %v %v", err, jsonResponse.Data)
	} else {
		order := &market.Order{}
		order.Pair = pair
		order.Rate = rate
		order.Quantity = quantity
		order.Side = "Sell"
		order.JsonResponse = jsonPlaceReturn
		if placeOrder.Data != "0" {
			order.OrderID = fmt.Sprintf("%d", placeOrder.Data)
			//order.FilledOrders = placeOrder.FilledOrders
			order.Status = market.New
		} /*else if len(placeOrder.FilledOrders) > 0 {
			order.OrderID = "Filled"
			order.FilledOrders = placeOrder.FilledOrders
			order.Status = market.Filled
		}*/

		return order, nil
	}

	return nil, nil
}

/*Place a limit Buy Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Fcoin) LimitBuy(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {

	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("Fcoin API Key or Secret Key are nil.")
	}

	jsonResponse := JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "orders"

	mapParams := make(map[string]string)
	mapParams["symbol"] = strings.ToLower(e.GetPairCode(pair))
	mapParams["type"] = "limit"
	mapParams["side"] = "buy"
	mapParams["price"] = fmt.Sprint(rate)
	mapParams["amount"] = fmt.Sprint(quantity)

	//log.Printf("===amount: %v", mapParams["amount"]) //=========================
	//log.Printf("===price: %v", mapParams["price"])   //=========================

	jsonPlaceReturn := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("Fcoin LimitBuy Unmarshal Err: %v %v", err, jsonPlaceReturn)
	} else if jsonResponse.Status != 0 {
		return nil, fmt.Errorf("Fcoin LimitBuy failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
	}

	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("Fcoin LimitBuy Data Unmarshal Err: %v %v", err, jsonResponse.Data)
	} else {
		order := &market.Order{}
		order.Pair = pair
		order.Rate = rate
		order.Quantity = quantity
		order.Side = "Buy"
		order.JsonResponse = jsonPlaceReturn
		if placeOrder.Data != "0" {
			order.OrderID = placeOrder.Data
			//order.FilledOrders = placeOrder.FilledOrders
			order.Status = market.New
		} /*else if len(placeOrder.FilledOrders) > 0 {
			order.OrderID = "Filled"
			order.FilledOrders = placeOrder.FilledOrders
			order.Status = market.Filled
		}*/

		return order, nil
	}

	return nil, nil
}

/*************** Signature Http Request ***************/
/*Method: GET and Signature is required  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Fcoin) ApiKeyGet(mapParams map[string]string, strRequestPath string) string {
	strMethod := "GET"
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)

	//Signature Request Params
	strUrl := API_URL + strRequestPath

	var strRequestUrl string
	if nil == mapParams {
		strRequestUrl = strUrl
	} else {
		strParams := Map2UrlQuery(mapParams)
		strRequestUrl = strUrl + "?" + strParams
	}

	// signMessage + POST request data
	signMessage := strMethod + strRequestUrl + timestamp
	//log.Printf("signMessage: %s", signMessage) //======================================
	Signature := base64.StdEncoding.EncodeToString([]byte(signMessage))
	Signature1 := ComputeHmac1(Signature, e.API_SECRET)
	// -todo-

	httpClient := &http.Client{}
	request, err := http.NewRequest("GET", strRequestUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("FC-ACCESS-KEY", e.API_KEY)
	request.Header.Add("FC-ACCESS-SIGNATURE", Signature1)
	request.Header.Add("FC-ACCESS-TIMESTAMP", timestamp)

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

/*Method: POST and Signature is required  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Fcoin) ApiKeyPost(mapParams map[string]string, strRequestPath string) string {
	strMethod := "POST"
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10) //time.Now().UTC().Format("2006-01-02T15:04:05")

	//Signature Request Params
	strUrl := API_URL + strRequestPath

	jsonParams := ""
	bytesParams, err := json.Marshal(mapParams)
	if nil != mapParams {
		jsonParams = string(bytesParams)
	}

	var strParams, strRequestUrl string
	if nil == mapParams {
		strRequestUrl = strUrl
	} else {
		strParams = Map2UrlQuery(mapParams)
		strRequestUrl = strUrl //+ "?" + strParams
	}

	// signMessage + POST request data
	signMessage := strMethod + strRequestUrl + timestamp + strParams //jsonParams
	//log.Printf("signMessage: %s", signMessage)                       //====================================
	Signature := base64.StdEncoding.EncodeToString([]byte(signMessage))
	Signature1 := ComputeHmac1(Signature, e.API_SECRET)

	// -todo-

	httpClient := &http.Client{}
	request, err := http.NewRequest("POST", strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("FC-ACCESS-KEY", e.API_KEY)
	request.Header.Add("FC-ACCESS-SIGNATURE", Signature1)
	request.Header.Add("FC-ACCESS-TIMESTAMP", timestamp)
	request.Header.Add("content-type", "application/json;charset=UTF-8")

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
}

//Signature加密
func ComputeHmac1(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 将map格式的请求参数转换为字符串格式的
// mapParams: map格式的参数键值对
// return: 查询字符串
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
