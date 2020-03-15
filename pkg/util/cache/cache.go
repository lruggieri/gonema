/*
Package Cache provides a general purpose thread-safe caching system, with support for custom "root" duration.

A cache KEY is basically a string composed by a "root" and a "key".
The root part is used to customize how different cache elements have to behave.
The key part is the actual element key, but belongs to a root.
Different roots can have the same key.
Having several roots is useful in order to set different durations to a different set of keys.
In the end, the actual cache hashmap will have 'root_key' as entry.
*/
package cache

import (
	"fmt"
	"sync"
	"time"
)


type (
	//Cache is the main cache struct
	Cache struct{
		sync.RWMutex
		
		elements map[string]cacheElement
		cacheRootElementsDuration map[CacheElementRoot]time.Duration
	}
	//cacheElement is the actual cache element
	cacheElement struct {
		deadline time.Time
		value    interface{}
	}

	//CacheElementRoot is the root of a cache entry
	CacheElementRoot string
	//CacheElementKey is a key belonging to a CacheElementRoot. Several CacheElementRoot can have the same CacheElementKey.
	CacheElementKey string
)
const(
	standardCacheElementDuration = 10 * time.Minute
)

var(
	localCacheExpirationCheckStarted = false
	cacheRefreshTime = 1 * time.Second
)

//NewCache returns a new Cache object
func NewCache() *Cache {
	c := &Cache{
		elements:make(map[string]cacheElement,1000),
		cacheRootElementsDuration:make(map[CacheElementRoot]time.Duration),
	}
	go c.startLocalExpirationFetch()
	return c
}
//SetNewRootElementDuration sets the maximum duration for each new element inserted in the cache
//under the input CacheElementRoot
func (c *Cache) SetNewRootElementDuration(iRoot CacheElementRoot, iDuration time.Duration){
	c.Lock()
	defer c.Unlock()

	c.cacheRootElementsDuration[iRoot] = iDuration
}
//loops through localCache elements and remove expired one
func (c *Cache) startLocalExpirationFetch(){
	if localCacheExpirationCheckStarted {
		fmt.Println("local cache expiration check already started")
		return
	}


	//not really thread safe, but a lock here is an overkill
	localCacheExpirationCheckStarted = true
	//just in case this function stops
	defer func(){localCacheExpirationCheckStarted = false}()

	for{
		time.Sleep(cacheRefreshTime)
		c.Lock()
		for key, element := range c.elements{
			if element.Deadline(){
				delete(c.elements,key)
			}
		}
		c.Unlock()
	}
}

//String returns the string version of a CacheElementRoot
func (cer CacheElementRoot) String() string{
	return string(cer)
}
//String returns the string version of a CacheElementKey
func (cek CacheElementKey) String() string{
	return string(cek)
}


//Deadline returns true if the cacheElement reached the deadline
func (ce *cacheElement) Deadline() bool{
	return !(ce.deadline.After(time.Now())) //true if it now or before now
}
//Value returns the cacheElement value
func (ce *cacheElement) Value() interface{}{
	return ce.value
}


//Insert adds a new element to Cache, forming the hashmap entry from CacheElementRoot and CacheElementKey
func (c *Cache) Insert(iCacheRoot CacheElementRoot,iCacheKey CacheElementKey, iValue interface{}){
	c.Lock()
	defer c.Unlock()

	elementDuration,ok := c.cacheRootElementsDuration[iCacheRoot]
	if !ok{
		elementDuration = standardCacheElementDuration
	}

	c.elements[generateCacheKey(iCacheRoot,iCacheKey)] = cacheElement{
		deadline:time.Now().Add(elementDuration),
		value:iValue,
	}
}
//Fetch return the element associated with the hashmap entry formed by CacheElementRoot and CacheElementKey.
//Returns nil if the entry is not found.
func (c *Cache) Fetch(iCacheRoot CacheElementRoot,iCacheKey CacheElementKey) (oValue interface{}){
	c.RLock()
	defer c.RUnlock()

	cachedElement, ok := c.elements[generateCacheKey(iCacheRoot,iCacheKey)]
	if !ok{
		return nil
	}
	return cachedElement.value
}
//Reset empties the whole cache
func (c *Cache) Reset(){
	c.Lock()
	defer c.Unlock()
	c.elements = make(map[string]cacheElement)
	c.cacheRootElementsDuration = make(map[CacheElementRoot]time.Duration)
}
//generateCacheKey uses CacheElementRoot and CacheElementKey to generate a hashmap entry to a Cache
func generateCacheKey(iCacheRoot CacheElementRoot,iCacheKey CacheElementKey) string{
	//TODO string concatenation is slow, maybe use a buffer?
	return iCacheRoot.String()+"_"+iCacheKey.String()
}
