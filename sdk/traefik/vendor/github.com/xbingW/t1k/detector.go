package t1k

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/xbingW/t1k/pkg/datetime"
	"github.com/xbingW/t1k/pkg/rand"
	"github.com/xbingW/t1k/pkg/t1k"
)

type Detector struct {
	cfg Config
}

type Config struct {
	Addr string `json:"addr"`
	// Get ip from header, if not set, get ip from remote addr
	IpHeader string `json:"ip_header"`
	// When ip_header has multiple ip, use this to get ip from right
	//
	//for example, X-Forwarded-For: ip1, ip2, ip3
	// 	when ip_last_index is 0, the client ip is ip3
	// 	when ip_last_index is 1, the client ip is ip2
	// 	when ip_last_index is 2, the client ip is ip1
	IPRightIndex uint `json:"ip_right_index"`
}

func NewDetector(cfg Config) *Detector {
	return &Detector{
		cfg: cfg,
	}
}

func (d *Detector) GetConn() (net.Conn, error) {
	_, _, err := net.SplitHostPort(d.cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("detector add %s is invalid because %v", d.cfg.Addr, err)
	}
	return net.Dial("tcp", d.cfg.Addr)
}

func (d *Detector) DetectorRequestStr(req string) (*t1k.DetectorResponse, error) {
	httpReq, err := http.ReadRequest(bufio.NewReader(strings.NewReader(req)))
	if err != nil {
		return nil, fmt.Errorf("read request failed: %v", err)
	}
	return d.DetectorRequest(httpReq)
}

func (d *Detector) DetectorRequest(req *http.Request) (*t1k.DetectorResponse, error) {
	extra, err := d.GenerateExtra(req)
	if err != nil {
		return nil, fmt.Errorf("generate extra failed: %v", err)
	}
	dc := t1k.NewHttpDetector(req, extra)
	conn, err := d.GetConn()
	if err != nil {
		return nil, fmt.Errorf("get conn failed: %v", err)
	}
	defer conn.Close()
	return dc.DetectRequest(conn)
}

func (d *Detector) DetectorResponseStr(req string, resp string) (*t1k.DetectorResponse, error) {
	httpReq, err := http.ReadRequest(bufio.NewReader(strings.NewReader(req)))
	if err != nil {
		return nil, fmt.Errorf("read request failed: %v", err)
	}
	httpResp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(resp)), httpReq)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %v", err)
	}
	extra, err := d.GenerateExtra(httpReq)
	if err != nil {
		return nil, fmt.Errorf("generate extra failed: %v", err)
	}
	conn, err := d.GetConn()
	if err != nil {
		return nil, fmt.Errorf("get conn failed: %v", err)
	}
	defer conn.Close()
	return t1k.NewHttpDetector(httpReq, extra).SetResponse(httpResp).DetectResponse(conn)
}

func (d *Detector) DetectorResponse(req *http.Request, resp *http.Response) (*t1k.DetectorResponse, error) {
	extra, err := d.GenerateExtra(req)
	if err != nil {
		return nil, fmt.Errorf("generate extra failed: %v", err)
	}
	conn, err := d.GetConn()
	if err != nil {
		return nil, fmt.Errorf("get conn failed: %v", err)
	}
	defer conn.Close()
	return t1k.NewHttpDetector(req, extra).SetResponse(resp).DetectResponse(conn)
}

func (d *Detector) GenerateExtra(req *http.Request) (*t1k.HttpExtra, error) {
	clientHost, err := d.getClientIP(req)
	if err != nil {
		return nil, err
	}
	serverHost, serverPort := req.Host, "80"
	if hasPort(req.Host) {
		serverHost, serverPort, err = net.SplitHostPort(req.Host)
		if err != nil {
			return nil, err
		}
	}
	return &t1k.HttpExtra{
		UpstreamAddr:  "",
		RemoteAddr:    clientHost,
		RemotePort:    d.getClientPort(req),
		LocalAddr:     serverHost,
		LocalPort:     serverPort,
		ServerName:    "",
		Schema:        req.URL.Scheme,
		ProxyName:     "",
		UUID:          rand.String(32),
		HasRspIfOK:    "y",
		HasRspIfBlock: "n",
		ReqBeginTime:  strconv.FormatInt(datetime.Now(), 10),
		ReqEndTime:    "",
		RspBeginTime:  strconv.FormatInt(datetime.Now(), 10),
		RepEndTime:    "",
	}, nil
}

func (d *Detector) getClientIP(req *http.Request) (string, error) {
	if d.cfg.IpHeader != "" {
		ips := req.Header.Get(d.cfg.IpHeader)
		if ips != "" {
			ipList := reverseStrSlice(strings.Split(ips, ","))
			if len(ipList) > int(d.cfg.IPRightIndex) {
				return strings.TrimSpace(ipList[d.cfg.IPRightIndex]), nil
			}
			return ipList[0], nil
		}
	}
	if !hasPort(req.RemoteAddr) {
		return req.RemoteAddr, nil
	}
	clientHost, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", err
	}
	return clientHost, nil
}

func (d *Detector) getClientPort(req *http.Request) string {
	_, clientPort, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return ""
	}
	return clientPort
}

// has port check if host has port
func hasPort(host string) bool {
	return strings.LastIndex(host, ":") > strings.LastIndex(host, "]")
}

func reverseStrSlice(arr []string) []string {
	if len(arr) == 0 {
		return arr
	}
	return append(reverseStrSlice(arr[1:]), arr[0])
}
