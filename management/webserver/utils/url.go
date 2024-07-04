package utils

import (
	"fmt"
	"net"
	"strings"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
)

func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}

func BuildUrl(protocol int, host string, port uint, path string) string {
	portStr := ""
	if len(host) > 0 && IsIPv6(host) && host[:1] != "[" {
		host = fmt.Sprintf("[%s]", host)
	}
	if protocol == constants.ProtocolHTTP {
		if port != 80 {
			portStr = fmt.Sprintf(":%d", port)
		}
	} else if protocol == constants.ProtocolHTTPS || protocol == constants.ProtocolHTTP2 {
		if port != 443 {
			portStr = fmt.Sprintf(":%d", port)
		}
	} else {
		// use HTTP as default protocol
		protocol = constants.ProtocolHTTP
	}
	return fmt.Sprintf("%s://%s%s%s", constants.HTTPProtocol[protocol], host, portStr, path)
}
