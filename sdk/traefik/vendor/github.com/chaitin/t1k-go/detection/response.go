package detection

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/chaitin/t1k-go/misc"
)

type Response interface {
	RequestHeader() ([]byte, error)
	Header() ([]byte, error)
	Body() (uint32, io.ReadCloser, error)
	Extra() ([]byte, error)
	T1KContext() ([]byte, error)
}

type HttpResponse struct {
	rsp *http.Response
	dc  *DetectionContext // this is a must-have
}

func MakeHttpResponseInCtx(rsp *http.Response, dc *DetectionContext) *HttpResponse {
	ret := &HttpResponse{
		rsp: rsp,
		dc:  dc,
	}
	dc.Response = ret
	dc.RspBeginTime = misc.Now()
	return ret
}

func (r *HttpResponse) RequestHeader() ([]byte, error) {
	return r.dc.Request.Header()
}

func (r *HttpResponse) Header() ([]byte, error) {
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

func (r *HttpResponse) Body() (uint32, io.ReadCloser, error) {
	bodyBytes, err := io.ReadAll(r.rsp.Body)
	if err != nil {
		return 0, nil, misc.ErrorWrapf(err, "get body size %d", len(bodyBytes))
	}
	r.rsp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return uint32(len(bodyBytes)), io.NopCloser(bytes.NewReader(bodyBytes)), nil
}

func (r *HttpResponse) Extra() ([]byte, error) {
	return GenResponseExtra(r.dc), nil
}

func (r *HttpResponse) T1KContext() ([]byte, error) {
	return r.dc.T1KContext, nil
}
