package detection

import (
	"fmt"

	"github.com/chaitin/t1k-go/misc"
)

func MakeRequestExtra(
	scheme string,
	proxyName string,
	remoteAddr string,
	remotePort uint16,
	localAddr string,
	localPort uint16,
	uuid string,
	hasRspIfOK string,
	hasRspIfBlock string,
	reqBeginTime int64,
) []byte {
	format := "Scheme:%s\n" +
		"ProxyName:%s\n" +
		"RemoteAddr:%s\n" +
		"RemotePort:%d\n" +
		"LocalAddr:%s\n" +
		"LocalPort:%d\n" +
		"UUID:%s\n" +
		"HasRspIfOK:%s\n" +
		"HasRspIfBlock:%s\n" +
		"ReqBeginTime:%d\n"

	return []byte(fmt.Sprintf(
		format,
		scheme,
		proxyName,
		remoteAddr,
		remotePort,
		localAddr,
		localPort,
		uuid,
		hasRspIfOK,
		hasRspIfBlock,
		reqBeginTime,
	))
}

func MakeResponseExtra(
	scheme string,
	proxyName string,
	remoteAddr string,
	remotePort uint16,
	localAddr string,
	localPort uint16,
	uuid string,
	rspBeginTime int64,
) []byte {
	format := "Scheme:%s\n" +
		"ProxyName:%s\n" +
		"RemoteAddr:%s\n" +
		"RemotePort:%d\n" +
		"LocalAddr:%s\n" +
		"LocalPort:%d\n" +
		"UUID:%s\n" +
		"RspBeginTime:%d\n"

	return []byte(fmt.Sprintf(
		format,
		scheme,
		proxyName,
		remoteAddr,
		remotePort,
		localAddr,
		localPort,
		uuid,
		rspBeginTime,
	))
}

func PlaceholderRequestExtra(uuid string) []byte {
	return MakeRequestExtra("http", "go-sdk", "127.0.0.1", 30001, "127.0.0.1", 80, uuid, "n", "n", misc.Now())
}

func GenRequestExtra(dc *DetectionContext) []byte {
	hasRsp := "u"
	if dc.Response != nil {
		hasRsp = "y"
	}
	return MakeRequestExtra(dc.Scheme, dc.ProxyName, dc.RemoteAddr, dc.RemotePort, dc.LocalAddr, dc.LocalPort, dc.UUID, hasRsp, hasRsp, dc.ReqBeginTime)
}

func GenResponseExtra(dc *DetectionContext) []byte {
	return MakeResponseExtra(dc.Scheme, dc.ProxyName, dc.RemoteAddr, dc.RemotePort, dc.LocalAddr, dc.LocalPort, dc.UUID, dc.RspBeginTime)
}
