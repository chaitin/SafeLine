package misc

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
	"time"
)

const (
	n               = 312
	m               = 156
	seedMask        = 170298
	hiMask   uint64 = 0xffffffff80000000
	loMask   uint64 = 0x000000007fffffff
	matrixA  uint64 = 0xB5026F5AA96619E9
)

type MT19937 struct {
	state []uint64
	index int
	mut   sync.Mutex
}

func NewMT19937WithSeed(seed int64) *MT19937 {
	mt := &MT19937{
		state: make([]uint64, n),
		index: n,
	}

	mt.mut.Lock()
	defer mt.mut.Unlock()

	mt.state[0] = uint64(seed)
	for i := uint64(1); i < n; i++ {
		mt.state[i] = 6364136223846793005*(mt.state[i-1]^(mt.state[i-1]>>62)) + i
	}

	return mt
}

func NewMT19937() *MT19937 {
	var seed int64

	b := make([]byte, 8)
	if _, err := rand.Read(b); err == nil {
		seed = int64(binary.LittleEndian.Uint64(b[:]))
	}

	seed = seed ^ time.Now().UnixNano() ^ seedMask

	return NewMT19937WithSeed(seed)
}

func (mt *MT19937) Uint64() uint64 {
	mt.mut.Lock()
	defer mt.mut.Unlock()

	x := mt.state
	if mt.index >= n {
		for i := 0; i < n-m; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+m] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		for i := n - m; i < n-1; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+(m-n)] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		y := (x[n-1] & hiMask) | (x[0] & loMask)
		x[n-1] = x[m-1] ^ (y >> 1) ^ ((y & 1) * matrixA)
		mt.index = 0
	}
	y := x[mt.index]
	y ^= (y >> 29) & 0x5555555555555555
	y ^= (y << 17) & 0x71D67FFFEDA60000
	y ^= (y << 37) & 0xFFF7EEE000000000
	y ^= (y >> 43)
	mt.index++
	return y
}

func (mt *MT19937) RandBytes(p []byte) {
	for len(p) >= 8 {
		val := mt.Uint64()
		p[0] = byte(val)
		p[1] = byte(val >> 8)
		p[2] = byte(val >> 16)
		p[3] = byte(val >> 24)
		p[4] = byte(val >> 32)
		p[5] = byte(val >> 40)
		p[6] = byte(val >> 48)
		p[7] = byte(val >> 56)
		p = p[8:]
	}
	if len(p) > 0 {
		val := mt.Uint64()
		for i := 0; i < len(p); i++ {
			p[i] = byte(val)
			val >>= 8
		}
	}
}
