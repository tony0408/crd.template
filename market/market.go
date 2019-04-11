package market

import (
	"sync"

	"../pair"
)

type OrderStatus string

const (
	New       OrderStatus = "New"
	Filled    OrderStatus = "Filled"
	Partial   OrderStatus = "Partial"
	Canceling OrderStatus = "Canceling"
	Canceled  OrderStatus = "Canceled"
	Rejected  OrderStatus = "Rejected"
	Expired   OrderStatus = "Expired"
	Other     OrderStatus = "Other"
)

type Order struct {
	Pair          *pair.Pair
	OrderID       string
	FilledOrders  []int64
	Rate          float64 `bson:"Rate"`
	Quantity      float64 `bson:"Quantity"`
	Side          string
	Status        OrderStatus `json:"status"`
	StatusMessage string
	DealRate      float64
	DealQuantity  float64
	JsonResponse  string
}

//Pair is the common name pairs across diff excahnges
type Maker struct {
	WorkerIP        string  `bson:"workerip"`
	BeforeTimestamp float64 `bson:"beforetimestamp"`
	AfterTimestamp  float64 `bson:"aftertimestamp"`
	KafkaTimestamp  float64 `bson:"kafkatimestamp"`
	Timestamp       float64 `bson:"timestamp"`
	Nounce          int     `bson:"Nounce"`
	LastUpdateID    int     `json:"lastUpdateId"`
	Bids            []Order `json:"bids"`
	Asks            []Order `json:"asks"`
}

type Market struct{}

// var makerMap = make(map[*pair.Pair]*Maker)
var instanceMarket *Market
var once sync.Once

func CreateMarket() *Market {
	once.Do(func() {
		instanceMarket = &Market{}

	})
	return instanceMarket
}

// func (m *Market) GetMaker(key *pair.Pair) *Maker {

// 	if makerMap[key] == nil {
// 		makerMap[key] = &Maker{}
// 	}
// 	return makerMap[key]
// }

// func (m *Market) UpdateMaker(key *pair.Pair, maker *Maker) {

// 	makerMap[key] = maker

// }
