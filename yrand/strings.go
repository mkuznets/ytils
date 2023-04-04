package yrand

import (
	"crypto/rand"
	"io"
	"mkuznets.com/go/ytils/ybase"
	"sync"
)

var randMutex = sync.Mutex{}

func Base62(nBytes int) string {
	b := make([]byte, nBytes)
	randMutex.Lock()
	_, err := io.ReadFull(rand.Reader, b)
	randMutex.Unlock()
	if err != nil {
		panic("failed to read random bytes: " + err.Error())
	}
	return ybase.EncodeBase62(b)
}
