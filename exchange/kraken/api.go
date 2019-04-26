package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
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
	API_URL string = "https://api.kraken.com/0"
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
func (e *Kraken) OrderBook(p *pair.Pair) (*market.Maker, error) {
	response := ResponseReturn{}

	strRequestUrl := "/public/Depth"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["pair"] = e.GetPairCode(p)

	jsonResponseReturn := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonResponseReturn), &response); err != nil {
		return nil, fmt.Errorf("Kraken Unmarshal Response error: %s", err)
	}
	if len(response.Error) != 0 {
		return nil, fmt.Errorf("Kraken OrderBook error: %s", response.Error)
	}

	data := make(map[string]*OrderBook)
	if err := json.Unmarshal(response.Result, &data); err != nil {
		return nil, fmt.Errorf("Kraken Unmarshal Result error: %s", err)
	}

	//Convert Exchange Struct to Maker
	maker := &market.Maker{}
	maker.Timestamp = float64(time.Now().UnixNano() / 1e6)
	for _, orderBook := range data {
		var err error
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
	}
	return maker, nil
}

/*Get Coins Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func GetKrakenCoin() map[string]*CoinData {
	response := ResponseReturn{}

	strRequestUrl := "/public/Assets"
	strUrl := API_URL + strRequestUrl

	jsonResponseReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonResponseReturn), &response); err != nil {
		log.Printf("Kraken Unmarshal Response error: %s", err)
		return make(map[string]*CoinData)
	}
	if len(response.Error) != 0 {
		log.Printf("Kraken Get Coin error: %s", response.Error)
		return make(map[string]*CoinData)
	}

	coinsInfo := make(map[string]*CoinData)
	if err := json.Unmarshal(response.Result, &coinsInfo); err != nil {
		log.Printf("Kraken Unmarshal Result error: %s", err)
		return coinsInfo
	}

	return coinsInfo
}

/*Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func GetKrakenPair() map[string]*PairData {
	response := ResponseReturn{}

	strRequestUrl := "/public/AssetPairs"
	strUrl := API_URL + strRequestUrl

	jsonResponseReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonResponseReturn), &response); err != nil {
		log.Printf("Kraken Unmarshal Response error: %s", err)
		return make(map[string]*PairData)
	}
	if len(response.Error) != 0 {
		log.Printf("Kraken Get Coin error: %s", response.Error)
		return make(map[string]*PairData)
	}

	pairsInfo := make(map[string]*PairData)
	if err := json.Unmarshal(response.Result, &pairsInfo); err != nil {
		log.Printf("Kraken Unmarshal Result error: %s", err)
		return pairsInfo
	}

	return pairsInfo
}

/*************** Private API ***************/
func (e *Kraken) UpdateAllBalances() {
	e.UpdateAllBalancesByUser(nil)
}

