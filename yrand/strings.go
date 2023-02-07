package yrand

import (
	"crypto/rand"
	"io"
	"math/big"
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

	var i big.Int
	i.SetBytes(b[:])
	return i.Text(62)
}
