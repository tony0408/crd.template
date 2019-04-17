package fcoin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
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
	orderBook := OrderBook{}
	symbol := strings.ToLower(e.GetPairCode(p))

	strRequestUrl := fmt.Sprintf("market/depth/L150/%s", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &market.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonFcoinOrderbook := exchange.HttpGetRequest(strUrl, nil)


	err := json.Unmarshal([]byte(jsonFcoinOrderbook), &orderBook)
	if err != nil {
		return nil, fmt.Errorf("Fcoin OrderBook json Unmarshal error:%v", err)
	}


	//Convert Exchange Struct to Maker
	maker.Timestamp = float64(orderBook.Data.Ts)
	var buyRate float64
	for i, bid := range orderBook.Data.Bids {
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
	for i, ask := range orderBook.Data.Asks {
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
func GetFcoinCoin() *CoinsData {

	coinsData := &CoinsData{}

	strRequestUrl := "public/currencies"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		log.Printf("Fcoin Get Coin Json Unmarshal Err: %v %v", err, jsonCurrencyReturn)
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

	pairsData := &PairsData{}

	strRequestUrl := "public/symbols"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)

	err := json.Unmarshal([]byte(jsonCurrencyReturn), &pairsData)
	if err != nil {
		log.Printf("Fcoin Get Pairs Json Unmarshal Err: %v %v", err, jsonCurrencyReturn)
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
		log.Printf("Fcoin API Key or Secret Key are nil.")
		return
	}

	/*	//TODO: GetBalance
		jsonResponse := JsonResponse{}
		accountBalance := AccountBalances{}
		strRequestUrl := "/accounts/balance"

		jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
		err := json.Unmarshal([]byte(jsonCurrencyReturn), &accountBalance)
		if err != nil {
			log.Printf("Fcoin Get Balance Json Unmarshal Err: %v %v", err, jsonCurrencyReturn)
			return nil
		}
	*/
	/*
		jsonBalanceReturn := uInstance.ApiKeyPost(make(map[string]interface{}), strRequestUrl)
		if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
			log.Printf("Fcoin Get Balance Json Unmarshal Err: %v %v", err, jsonBalanceReturn)
			return
		} else if !jsonResponse.Success {
			log.Printf("Fcoin Get Balance Err: %v %v", jsonResponse.Error, jsonResponse.Message)
			return
		}

		if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
			log.Printf("Fcoin Get Balance Data Unmarshal Err: %v %v", err, jsonResponse.Data)
			return
		} else {

			for _, data := range accountBalance.Data {
				c := coin.GetCoin(e.GetCode(data.Currency))
				if c != nil {
					balanceMap.Set(c.Code, data.Available)
				}
			}
		}
	*/
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
	/*
		JsonResponse := JsonResponse{}
		strRequest := "/assets/accounts/assets-to-spot"

		mapParams := make(map[string]interface{})
		mapParams["Currency"] = e.GetSymbol(coin.Code)
		mapParams["Address"] = addr
		mapParams["PaymentId"] = coin.Code
		mapParams["Amount"] = fmt.Sprintf("%f", quantity)

		jsonSubmitWithdraw := e.ApiKeyPost(mapParams, strRequest)
		if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
			log.Printf("Fcoin Withdraw Json Unmarshal failed: %v %v", err, jsonSubmitWithdraw)
			return false
		} else if !jsonResponse.Success {
			log.Printf("Fcoin Withdraw failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
			return false
		}
	*/
	return true
}

/*Get the Status of a Singal Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (Status reference ../market/market.go)*/
func (e *Fcoin) OrderStatus(order *market.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("Fcoin API Key or Secret Key are nil.")
	}
	/*
		jsonResponse := JsonResponse{}
		orderStatus := TradeHistory{}
		strRequest := "/orders/" + order.OrderID

		mapParams := make(map[string]interface{})
		mapParams["Market"] = fmt.Sprintf("%s/%s", e.GetSymbol(order.Pair.Target.Code), e.GetSymbol(order.Pair.Base.Code))

		jsonOrderStatus := e.ApiKeyPost(mapParams, strRequest)
		if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
			return fmt.Errorf("Fcoin OrderStatus Unmarshal Err: %v %v", err, jsonOrderStatus)
		} else if !jsonResponse.Success {
			return fmt.Errorf("Fcoin Get OrderStatus failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
		}

		if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
			return fmt.Errorf("Fcoin Get OrderStatus Data Unmarshal Err: %v %v", err, jsonResponse.Data)
		} else {
			for _, list := range orderStatus {
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
			}
		}
	*/
	return nil
}

func (e *Fcoin) ListOrders() (*[]market.Order, error) {
	return nil, nil
}

/*Cancel an Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (order.Status = market.Canceling)*/
func (e *Fcoin) CancelOrder(order *market.Order) error {
	return nil
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

/*Place a limit Sell Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Fcoin) LimitSell(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	return nil, nil
}

/*Place a limit Buy Order  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Fcoin) LimitBuy(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	return nil, nil
}

/*************** Signature Http Request ***************/
/*Method: GET and Signature is required  --reference Cryptopia
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Fcoin) ApiKeyGet(mapParams map[string]string, strRequestPath string) string {
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
func (e *Fcoin) ApiKeyPost(mapParams map[string]string, strRequestPath string) string {
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
