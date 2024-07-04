package t1k

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Request interface {
	Header() ([]byte, error)
	Body() ([]byte, error)
	Extra() ([]byte, error)
	Serialize() ([]byte, error)
}

type HttpRequest struct {
	extra *HttpExtra
	req   *http.Request
}

func NewHttpRequest(req *http.Request, extra *HttpExtra) *HttpRequest {
	return &HttpRequest{
		req:   req,
		extra: extra,
	}
}

func NewHttpRequestRead(req string) *HttpRequest {
	httpReq, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(req)))
	return &HttpRequest{
		req: httpReq,
	}
}

func (r *HttpRequest) Header() ([]byte, error) {
	var buf bytes.Buffer
	proto := r.req.Proto
	startLine := fmt.Sprintf("%s %s %s\r\n", r.req.Method, r.req.URL.RequestURI(), proto)
	_, err := buf.Write([]byte(startLine))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write([]byte(fmt.Sprintf("Host: %s\r\n", r.req.Host)))
	if err != nil {
		return nil, err
	}
	err = r.req.Header.Write(&buf)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write([]byte("\r\n"))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *HttpRequest) Body() ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.req.Body)
	if err != nil {
		return nil, err
	}
	r.req.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	return buf.Bytes(), nil
}

func (r *HttpRequest) Extra() ([]byte, error) {
	return r.extra.ReqSerialize()
}

func (r *HttpRequest) Version() []byte {
	return []byte("Proto:3\n")
}

func (r *HttpRequest) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	{
		raw, err := r.Header()
		if err != nil {
			return nil, err
		}
		packet := NewHttpPacket(TAG_HEADER|MASK_FIRST, raw)
		_, err = buf.Write(packet.Serialize())
		if err != nil {
			return nil, err
		}
	}
	{
		raw, err := r.Body()
		if err != nil {
			return nil, err
		}
		packet := NewHttpPacket(TAG_BODY, raw)
		_, err = buf.Write(packet.Serialize())
		if err != nil {
			return nil, err
		}
	}
	{
		raw, err := r.Extra()
		if err != nil {
			return nil, err
		}
		packet := NewHttpPacket(TAG_EXTRA, raw)
		_, err = buf.Write(packet.Serialize())
		if err != nil {
			return nil, err
		}
	}
	{
		packet := NewHttpPacket(TAG_VERSION|MASK_LAST, r.Version())
		_, err := buf.Write(packet.Serialize())
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
