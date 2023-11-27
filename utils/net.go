package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

func SpliceUrl(addr string, port int, request string) string {
	return "http://" + addr + ":" + strconv.Itoa(port) + "/" + request
}
