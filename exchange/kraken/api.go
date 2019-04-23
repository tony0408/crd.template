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


	strRequest := "/private/Balance"

	jsonBanlanceReturn := uInstance.ApiKeyPost(make(map[string]string), strRequest)
	log.Printf("Balance Json: %v", jsonBanlanceReturn) 

	//TODO: GetBalance
}

/*Withdraw the coin to another address
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Call ApiKey Function (Depend on API request)
Step 5: Check the success of withdraw*/
func (e *Kraken) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	return false
}

/*Get the Status of a Singal Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Change Order Status (Status reference ../market/market.go)*/
func (e *Kraken) OrderStatus(order *market.Order) error {
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
	return nil
}

/*Place a limit Sell Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Kraken) LimitSell(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
	return nil, nil
}

/*Place a limit Buy Order  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)
Step 4: Create mapParams & Call ApiKey Function (Depend on API request)
Step 5: Create a new Order*/
func (e *Kraken) LimitBuy(pair *pair.Pair, quantity, rate float64) (*market.Order, error) {
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
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
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
