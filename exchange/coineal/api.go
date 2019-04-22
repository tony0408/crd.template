package coineal

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"../../coin"
	"../../exchange"
	"../../market"
	"../../pair"
	"../../user"
)

/*The Base Endpoint URL*/
const (
	API_URL string = "https://exchange-open-api.coineal.com"
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
func (e *Coineal) OrderBook(p *pair.Pair) (*market.Maker, error) {
	coinealOrderbook := CoinealOrderBook{}
	symbol := e.GetPairCode(p)

	strRequestUrl := "/open/api/market_dept"
	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["type"] = "step0"
	strUrl := API_URL + strRequestUrl

	jsonCoinealOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	err := json.Unmarshal([]byte(jsonCoinealOrderbook), &coinealOrderbook)
	if err != nil {
		return nil, fmt.Errorf("Coineal OrderBook json Unmarshal error:%v", err)
	}

	//Convert Exchange Struct to Maker
	maker := &market.Maker{}
	maker.Timestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range coinealOrderbook.Data.Tick.Bids {
		var buydata market.Order
		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid[0].(string), 64) //price
		if err != nil {
			return nil, err
		}
		//buydata.Quantity, err = strconv.ParseFloat(bid[1].(string), 64) //amount
		buydata.Quantity = bid[1].(float64)

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range coinealOrderbook.Data.Tick.Asks {
		var selldata market.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask[0].(string), 64)
		if err != nil {
			return nil, err
		}
		selldata.Quantity = ask[1].(float64)

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*Get Coins Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func GetCoinealCoins() CoinealCoins {
	coinealCoins := CoinealCoins{}

	strRequestUrl := "/open/api/common/symbols"
	strUrl := API_URL + strRequestUrl

	jsoncoinealCoins := exchange.HttpGetRequest(strUrl, nil)
	json.Unmarshal([]byte(jsoncoinealCoins), &coinealCoins)

	return coinealCoins
}

/*Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func GetCoinealPair() CoinealCoins {
	pairsInfo := CoinealCoins{}

	strRequestUrl := "/open/api/common/symbols"
	strUrl := API_URL + strRequestUrl

	jsonPairsInfo := exchange.HttpGetRequest(strUrl, nil)
	json.Unmarshal([]byte(jsonPairsInfo), &pairsInfo)

	return pairsInfo
}

/*************** Private API ***************/
func (e *Coineal) UpdateAllBalances() {
	e.UpdateAllBalancesByUser(nil)
}

/*Get Exchange Account All Coins Balance  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Call ApiKey Function (Depend on API request)
Step 5: Get Coin Availamount and store in balanceMap*/
func (e *Coineal) UpdateAllBalancesByUser(u *user.User) {
	var uInstance *Coineal
	if u != nil {
		uInstance = &Coineal{}
		uInstance.API_KEY = u.API_KEY
		uInstance.API_SECRET = u.API_SECRET
		// uInstance.MakerDB = e.MakerDB
	} else {
		uInstance = e
	}

	if uInstance.API_KEY == "" || uInstance.API_SECRET == "" {
		log.Printf("Coineal API Key or Secret Key are nil.")
		return
	}
	strRequestUrl := "/open/api/user/account"
	jsonResponse := JsonResponse{}
	coinealBalance := CoinealBalance{}
	jsonCoinealBalance := e.ApiKeyRequest("GET", make(map[string]string), strRequestUrl)
	//log.Printf("jsonCoinealBalance :%v", jsonCoinealBalance)
	err := json.Unmarshal([]byte(jsonCoinealBalance), &jsonResponse)
	if err != nil {
		log.Printf("Coineal get Balance jsonUnmarshal Error: %v", jsonCoinealBalance)
	} else if jsonResponse.Code != "0" {
		log.Printf("Coineal get balance Error :%v", jsonResponse.Msg)
	} else {
		err := json.Unmarshal([]byte(jsonResponse.Data), &coinealBalance)
		if err != nil {
			log.Printf("Coineal get Balance jsonUnmarshal Error :%v", jsonCoinealBalance)
		}
		for _, balance := range coinealBalance.CoinList {
			freeAmount, err := strconv.ParseFloat(fmt.Sprintf("%v", balance.Normal), 64)
			if err == nil {
				c := coin.GetCoin(e.GetCode(fmt.Sprintf("%v", balance.Coin)))
				if c != nil {
					balanceMap.Set(c.Code, freeAmount)
				} else {
					c = &coin.Coin{}
					c.Code = e.GetCode(fmt.Sprintf("%v", balance.Coin))
					coin.AddCoin(c)
					balanceMap.Set(c.Code, freeAmount)

				}
			}
		}
	}

}

/*Withdraw the coin to another address
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Call ApiKey Function (Depend on API request)
Step 5: Check the success of withdraw*/
func (e *Coineal) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	return false
}

/*Get the Status of a Singal Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (Status reference ../market/market.go)*/
func (e *Coineal) OrderStatus(order *market.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("coineal API Key or Secret Key are nil")
	}

	strRequestUrl := "/open/api/order_info"

	mapParams := make(map[string]string)
	mapParams["order_id"] = order.OrderID
	mapParams["symbol"] = e.GetPairCode(order.Pair)
	jsonResponse := JsonResponse{}
	coinealOrderStatus := CoinealOrderStatus{}

	jsonOrderStatus := e.ApiKeyRequest("GET", mapParams, strRequestUrl)
	//log.Printf("orderstatus :%v", jsonOrderStatus)
	err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse)
	if err != nil {
		return fmt.Errorf("coineal OrderStatus Json Unmarshal Err: %v, %v", err, jsonOrderStatus)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("coineal OrderStatus Failed: %v", jsonResponse.Msg)
	} else {
		err := json.Unmarshal([]byte(jsonResponse.Data), &coinealOrderStatus)
		if err != nil {
			return fmt.Errorf("coineal OrderStatus Json Unmarshal Err: %v, %v", err, jsonOrderStatus)
		}
		status := coinealOrderStatus.Data.OrderInfo.Status
		if status == 4 {
			order.Status = market.Canceled
		} else if status == 5 {
			order.Status = market.Canceling
		} else if status == 6 {
			order.Status = market.Other
		} else if status == 2 {
			order.Status = market.Filled
		} else if status == 3 {
			order.Status = market.Partial
		} else if status == 1 || status == 0 {
			order.Status = market.New
		} else {
			order.Status = market.Other
		}
	}
	return nil
}

