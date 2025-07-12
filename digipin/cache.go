package digipin

import (
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
)

var (
	cacheOnce   sync.Once
	encodeCache *lru.Cache[string, string]
	decodeCache *lru.Cache[string, LatLng]
)

func initCaches() {
	cacheOnce.Do(func() {
		encodeCache, _ = lru.New[string, string](10000)
		decodeCache, _ = lru.New[string, LatLng](10000)
	})
}

// GetCachedEncode returns cached DIGIPIN for lat/lng if present
func GetCachedEncode(lat, lng float64) (string, bool) {
	initCaches()
	key := fmt.Sprintf("%f,%f", lat, lng)
	val, ok := encodeCache.Get(key)
	return val, ok
}

// SetCachedEncode stores DIGIPIN for lat/lng
func SetCachedEncode(lat, lng float64, pin string) {
	initCaches()
	key := fmt.Sprintf("%f,%f", lat, lng)
	encodeCache.Add(key, pin)
}

// GetCachedDecode returns cached LatLng for DIGIPIN if present
func GetCachedDecode(pin string) (LatLng, bool) {
	initCaches()
	val, ok := decodeCache.Get(pin)
	return val, ok
}

// SetCachedDecode stores LatLng for DIGIPIN
func SetCachedDecode(pin string, coord LatLng) {
	initCaches()
	decodeCache.Add(pin, coord)
}

// Clear all caches
func ClearCache() {
	initCaches()
	encodeCache.Purge()
	decodeCache.Purge()
}
