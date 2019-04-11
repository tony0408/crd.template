package db

import (
	"errors"
	"sync"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
	lock   sync.RWMutex
}

func CreateRedis() *Redis {

	r := &Redis{}
	r.lock = sync.RWMutex{}
	return r
}

func (r *Redis) Init(addr string, no int) *Redis {
	r.client = redis.NewClient(&redis.Options{
		Addr:     addr, // "45.63.39.254:6379",
		Password: "",   // no password set
		DB:       no,   // use default DB
	})
	return r
}

func (r *Redis) Set(key string, val interface{}) error {
	// r.lock.RLock()
	// defer r.lock.RUnlock()

	//log.Printf("key:%v    val:%v", key, val)
	return r.client.Set(key, val, 0).Err()
}

func (r *Redis) GetSet(key string, val interface{}) error {
	return r.client.GetSet(key, val).Err()
}

func (r *Redis) Get(key string) (interface{}, error) {
	// r.lock.RLock()
	// defer r.lock.RUnlock()
	result := r.client.Get(key)
	// log.Printf("key:%v    val:%v", key, result)
	if result != nil {
		return result.Result()
	}

	return nil, errors.New("the key is not found in DB")
}

func (r *Redis) Close() {
	r.client.Close()
	r.client = nil
}

/*

func (d *Redis) UpdateMaker(exchange string, pair *pair.Pair, maker *market.Maker) error {
	m, err := json.Marshal(maker)
	if err != nil {
		return err
	}
	// fmt.Println(string(m))
	return client.Set(exchange+"-"+pair.Name, m, 0).Err()
}

func (d *Redis) GetMaker(exchange string, pair *pair.Pair) (maker *market.Maker, err error) {
	val, err := client.Get(exchange + "-" + pair.Name).Result()

	if err := json.Unmarshal([]byte(val), &maker); err != nil {
		return nil, err
	}

	return maker, err

}

func ExampleNewClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "45.63.39.254:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	err = client.Set("key2", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

	// Output: key value
	// key2 does not exist
}

*/
