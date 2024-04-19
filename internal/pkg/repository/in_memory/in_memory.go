package in_memory

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"pvz_controller/internal/pkg/repository"
)

type InMemoryCache struct {
	PVZ map[repository.PVZDbId]struct {
		repository.PvzDb
		time.Time
	}
	mxPVZ sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		PVZ: make(map[repository.PVZDbId]struct {
			repository.PvzDb
			time.Time
		}, 100),
		mxPVZ: sync.RWMutex{},
	}

	go func() {
		t := time.NewTicker(10 * time.Minute)
		for {
			select {
			case <-t.C:
				now := time.Now()
				for key, v := range cache.PVZ {
					if v.Time.Sub(now).Hours() > 12 {
						delete(cache.PVZ, key)
					}
				}
			default:

			}
		}
	}()
	return cache
}

func (c *InMemoryCache) Set(id repository.PVZDbId, item repository.PvzDb) {
	c.mxPVZ.RLock()
	defer c.mxPVZ.Unlock()
	c.PVZ[id] = struct {
		repository.PvzDb
		time.Time
	}{item, time.Now()}
}

func (c *InMemoryCache) Get(id repository.PVZDbId) (repository.PvzDb, error) {
	c.mxPVZ.RLock()
	defer c.mxPVZ.RUnlock()
	item, ok := c.PVZ[id]
	if !ok {
		return repository.PvzDb{}, errors.New("cant find item by id")
	}
	return item.PvzDb, nil
}

func (c *InMemoryCache) Delete(id repository.PVZDbId) {
	c.mxPVZ.Lock()
	defer c.mxPVZ.Unlock()
	delete(c.PVZ, id)
}
