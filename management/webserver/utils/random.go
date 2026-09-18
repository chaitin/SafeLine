package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStr returns n letters drawn from a cryptographically secure source.
//
// The values produced here are used as session signing keys, temporary file
// names and event ids. math/rand is seeded from a predictable value (and its
// output can be reconstructed from a single sample), so it must not be used
// for any of them.
func RandStr(n int) string {
	b := make([]rune, n)
	limit := big.NewInt(int64(len(letters)))
	for i := range b {
		idx, err := rand.Int(rand.Reader, limit)
		if err != nil {
			// crypto/rand fails only when the system entropy source is
			// unusable. Continuing with a predictable value is not an option,
			// because callers use the result as a secret.
			panic("utils: cannot read from the system entropy source: " + err.Error())
		}
		b[i] = letters[idx.Int64()]
	}
	return string(b)
}

// RandHex returns a hex string of 2*n characters read from crypto/rand. Every
// byte of the result carries eight bits of entropy, which makes it the right
// helper for keys and secrets.
func RandHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("utils: cannot read from the system entropy source: " + err.Error())
	}
	return hex.EncodeToString(b)
}
