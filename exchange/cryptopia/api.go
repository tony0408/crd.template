package cryptopia

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"../../coin"
	"../../exchange"
	"../../market"
	"../../pair"
	"../../user"
)

const (
	API_URL string = "https://www.cryptopia.co.nz"
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
func (e *Cryptopia) OrderBook(p *pair.Pair) (*market.Maker, error) {
	jsonResponse := JsonResponse{}
	orderbook := OrderBook{}
	symbol := e.GetPairCode(p)

	strRequestUrl := fmt.Sprintf("/api/GetMarketOrders/%s", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &market.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonMarketDepthReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonMarketDepthReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("Cryptopia OrderBook json Unmarshal error: %v %v", err, jsonMarketDepthReturn)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("Cryptopia OrderBook failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
	}

	if err := json.Unmarshal(jsonResponse.Data, &orderbook); err != nil {
		return nil, fmt.Errorf("Cryptopia OrderBook Data Unmarshal error: %v %s", err, jsonResponse.Data)
	} else {
		//Convert Exchange Struct to Maker
		for _, bid := range orderbook.Buy {
			var buydata market.Order

			buydata.Rate = bid.Price
			buydata.Quantity = bid.Volume

			maker.Bids = append(maker.Bids, buydata)
		}
		for _, ask := range orderbook.Sell {
			var selldata market.Order

			selldata.Rate = ask.Price
			selldata.Quantity = ask.Volume

			maker.Asks = append(maker.Asks, selldata)
		}
		return maker, nil
	}
}

/*Get Coins Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func GetCryptopiaCoin() *CoinsData {
	jsonResponse := JsonResponse{}
	coinsData := &CoinsData{}

	strRequestUrl := "/api/GetCurrencies"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		log.Printf("Cryptopia Get Coin Json Unmarshal Err: %v %v", err, jsonCurrencyReturn)
		return nil
	} else if !jsonResponse.Success {
		log.Printf("Cryptopia Get Coin Err: %v %v", jsonResponse.Error, jsonResponse.Message)
		return nil
	}

	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		log.Printf("Cryptopia Get Coin Data Unmarshal error: %v %s", err, jsonResponse.Data)
		return nil
	}
	return coinsData
}

/*Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func GetCryptopiaPair() *PairsData {
	jsonResponse := JsonResponse{}
	pairsData := &PairsData{}

	strRequestUrl := "/api/GetTradePairs"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		log.Printf("Cryptopia Get Pairs Json Unmarshal Err: %v %v", err, jsonCurrencyReturn)
		return nil
	} else if !jsonResponse.Success {
		log.Printf("Cryptopia Get Pairs Err: %v %v", jsonResponse.Error, jsonResponse.Message)
		return nil
	}

	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		log.Printf("Cryptopia Get Pairs Data Unmarshal error: %v %s", err, jsonResponse.Data)
		return nil
	}
	return pairsData
}

/*************** Private API ***************/
func (e *Cryptopia) UpdateAllBalances() { // Get Exchange Account All Coins Balance
	e.UpdateAllBalancesByUser(nil)
}

/*Get Exchange Account All Coins Balance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Call ApiKey Function (Depend on API request)
Step 5: Get Coin Availamount and store in balanceMap*/
func (e *Cryptopia) UpdateAllBalancesByUser(u *user.User) {
	var uInstance *Cryptopia
	if u != nil {
		uInstance = &Cryptopia{}
		uInstance.API_KEY = u.API_KEY
		uInstance.API_SECRET = u.API_SECRET
		// uInstance.MakerDB = e.MakerDB
	} else {
		uInstance = e
	}

	if uInstance.API_KEY == "" || uInstance.API_SECRET == "" {
		log.Printf("Cryptopia API Key or Secret Key are nil.")
		return
	}

	jsonResponse := JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/api/GetBalance"

	jsonBalanceReturn := uInstance.ApiKeyPost(make(map[string]interface{}), strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("Cryptopia Get Balance Json Unmarshal Err: %v %v", err, jsonBalanceReturn)
		return
	} else if !jsonResponse.Success {
		log.Printf("Cryptopia Get Balance Err: %v %v", jsonResponse.Error, jsonResponse.Message)
		return
	}

	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("Cryptopia Get Balance Data Unmarshal Err: %v %s", err, jsonResponse.Data)
		return
	} else {
		for _, data := range accountBalance {
			c := coin.GetCoin(e.GetCode(data.Symbol))
			if c != nil {
				balanceMap.Set(c.Code, data.Available)
			} else {
				// TODO: Add new coins
				// log.Printf("%s %v", e.GetCode(data.Symbol), c)
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
func (e *Cryptopia) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("Cryptopia API Key or Secret Key are nil.")
		return false
	}

	jsonResponse := JsonResponse{}
	strRequest := "/api/SubmitWithdraw"

	mapParams := make(map[string]interface{})
	mapParams["Currency"] = e.GetSymbol(coin.Code)
	mapParams["Address"] = addr
	mapParams["PaymentId"] = coin.Code
	mapParams["Amount"] = fmt.Sprintf("%f", quantity)

	jsonSubmitWithdraw := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		log.Printf("Cryptopia Withdraw Json Unmarshal failed: %v %v", err, jsonSubmitWithdraw)
		return false
	} else if !jsonResponse.Success {
		log.Printf("Cryptopia Withdraw failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
		return false
	}

	return true
}

/*Get the Status of a Singal Order
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (Status reference ../market/market.go)*/
func (e *Cryptopia) OrderStatus(order *market.Order) error { // Get the Status of a Singal Order
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("Cryptopia API Key or Secret Key are nil.")
	}

	jsonResponse := JsonResponse{}
	orderStatus := TradeHistory{}
	strRequest := "/api/GetOpenOrders"

	mapParams := make(map[string]interface{})
	mapParams["Market"] = fmt.Sprintf("%s/%s", e.GetSymbol(order.Pair.Target.Code), e.GetSymbol(order.Pair.Base.Code))

	jsonOrderStatus := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("Cryptopia OrderStatus Unmarshal Err: %v %v", err, jsonOrderStatus)
	} else if !jsonResponse.Success {
		return fmt.Errorf("Cryptopia Get OrderStatus failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
	}

	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("Cryptopia Get OrderStatus Data Unmarshal Err: %v %s", err, jsonResponse.Data)
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
	return nil
}

func (e *Cryptopia) ListOrders() (*[]market.Order, error) {
	return nil, nil
}

func (e *Cryptopia) CancelAllOrder() error {
	return nil
}

/*Cancel an Order
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (order.Status = market.Canceling)*/
func (e *Cryptopia) CancelOrder(order *market.Order) error { // Get/Post(Depend on API) Cancel
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("Cryptopia API Key or Secret Key are nil.")
	}

	jsonResponse := JsonResponse{}
	strRequest := "/api/CancelTrade"

	mapParams := make(map[string]interface{})
	mapParams["Type"] = "Trade"
	mapParams["OrderId"], _ = strconv.Atoi(order.OrderID)

	jsonCancelOrder := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("Cryptopia CancelOrder Unmarshal Err: %v %v", err, jsonCancelOrder)
	} else if !jsonResponse.Success {
		return fmt.Errorf("Cryptopia CancelOrder failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
	}

	order.Status = market.Canceling

	return nil
}

/*Place a limit Sell Order
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Cryptopia) LimitSell(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("Cryptopia API Key or Secret Key are nil.")
	}

	jsonResponse := JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/SubmitTrade"

	mapParams := make(map[string]interface{})
	mapParams["Market"] = fmt.Sprintf("%s/%s", e.GetSymbol(pair.Target.Code), e.GetSymbol(pair.Base.Code))
	mapParams["Type"] = "Sell"
	mapParams["Rate"] = rate
	mapParams["Amount"] = quantity

	jsonPlaceReturn := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("Cryptopia LimitSell Unmarshal Err: %v %v", err, jsonPlaceReturn)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("Cryptopia LimitSell failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
	}

	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("Cryptopia LimitSell Data Unmarshal Err: %v %s", err, jsonResponse.Data)
	} else {
		order := &market.Order{}
		order.Pair = pair
		order.Rate = rate
		order.Quantity = quantity
		order.Side = "Sell"
		order.JsonResponse = jsonPlaceReturn
		if placeOrder.OrderID != 0 {
			order.OrderID = fmt.Sprintf("%d", placeOrder.OrderID)
			order.FilledOrders = placeOrder.FilledOrders
			order.Status = market.New
		} else if len(placeOrder.FilledOrders) > 0 {
			order.OrderID = "Filled"
			order.FilledOrders = placeOrder.FilledOrders
			order.Status = market.Filled
		}

		return order, nil
	}
}

/*Place a limit Buy Order
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Cryptopia) LimitBuy(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("Cryptopia API Key or Secret Key are nil.")
	}

	jsonResponse := JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/SubmitTrade"

	mapParams := make(map[string]interface{})
	mapParams["Market"] = fmt.Sprintf("%s/%s", e.GetSymbol(pair.Target.Code), e.GetSymbol(pair.Base.Code))
	mapParams["Type"] = "Buy"
	mapParams["Rate"] = rate
	mapParams["Amount"] = quantity

	jsonPlaceReturn := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("Cryptopia LimitBuy Unmarshal Err: %v %v", err, jsonPlaceReturn)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("Cryptopia LimitBuy failed:%v Message:%v", jsonResponse.Error, jsonResponse.Message)
	}

	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("Cryptopia LimitBuy Data Unmarshal Err: %v %s", err, jsonResponse.Data)
	} else {
		order := &market.Order{}
		order.Pair = pair
		order.Rate = rate
		order.Quantity = quantity
		order.Side = "Buy"
		order.JsonResponse = jsonPlaceReturn
		if placeOrder.OrderID != 0 {
			order.OrderID = fmt.Sprintf("%d", placeOrder.OrderID)
			order.FilledOrders = placeOrder.FilledOrders
			order.Status = market.New
		} else if len(placeOrder.FilledOrders) > 0 {
			order.OrderID = "Filled"
			order.FilledOrders = placeOrder.FilledOrders
			order.Status = market.Filled
		}

		return order, nil
	}
}

/*************** Signature Http Request ***************/
/*Method: POST and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Cryptopia) ApiKeyPost(mapParams map[string]interface{}, strRequestPath string) string {
	strMethod := "POST"
	Nonce := strconv.FormatInt(time.Now().UnixNano(), 10)

	//Signature Request Params
	strUrl := API_URL + strRequestPath

	jsonParams := ""
	bytesParams, err := json.Marshal(mapParams)
	if nil != mapParams {
		jsonParams = string(bytesParams)
	}
	md5 := md5.Sum([]byte(bytesParams))
	signMessage := e.API_KEY + strMethod + strings.ToLower(url.QueryEscape(strUrl)) + Nonce + base64.StdEncoding.EncodeToString(md5[:])
	Signature := ComputeHmac256(signMessage, e.API_SECRET)
	Authentication := "amx " + e.API_KEY + ":" + Signature + ":" + Nonce

	httpClient := &http.Client{}

	request, err := http.NewRequest(strMethod, strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	request.Header.Add("Authorization", Authentication)

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
func ComputeHmac256(strMessage string, strSecret string) string {
	key, _ := base64.StdEncoding.DecodeString(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
