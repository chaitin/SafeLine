package t1k

import (
	"bytes"
	"io"

	"github.com/chaitin/t1k-go/misc"
)

type Section interface {
	Header() Header
	WriteBody(io.Writer) error
}

type SimpleSection struct {
	tag  Tag
	body []byte
}

func MakeSimpleSection(tag Tag, body []byte) *SimpleSection {
	return &SimpleSection{
		tag:  tag,
		body: body,
	}
}

func (msg *SimpleSection) Header() Header {
	return MakeHeader(msg.tag, uint32(len(msg.body)))
}

func (msg *SimpleSection) WriteBody(w io.Writer) error {
	_, err := w.Write(msg.body)
	return err
}

type ReaderSection struct {
	tag    Tag
	size   uint32
	reader io.Reader
}

func MakeReaderSection(tag Tag, size uint32, reader io.Reader) *ReaderSection {
	return &ReaderSection{
		tag:    tag,
		size:   size,
		reader: reader,
	}
}

func (msg *ReaderSection) Header() Header {
	return MakeHeader(msg.tag, msg.size)
}

func (msg *ReaderSection) WriteBody(w io.Writer) error {
	_, err := io.CopyN(w, msg.reader, int64(msg.size))
	return misc.ErrorWrap(err, "")
}

func WriteSection(s Section, w io.Writer) error {
	h := s.Header()
	_, err := w.Write(h.Serialize())
	if err != nil {
		return misc.ErrorWrap(err, "")
	}
	return s.WriteBody(w)
}

// returns a *ReaderSection, must call its WriteBody
// before next call to r
func ReadSection(r io.Reader) (Section, error) {
	bHeader := make([]byte, T1K_HEADER_SIZE)
	_, err := io.ReadFull(r, bHeader)
	if err != nil {
		return nil, err
	}
	h := DeserializeHeader(bHeader)
	bodyReader := io.LimitReader(r, int64(h.Size))
	return MakeReaderSection(h.Tag, h.Size, bodyReader), nil
}

// returns a *SimpleSection
func ReadFullSection(r io.Reader) (Section, error) {
	bHeader := make([]byte, T1K_HEADER_SIZE)
	_, err := io.ReadFull(r, bHeader)
	if err != nil {
		return nil, misc.ErrorWrap(err, "")
	}
	h := DeserializeHeader(bHeader)
	var buf bytes.Buffer
	_, err = io.CopyN(&buf, r, int64(h.Size))
	if err != nil {
		return nil, misc.ErrorWrap(err, "")
	}
	return MakeSimpleSection(h.Tag, buf.Bytes()), nil
}
