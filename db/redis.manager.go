package db

import (
	"log"
	"sync"
	"time"

	"github.com/orcaman/concurrent-map"
)

const (
	MAX_CONN  = 3
	GC_PEROID = 30 // seconds
)

type RedisManager struct {
	lock     sync.RWMutex
	redisMap cmap.ConcurrentMap
}

var ShutdownSignal bool

func CreateRedisManager() *RedisManager {
	rm := &RedisManager{}
	rm.lock = sync.RWMutex{}
	rm.redisMap = cmap.New()

	ShutdownSignal = false
	go AutoGC(rm)
	return rm
}

func (rm *RedisManager) Add(key string, r *Redis) {
	if r == nil {
		log.Printf("r=nil error")
		return
	}
	s := rm.Get(key)
	if s != nil {
		s.Close()
	}
	rm.redisMap.Set(key, r)
}

func (rm *RedisManager) Get(key string) *Redis {
	if rm.redisMap.Count() > 0 {
		if tmp, ok := rm.redisMap.Get(key); ok {
			return tmp.(*Redis)
		}
	}
	return nil
}

func (rm *RedisManager) Close() {
	ShutdownSignal = true
}

func AutoGC(rm *RedisManager) {
	for {
		if ShutdownSignal {
			break
		}
		time.Sleep(time.Second * GC_PEROID)

		size := rm.redisMap.Count()
		if size > MAX_CONN {
			removeSize := MAX_CONN - size
			for item := range rm.redisMap.Iter() { //TODO need to add TTL or sorting to make sure want to remove the oldest, useless one
				item.Val.(*Redis).Close()
				rm.redisMap.Remove(item.Key)
				removeSize--
				if removeSize < 0 {
					break
				}
			}
		}

	}
}
