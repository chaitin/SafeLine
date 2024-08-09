package misc

import (
	"encoding/hex"
)

var (
	rng *MT19937 = NewMT19937()
)

func GenUUID() string {
	u := make([]byte, 16)
	rng.RandBytes(u)
	u[6] = (u[6] | 0x40) & 0x4F
	u[8] = (u[8] | 0x80) & 0xBF
	return hex.EncodeToString(u)
}
