package main

import (
	"log"
	"testing"

	"github.com/coocood/freecache"
)

func TestPingHandler(t *testing.T) {
	key := getRandNum()
	v, err := cacheInstance.Get([]byte(key))
	if err != nil && err.Error() != freecache.ErrNotFound.Error() {
		log.Println("freecache get error: ", err)
		return
	}
	if len(v) == 0 {
		if err := cacheInstance.Set([]byte(key), cacheValue, 2); err != nil {
			log.Println("ser error: ", err)
		}
	}
}