/*Get Exchange Account All Coins Balance  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Call ApiKey Function (Depend on API request)
Step 5: Get Coin Availamount and store in balanceMap*/
func (e *Kraken) UpdateAllBalancesByUser(u *user.User) {
	var uInstance *Kraken
	if u != nil {
		uInstance = &Kraken{}
		uInstance.API_KEY = u.API_KEY
		uInstance.API_SECRET = u.API_SECRET
		// uInstance.MakerDB = e.MakerDB
	} else {
		uInstance = e
	}

	if uInstance.API_KEY == "" || uInstance.API_SECRET == "" {
		log.Printf("Kraken API Key or Secret Key are nil.")
		return
	}

	jsonResponse := ResponseReturn{} //JsonResponse{}

	accountBalance := AccountBalances{}
	strRequest := "/private/Balance"
	mapParams := make(map[string]string)
	//mapParams["pair"] = e.GetPairCode(p)

	jsonBalanceReturn := uInstance.ApiKeyPost(mapParams, strRequest)
	log.Printf("==jsonBalanceReturn: %v", jsonBalanceReturn) //=====

	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("Kraken Get Balance Json Unmarshal Err: %v %v", err, jsonBalanceReturn)
		return
	}
	log.Printf("==jsonResponse: %v", jsonResponse) //=====
	if len(jsonResponse.Error) != 0 {
		log.Printf("Kraken Get Balance Err: %v ", jsonResponse.Error)
		return
	}

	if err := json.Unmarshal(jsonResponse.Result, &accountBalance); err != nil {
		log.Printf("Kraken Get Balance Data Unmarshal Err: %v %v", err, jsonResponse.Result)
		return
	} else {
		structInf := reflect.Indirect(reflect.ValueOf(accountBalance))
		for i := 0; i < structInf.NumField(); i++ {
			fieldName := structInf.Type().Field(i).Name
			fieldValue := structInf.Field(i)
			c := coin.GetCoin(e.GetCode(fieldName))
			if c != nil {
				available, err := strconv.ParseFloat(fieldValue.String(), 64)
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
func (e *Kraken) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("Kraken API Key or Secret Key are nil.")
		return false
	}

	jsonResponse := ResponseReturn{}
	strRequest := "/private/Withdraw"

	mapParams := make(map[string]string)

	//key = withdrawal key name, as set up on your account
	mapParams["key"] = addr
	//asset = asset being withdrawn
	mapParams["asset"] = e.GetSymbol(coin.Code)
	mapParams["Amount"] = fmt.Sprintf("%f", quantity)

	jsonSubmitWithdraw := e.ApiKeyPost(mapParams, strRequest)
	log.Printf("jsonSubmitWithdraw: %+v", jsonSubmitWithdraw)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		log.Printf("Kraken Withdraw Json Unmarshal failed: %v %v", err, jsonSubmitWithdraw)
		return false
	}
	if len(jsonResponse.Error) != 0 {
		log.Printf("Kraken jsonSubmitWithdraw Unmarshal Err: %v ", jsonResponse.Error)
		return false
	}

	return false
	//Note that the first withdrawal to an address will still have to be confirmed manually by
	//clicking a link sent to the user via e-mail, even if the withdrawal request is made via the API.
}

/*Get the Status of a Singal Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (Status reference ../market/market.go)*/
func (e *Kraken) OrderStatus(order *market.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("Kraken API Key or Secret Key are nil.")
	}

	jsonResponse := ResponseReturn{}
	orderStatus := Order{}
	strRequest := "/private/QueryOrders" //"/private/OpenOrders"

	mapParams := make(map[string]string)
	mapParams["order_txid"] = order.OrderID

	jsonOrderStatus := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("Kraken OrderStatus Unmarshal Err: %v %v", err, jsonOrderStatus)
	}
	if len(jsonResponse.Error) != 0 {
		log.Printf("Kraken jsonOrderStatus Unmarshal Err: %v ", jsonResponse.Error)
		return nil
	}

	if err := json.Unmarshal(jsonResponse.Result, &orderStatus); err != nil {
		return fmt.Errorf("Kraken Get OrderStatus Data Unmarshal Err: %v %s", err, jsonResponse.Result)
	} else {

		if orderStatus.TransactionID == order.OrderID {
			if strings.Contains(orderStatus.Misc, "partial") {
				order.Status = market.Partial
			}
		}
		// need change
		/* for _, list := range orderStatus {
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
		} */
	}

	return nil
}

func (e *Kraken) ListOrders() (*[]market.Order, error) {
	return nil, nil
}

func (e *Kraken) CancelAllOrder() error {
	return nil
}

/*Cancel an Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (order.Status = market.Canceling)*/
func (e *Kraken) CancelOrder(order *market.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("Kraken API Key or Secret Key are nil.")
	}

	jsonResponse := ResponseReturn{}
	strRequest := "/private/CancelOrder"

	mapParams := make(map[string]string)
	mapParams["txid"] = order.OrderID

	jsonCancelOrder := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("Kraken CancelOrder Unmarshal Err: %v %v", err, jsonCancelOrder)
	}
	if len(jsonResponse.Error) != 0 {
		return fmt.Errorf("Kraken CancelOrder Unmarshal Err: %+v", jsonResponse.Error)
	}

	order.Status = market.Canceling

	return nil
}

