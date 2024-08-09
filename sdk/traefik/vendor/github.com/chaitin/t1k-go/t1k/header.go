package t1k

import (
	"encoding/binary"
)

const (
	T1K_HEADER_SIZE uint64 = 5
)

type Header struct {
	Tag  Tag
	Size uint32
}

func MakeHeader(tag Tag, size uint32) Header {
	return Header{
		Tag:  tag,
		Size: size,
	}
}

func (h Header) Serialize() []byte {
	b := make([]byte, 5)
	b[0] = byte(uint8(h.Tag))
	binary.LittleEndian.PutUint32(b[1:], h.Size)
	return b
}

func DeserializeHeader(b []byte) Header {
	return MakeHeader(
		Tag(uint8(b[0])),
		binary.LittleEndian.Uint32(b[1:]),
	)
}