func (e *Coineal) ListOrders() ([]*market.Order, error) {
	return nil, nil
}

/*Cancel an Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (order.Status = market.Canceling)*/
func (e *Coineal) CancelOrder(order *market.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("coineal API Key or Secret Key are nil.")
	}
	jsonResponse := JsonResponse{}
	//cancelOrder := CoinealOrder{}

	strRequestUrl := "/open/api/cancel_order"
	mapParams := make(map[string]string)
	mapParams["order_id"] = order.OrderID
	mapParams["symbol"] = e.GetPairCode(order.Pair)
	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequestUrl)
	//log.Printf("cancel order :%v", jsonPlaceReturn)
	err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse)
	if err != nil {
		return fmt.Errorf("coineal cancel order Json Unmarshal Err: %v", jsonPlaceReturn)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("coineal cancel order failed: %v", jsonResponse.Msg)
	} else {
		// err := json.Unmarshal([]byte(jsonResponse.Data), &cancelOrder)
		// if err != nil {
		// 	log.Printf("coineal cancel order jsonUnmarshal Error: %v", jsonPlaceReturn)
		// }
		order.Status = market.Canceling
		return nil
	}
}

/*Cancel All Order
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (order.Status = market.Canceling)*/
func (e *Coineal) CancelAllOrder() error {
	return nil
}