/*Place a limit Sell Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Kraken) LimitSell(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("Kraken API Key or Secret Key are nil.")
	}

	jsonResponse := ResponseReturn{}
	placeOrder := AddOrderResponse{}
	strRequest := "/private/AddOrder"

	mapParams := make(map[string]string)
	mapParams["pair"] = strings.ToLower(e.GetPairCode(pair))
	mapParams["type"] = "sell"
	mapParams["ordertype"] = "limit"
	mapParams["price"] = fmt.Sprint(rate)
	mapParams["volume"] = fmt.Sprint(quantity)

	jsonPlaceReturn := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("Kraken LimitSell Unmarshal Err: %v %v", err, jsonPlaceReturn)
	}
	if len(jsonResponse.Error) != 0 {
		return nil, fmt.Errorf("Kraken LimitSell Unmarshal Err: %+v", jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("Kraken LimitSell Data Unmarshal Err: %v %v", err, jsonResponse.Result)
	} else {
		order := &market.Order{
			OrderID:      placeOrder.TransactionIds[0],
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
func (e *Kraken) LimitBuy(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("Kraken API Key or Secret Key are nil.")
	}

	jsonResponse := ResponseReturn{}
	placeOrder := AddOrderResponse{}
	strRequest := "/private/AddOrder"

	mapParams := make(map[string]string)
	mapParams["pair"] = strings.ToLower(e.GetPairCode(pair))
	mapParams["type"] = "buy"
	mapParams["ordertype"] = "limit"
	mapParams["price"] = fmt.Sprint(rate)
	mapParams["volume"] = fmt.Sprint(quantity)

	jsonPlaceReturn := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("Kraken LimitBuy Unmarshal Err: %v %v", err, jsonPlaceReturn)
	}
	if len(jsonResponse.Error) != 0 {
		return nil, fmt.Errorf("Kraken LimitBuy Unmarshal Err: %+v", jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("Kraken LimitBuy Data Unmarshal Err: %v %v", err, jsonResponse.Result)
	} else {
		order := &market.Order{
			OrderID:      placeOrder.TransactionIds[0],
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

/*Method: POST and Signature is required  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Kraken) ApiKeyPost(mapParams map[string]string, strRequestPath string) string {
	strMethod := "POST"

	//Signature Request Params
	mapParams["nonce"] = fmt.Sprintf("%d", time.Now().UnixNano())
	if e.Two_Factor != "" {
		mapParams["otp"] = e.Two_Factor
	}
	// two-factor password added
	mapParams["otp"] = "kraken123."
	Signature := ComputeHmac512(strRequestPath, mapParams, e.API_SECRET)

	strUrl := API_URL + strRequestPath

	httpClient := &http.Client{}

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	request, err := http.NewRequest(strMethod, strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	//request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	//Signature = Signature + "111"
	//Signature correct
	log.Printf("Key!!!: %v", e.API_KEY)       //============
	log.Printf("Secret!!!: %v", e.API_SECRET) //============
	log.Printf("Signature!!!: %v", Signature) //============
	request.Header.Add("API-Key", e.API_KEY)
	request.Header.Add("API-Sign", Signature)

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
func ComputeHmac512(strPath string, mapParams map[string]string, strSecret string) string {
	bytesParams, _ := json.Marshal(mapParams)
	sha := sha256.New()
	sha.Write(bytesParams)
	shaSum := sha.Sum(nil)

	strMessage := fmt.Sprintf("%s%s", strPath, string(shaSum))
	decodeSecret, _ := base64.StdEncoding.DecodeString(strSecret)

	h := hmac.New(sha512.New, decodeSecret)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
