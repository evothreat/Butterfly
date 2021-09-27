package cnc

import (
	"sync"
	"time"
)

type CookieData struct {
	Expires time.Time
}

type CookieDataStore struct {
	mutex      sync.RWMutex
	cookieData map[string]*CookieData
}

func NewCookieData() *CookieData {
	return &CookieData{}
}

func NewCookieDataStore() *CookieDataStore {
	return &CookieDataStore{
		cookieData: map[string]*CookieData{},
	}
}

func (cds *CookieDataStore) Put(key string, data *CookieData) {
	cds.mutex.Lock()
	cds.cookieData[key] = data
	cds.mutex.Unlock()
}

func (cds *CookieDataStore) Get(key string) *CookieData {
	cds.mutex.RLock()
	data, ok := cds.cookieData[key]
	cds.mutex.RUnlock()
	if ok {
		return data
	}
	return nil
}

func (cds *CookieDataStore) Delete(key string) {
	cds.mutex.Lock()
	delete(cds.cookieData, key)
	cds.mutex.Unlock()
}
