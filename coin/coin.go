package coin

import (
	"errors"
	"strings"
)

// modeling structure and functions,
// don't modify unless bugs or new features

type Coin struct { // might also understand as public chain or token
	ID         string `json:"_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Website    string `json:"website"`
	Explorer   string `json:"explorer"`
	Health     string // the health of the chain
	Blockheigh int
	Blocktime  int // in seconds
	Blocklast  int // the last blocktime in UTC

}

var coinmap = make(map[string]*Coin)

func GetCoin(code string) *Coin {
	return coinmap[strings.ToUpper(code)]
}

func GetCoins() map[string]*Coin {
	return coinmap
}

func AddCoin(coin *Coin) error {

	if coin != nil && coin.Code != "" {
		coin.Code = strings.ToUpper(coin.Code)
		coinmap[coin.Code] = coin
	} else {
		return errors.New("code is not assign yet")
	}
	return nil
}

/* func SetCoin() error {
	e := redis.CreateRedis()
	e.Init("37.59.46.161:16379", 0)
	c, err := json.Marshal(coinmap)
	if err != nil {
		return err
	}
	return e.Set("CoinList", string(c))
} */
