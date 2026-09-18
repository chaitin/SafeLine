package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStr returns n letters drawn from a cryptographically secure source.
//
// The value is used to name the temporary file created by RenameWriteFile, so
// it must not be predictable: math/rand is seeded from the clock and its
// output can be reconstructed from a single sample.
func RandStr(n int) string {
	b := make([]rune, n)
	limit := big.NewInt(int64(len(letters)))
	for i := range b {
		idx, err := rand.Int(rand.Reader, limit)
		if err != nil {
			panic("utils: cannot read from the system entropy source: " + err.Error())
		}
		b[i] = letters[idx.Int64()]
	}
	return string(b)
}

// RandHex returns a hex string of 2*n characters read from crypto/rand.
func RandHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("utils: cannot read from the system entropy source: " + err.Error())
	}
	return hex.EncodeToString(b)
}
