package t1k

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Packet interface {
	Serialize() []byte
	Last() bool
	Tag() Tag
	PayLoad() []byte
}

type HttpPacket struct {
	tag     Tag
	payload []byte
}

func (p *HttpPacket) Last() bool {
	return p.tag.Last()
}

func (p *HttpPacket) Tag() Tag {
	return p.tag
}

func (p *HttpPacket) PayLoad() []byte {
	return p.payload
}

func NewHttpPacket(tag Tag, payload []byte) Packet {
	return &HttpPacket{
		tag:     tag,
		payload: payload,
	}
}

func (p *HttpPacket) SizeBytes() []byte {
	sizeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizeBytes, uint32(len(p.payload)))
	return sizeBytes
}

func (p *HttpPacket) Serialize() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(p.tag))
	buf.Write(p.SizeBytes())
	buf.Write(p.payload)
	return buf.Bytes()
}

func ReadPacket(r io.Reader) (Packet, error) {
	tag := make([]byte, 1)
	if _, err := io.ReadFull(r, tag); err != nil {
		return nil, err
	}
	sizeBytes := make([]byte, 4)
	if _, err := io.ReadFull(r, sizeBytes); err != nil {
		return nil, err
	}
	size := binary.LittleEndian.Uint32(sizeBytes)
	payload := make([]byte, size)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}
	return NewHttpPacket(Tag(tag[0]), payload), nil
}
