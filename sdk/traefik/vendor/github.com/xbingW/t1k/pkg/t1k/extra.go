package t1k

import "fmt"

type HttpExtra struct {
	UpstreamAddr  string
	RemoteAddr    string
	RemotePort    string
	LocalAddr     string
	LocalPort     string
	ServerName    string
	Schema        string
	ProxyName     string
	UUID          string
	HasRspIfOK    string
	HasRspIfBlock string
	ReqBeginTime  string
	ReqEndTime    string
	RspBeginTime  string
	RepEndTime    string
}

func (h *HttpExtra) ReqSerialize() ([]byte, error) {
	format := "UpstreamAddr:%s\n" +
		"RemotePort:%s\n" +
		"LocalPort:%s\n" +
		"RemoteAddr:%s\n" +
		"LocalAddr:%s\n" +
		"ServerName:%s\n" +
		"Schema:%s\n" +
		"ProxyName:%s\n" +
		"UUID:%s\n" +
		"HasRspIfOK:%s\n" +
		"HasRspIfBlock:%s\n" +
		"ReqBeginTime:%s\n" +
		"ReqEndTime:%s\n"
	return []byte(fmt.Sprintf(
		format,
		h.UpstreamAddr,
		h.RemotePort,
		h.LocalPort,
		h.RemoteAddr,
		h.LocalAddr,
		h.ServerName,
		h.Schema,
		h.ProxyName,
		h.UUID,
		h.HasRspIfOK,
		h.HasRspIfBlock,
		h.ReqBeginTime,
		h.ReqEndTime,
	)), nil
}

func (h *HttpExtra) RspSerialize() ([]byte, error) {
	format := "Scheme:%s\n" +
		"ProxyName:%s\n" +
		"RemoteAddr:%s\n" +
		"RemotePort:%s\n" +
		"LocalAddr:%s\n" +
		"LocalPort:%s\n" +
		"UUID:%s\n" +
		"RspBeginTime:%s\n"
	return []byte(fmt.Sprintf(
		format,
		h.Schema,
		h.ProxyName,
		h.RemoteAddr,
		h.RemotePort,
		h.LocalAddr,
		h.LocalPort,
		h.UUID,
		h.RspBeginTime,
	)), nil
}
