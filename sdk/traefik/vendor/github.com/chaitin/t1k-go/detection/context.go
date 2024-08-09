package detection

import (
	"errors"
	"net/http"

	"github.com/chaitin/t1k-go/misc"
)

type DetectionContext struct {
	UUID         string
	Scheme       string
	ProxyName    string
	RemoteAddr   string
	Protocol     string
	RemotePort   uint16
	LocalAddr    string
	LocalPort    uint16
	ReqBeginTime int64
	RspBeginTime int64

	T1KContext []byte

	Request  Request
	Response Response
}

func New() *DetectionContext {
	return &DetectionContext{
		UUID:       misc.GenUUID(),
		Scheme:     "http",
		ProxyName:  "go-sdk",
		RemoteAddr: "127.0.0.1",
		RemotePort: 30001,
		LocalAddr:  "127.0.0.1",
		LocalPort:  80,
		Protocol:   "HTTP/1.1",
	}
}

func MakeContextWithRequest(req *http.Request) (*DetectionContext, error) {
	if req == nil {
		return nil, errors.New("nil http.request or response")
	}
	wrapReq := &HttpRequest{
		req: req,
	}

	// ignore GetRemoteIP error,not sure request record remote ip
	remoteIP, _ := wrapReq.GetRemoteIP()
	remotePort, _ := wrapReq.GetRemotePort()

	localAddr, err := wrapReq.GetUpstreamAddress()
	if err != nil {
		return nil, err
	}

	localPort, err := wrapReq.GetUpstreamPort()
	if err != nil {
		return nil, err
	}

	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}

	context := &DetectionContext{
		UUID:         misc.GenUUID(),
		Scheme:       scheme,
		ProxyName:    "go-sdk",
		RemoteAddr:   remoteIP,
		RemotePort:   remotePort,
		LocalAddr:    localAddr,
		LocalPort:    localPort,
		ReqBeginTime: misc.Now(),
		Request:      wrapReq,
		Protocol:     req.Proto,
	}
	wrapReq.dc = context
	return context, nil
}

func (dc *DetectionContext) ProcessResult(r *Result) {
	if r.Objective == RO_REQUEST {
		dc.T1KContext = r.T1KContext
	}
}
