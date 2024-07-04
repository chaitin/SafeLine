package t1k

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Response interface {
	RequestHeader() ([]byte, error)
	RspHeader() ([]byte, error)
	Body() ([]byte, error)
	Extra() ([]byte, error)
	Serialize() ([]byte, error)
}

type HttpResponse struct {
	req   Request
	extra *HttpExtra
	rsp   *http.Response
}

func NewHttpResponse(req Request, rsp *http.Response, extra *HttpExtra) *HttpResponse {
	return &HttpResponse{
		req:   req,
		rsp:   rsp,
		extra: extra,
	}
}

func (r *HttpResponse) RequestHeader() ([]byte, error) {
	return r.req.Header()
}

func (r *HttpResponse) RspHeader() ([]byte, error) {
	var buf bytes.Buffer
	statusLine := fmt.Sprintf("HTTP/1.1 %s\n", r.rsp.Status)
	_, err := buf.Write([]byte(statusLine))
	if err != nil {
		return nil, err
	}
	err = r.rsp.Header.Write(&buf)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write([]byte("\r\n"))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *HttpResponse) Body() ([]byte, error) {
	data, err := io.ReadAll(r.rsp.Body)
	if err != nil {
		return nil, err
	}
	r.rsp.Body = io.NopCloser(bytes.NewReader(data))
	return data, nil
}

func (r *HttpResponse) Extra() ([]byte, error) {
	return r.extra.RspSerialize()
}

func (r *HttpResponse) Version() []byte {
	return []byte("Proto:3\n")

}

func (r *HttpResponse) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	{
		header, err := r.RequestHeader()
		if err != nil {
			return nil, err
		}
		packet := NewHttpPacket(TAG_HEADER|MASK_FIRST, header)
		_, err = buf.Write(packet.Serialize())
		if err != nil {
			return nil, err
		}
	}
	{
		header, err := r.RspHeader()
		if err != nil {
			return nil, err
		}
		packet := NewHttpPacket(TAG_RSP_HEADER, header)
		_, err = buf.Write(packet.Serialize())
		if err != nil {
			return nil, err
		}
	}
	{
		body, err := r.Body()
		if err != nil {
			return nil, err
		}
		packet := NewHttpPacket(TAG_RSP_BODY, body)
		_, err = buf.Write(packet.Serialize())
		if err != nil {
			return nil, err
		}
	}
	{
		extra, err := r.Extra()
		if err != nil {
			return nil, err
		}
		packet := NewHttpPacket(TAG_RSP_EXTRA, extra)
		_, err = buf.Write(packet.Serialize())
		if err != nil {
			return nil, err
		}
	}
	{
		version := r.Version()
		packet := NewHttpPacket(TAG_VERSION|MASK_LAST, version)
		_, err := buf.Write(packet.Serialize())
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
