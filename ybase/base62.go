package ybase

import "math/big"

func EncodeBase62(b []byte) string {
	var i big.Int
	i.SetBytes(b[:])
	return i.Text(62)
}

func DecodeBase62(s string) []byte {
	var i big.Int
	i.SetString(s, 62)
	return i.Bytes()
}
