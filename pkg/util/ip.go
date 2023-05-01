package util

import (
	"net"
	"strconv"
	"strings"
)

var IP = new(addr)

type addr struct{}

// GetOutBoundIP 获取主机本机IP
func (addr) GetOutBoundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(localAddr.String(), ":")[0]
	return ip, err
}

// Parse addr -> ip, port
func (addr) Parse(addr string) (string, uint16) {
	if addr == "" {
		return "", 0
	}
	ipAndPort := strings.Split(addr, ":")
	if len(ipAndPort) < 1 { // Incorrect format
		return "", 0
	}
	port, err := strconv.Atoi(ipAndPort[1])
	if err != nil {
		return "", 0
	}
	return ipAndPort[0], uint16(port)
}
