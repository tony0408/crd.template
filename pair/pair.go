package pair

import (
	coin "../coin"
	cmap "github.com/orcaman/concurrent-map"
)

//Pair is the common name pairs across diff excahnges
type Pair struct {
	Name   string
	Base   *coin.Coin
	Target *coin.Coin
}

// var pairMap = make(map[string]*Pair)
var pairMap cmap.ConcurrentMap

func Init() {
	if pairMap == nil {
		pairMap = cmap.New()
	}
}

func GetString(pair *Pair) string {
	key := GetKey(pair.Base, pair.Target)
	return key
}

func GetKey(base, target *coin.Coin) string {
	key := ""
	if base != nil && target != nil {
		// log.Printf("base:%v  %s  target:%v", base.Code, coin.SEPARATOR, target.Code)
		key = (base.Code + coin.SEPARATOR + target.Code)
	}
	return key
}

func GetPairByKey(key string) *Pair {
	// return pairMap[key]

	if tmp, ok := pairMap.Get(key); ok {
		return tmp.(*Pair)
	}
	return nil
}

func GetPair(base, target *coin.Coin) *Pair {
	// log.Printf("base:%v	target:%v", base, target)
	key := GetKey(base, target)
	// if pairMap[key] == nil {
	// 	pairMap[key] = &Pair{key, base, target}
	// }
	// return pairMap[key]

	if tmp, ok := pairMap.Get(key); ok {
		return tmp.(*Pair)
	} else {
		p := &Pair{key, base, target}
		pairMap.Set(key, p)
		return p
	}
}

func GetPairs() []*Pair {
	pairs := []*Pair{}

	// for _, pair := range pairMap {
	// 	pairs = append(pairs, pair)
	// }
	// return pairs

	for _, key := range pairMap.Keys() {
		pairs = append(pairs, GetPairByKey(key))
	}
	return pairs

}
