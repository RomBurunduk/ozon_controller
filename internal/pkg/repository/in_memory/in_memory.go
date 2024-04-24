package in_memory

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"pvz_controller/internal/pkg/repository"
)

type InMemoryCache struct {
	PVZ   map[repository.PVZDbId]cachePVZ
	mxPVZ sync.RWMutex
}

type cachePVZ struct {
	v          repository.PvzDb
	expiration time.Duration
	createdAt  time.Time
}

func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		PVZ:   make(map[repository.PVZDbId]cachePVZ, 100),
		mxPVZ: sync.RWMutex{},
	}

	go func() {
		t := time.NewTicker(10 * time.Minute)
		for {
			select {
			case <-t.C:
				now := time.Now()
				newCache := make(map[repository.PVZDbId]cachePVZ)
				for key, v := range cache.PVZ {
					if v.createdAt.Sub(now).Hours() <= v.expiration.Hours() {
						newCache[key] = v
					}
				}
				cache.PVZ = newCache
			default:

			}
		}
	}()
	return cache
}

func (c *InMemoryCache) Set(id repository.PVZDbId, item repository.PvzDb, expiration time.Duration) {
	c.mxPVZ.RLock()
	defer c.mxPVZ.RUnlock()
	c.PVZ[id] = cachePVZ{
		v:          item,
		expiration: expiration,
		createdAt:  time.Now(),
	}
}

func (c *InMemoryCache) Get(id repository.PVZDbId) (repository.PvzDb, error) {
	c.mxPVZ.RLock()
	defer c.mxPVZ.RUnlock()
	item, ok := c.PVZ[id]
	if !ok {
		return repository.PvzDb{}, errors.New("cant find item by id")
	}
	return item.v, nil
}

func (c *InMemoryCache) Delete(id repository.PVZDbId) {
	c.mxPVZ.Lock()
	defer c.mxPVZ.Unlock()
	delete(c.PVZ, id)
}
