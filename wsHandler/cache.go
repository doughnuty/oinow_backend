package wsHandler

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"nhooyr.io/websocket"
)

type SyncCache struct {
	cache *cache.Cache
	mut   *sync.Mutex
}

func (c *SyncCache) Set(id string, val interface{}) {
	c.mut.Lock()
	c.cache.Set(id, val, time.Minute*30)
	c.mut.Unlock()
}
func (c *SyncCache) Get(id string) (interface{}, bool) {
	c.mut.Lock()
	res, exists := c.cache.Get(id)
	c.mut.Unlock()
	if !exists {
		return nil, exists
	}
	return res, exists
}

func (c *SyncCache) Delete(id string, conn *websocket.Conn) {
	c.mut.Lock()
	resIntf, _ := c.cache.Get(id)
	res := resIntf.([]*websocket.Conn)

	i := 0
	for i = range res {
		if res[i] == conn {
			break
		}
	}

	res[i] = res[len(res)-1]
	res = res[:len(res)-1]

	c.cache.Set(id, res, time.Minute*30)
	c.mut.Unlock()
}

var ConnsCache = &SyncCache{
	cache: cache.New(time.Minute*30, time.Minute*30),
	mut:   &sync.Mutex{},
}

var UsernamesCache = &SyncCache{
	cache: cache.New(time.Minute*30, time.Minute*30),
	mut:   &sync.Mutex{},
}
