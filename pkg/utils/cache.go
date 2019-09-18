package utils

import (
	"sync"
	"time"
)


/*
General purpose caching system, with support for custom "root" duration.

A cache key is basically a string composed by a "root" and a "key".
The root part is used to customize how different
cache elements have to behave.
*/


type (
	//map key ==> cacheElement
	cache struct{
		elements map[string]cacheElement
		lock sync.RWMutex

		cacheRootElementsDuration map[CacheElementRoot]time.Duration
	}
	cacheElement struct {
		deadline time.Time
		value    interface{}
	}

	/*
	I may want to differentiate between keys as imdbID or keys as direct name, so examples of
	root may be 'resourceImdbID' or 'resourceName'
	*/
	CacheElementRoot string
	CacheElementKey string
)
const(
	standardCacheElementDuration = 10 * time.Minute
)

var(
	localCacheExpirationCheckStarted = false
	cacheRefreshTime = 1 * time.Second
)

func NewCache() *cache{
	return &cache{
		elements:make(map[string]cacheElement,1000),
		cacheRootElementsDuration:make(map[CacheElementRoot]time.Duration),
	}
}

func (c *cache) SetNewRootElementDuration(iRoot CacheElementRoot, iDuration time.Duration){
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cacheRootElementsDuration[iRoot] = iDuration
}
//loops through localCache elements and remove expired one
func (c *cache) StartLocalExpirationFetch(){
	if localCacheExpirationCheckStarted {
		Logger.Warn("local cache expiration check already started")
	}


	//not really thread safe, but a lock here is an overkill
	localCacheExpirationCheckStarted = true
	//just in case this function stops
	defer func(){localCacheExpirationCheckStarted = false}()

	for{
		time.Sleep(cacheRefreshTime)
		c.lock.Lock()
		for key, element := range c.elements{
			if element.Deadline(){
				delete(c.elements,key)
			}
		}
		c.lock.Unlock()
	}
}


func (cer CacheElementRoot) String() string{
	return string(cer)
}
func (cek CacheElementKey) String() string{
	return string(cek)
}


//true if this element reached the deadline
func (ce *cacheElement) Deadline() bool{
	return !(ce.deadline.After(time.Now())) //true if it now or before now
}
func (ce *cacheElement) Value() interface{}{
	return ce.value
}




/*
I may want to differentiate between keys as imdbID or keys as direct name
*/
func (c *cache) Insert(iCacheRoot CacheElementRoot,iCacheKey CacheElementKey, iValue interface{}){
	c.lock.Lock()
	defer c.lock.Unlock()

	elementDuration,ok := c.cacheRootElementsDuration[iCacheRoot]
	if !ok{
		elementDuration = standardCacheElementDuration
	}

	c.elements[generateCacheKey(iCacheRoot,iCacheKey)] = cacheElement{
		deadline:time.Now().Add(elementDuration),
		value:iValue,
	}
}
func (c *cache) Fetch(iCacheRoot CacheElementRoot,iCacheKey CacheElementKey) (oValue interface{}){
	c.lock.RLock()
	defer c.lock.RUnlock()

	cachedElement, ok := c.elements[generateCacheKey(iCacheRoot,iCacheKey)]
	if !ok{
		return nil
	}
	return cachedElement.value
}
func (c *cache) Reset(){
	c.lock.Lock()
	defer c.lock.Unlock()
	c = NewCache()
}

func generateCacheKey(iCacheRoot CacheElementRoot,iCacheKey CacheElementKey) string{
	//TODO string concatenation is slow, maybe use a buffer?
	return iCacheRoot.String()+"_"+iCacheKey.String()
}