/*Place a limit Sell Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Coineal) LimitSell(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("coineal API Key or Secret Key are nil.")
	}
	jsonResponse := JsonResponse{}
	coinLimitSell := CoinealOrder{}
	strRequestUrl := "/open/api/create_order"
	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetPairCode(pair)
	mapParams["side"] = "SELL"
	mapParams["volume"] = fmt.Sprintf("%v", quantity)
	mapParams["price"] = fmt.Sprintf("%v", rate)
	mapParams["type"] = "1" //now only limit

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequestUrl)
	//log.Printf("jsonPlaceReturn :%v", jsonPlaceReturn)
	err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("coineal LimitSell Json Unmarshal Err: %v", jsonPlaceReturn)
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("coineal LimitSell failed: %v", jsonResponse.Msg)
	} else {
		err := json.Unmarshal([]byte(jsonResponse.Data), &coinLimitSell)
		if err != nil {
			return nil, fmt.Errorf("coineal LimitSell Json Unmarshal Error :%v", jsonPlaceReturn)
		}
		order := &market.Order{
			Pair:     pair,
			OrderID:  fmt.Sprintf("%v", coinLimitSell.OrderID),
			Rate:     rate,
			Quantity: quantity,
			Side:     "SELL",
			Status:   market.New,
		}
		return order, nil
	}
}

/*Place a limit Buy Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Coineal) LimitBuy(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("coineal API Key or Secret Key are nil.")
	}
	jsonResponse := JsonResponse{}
	coinLimitBuy := CoinealOrder{}
	strRequestUrl := "/open/api/create_order"
	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetPairCode(pair)
	mapParams["side"] = "BUY"
	mapParams["volume"] = fmt.Sprintf("%v", quantity)
	mapParams["price"] = fmt.Sprintf("%v", rate)
	mapParams["type"] = "1" //now only limit

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequestUrl)
	//log.Printf("limit buy :%v", jsonPlaceReturn)
	err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("coineal LimitBuy Json Unmarshal Err: %v", jsonPlaceReturn)
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("coineal LimitBuy failed: %v", jsonResponse.Msg)
	} else {
		err := json.Unmarshal([]byte(jsonResponse.Data), &coinLimitBuy)
		if err != nil {
			return nil, fmt.Errorf("coineal LimitBuy Json Unmarshal Error :%v", jsonPlaceReturn)
		}
		order := &market.Order{
			Pair:     pair,
			OrderID:  fmt.Sprintf("%v", coinLimitBuy.OrderID),
			Rate:     rate,
			Quantity: quantity,
			Side:     "BUY",
			Status:   market.New,
		}
		return order, nil
	}
}

/*************** Signature Http Request ***************/
/*Method: GET and Signature is required  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
// func (e *Coineal) ApiKeyGet(mapParams map[string]string, strRequestPath string) string {
// 	return nil
// }

/*Method: POST and Signature is required  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Coineal) ApiKeyRequest(requestMethod string, mapParams map[string]string, strRequestPath string) string {
	timestamp := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	log.Printf("timestamp :%v", timestamp)
	//Signature Request Params
	mapParams["api_key"] = e.API_KEY
	mapParams["time"] = fmt.Sprintf("%v", timestamp)
	mapParams["sign"] = MD5(CreatePayload(mapParams), e.API_SECRET)

	strUrl := API_URL + strRequestPath

	if requestMethod == "GET" {
		return exchange.HttpGetRequest(strUrl, mapParams)
	} else {
		return e.PostReq(strUrl, mapParams)
	}
}

func (e *Coineal) PostReq(resource string, payload map[string]string) string {

	c := &http.Client{}

	body := []byte{}

	var rawurl string
	if strings.HasPrefix(resource, "http") {
		rawurl = resource
	} else {
		rawurl = fmt.Sprintf("%s/%s", API_URL, resource)
	}

	formValues := url.Values{}
	for key, value := range payload {
		formValues.Add(key, value)
	}

	formData := formValues.Encode()

	req, err := http.NewRequest("POST", rawurl, strings.NewReader(formData))
	if err != nil {
		return err.Error()
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Add("Accept", "application/json")

	resp, err := c.Do(req)

	if err != nil {
		log.Printf("err=%v", err)
		return err.Error()
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	return string(body)

}

func CreatePayload(mapParams map[string]string) string {
	keys := make([]string, 0, len(mapParams))
	for key := range mapParams {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	payload := ""
	for _, key := range keys {
		payload += (key + mapParams[key])
	}
	return payload
}
func MD5(strMessage string, strSecret string) string {
	h := md5.New()
	h.Write([]byte(strMessage + strSecret))
	return hex.EncodeToString(h.Sum(nil))
}
