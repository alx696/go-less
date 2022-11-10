package lilu_net

import (
	"fmt"
	"log"
	"net"
	"sort"
	"time"
)

// GetFreePort 获取可用端口
func GetFreePort() (port int, e error) {
	var a *net.TCPAddr
	a, e = net.ResolveTCPAddr("tcp", "localhost:0")
	if e != nil {
		return 0, e
	}

	var l *net.TCPListener
	l, e = net.ListenTCP("tcp", a)
	if e != nil {
		return 0, e
	}
	defer func() {
		_ = l.Close()
	}()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// CheckPortFree 检查端口是否可用
func CheckPortFree(port int) bool {
	conn, e := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprint(port)), time.Second)
	if e != nil {
		return true
	}
	conn.Close()
	return false
}

// GetIp 获取局域网IPv4或公网IPv6
//
// 注意: Android的targetSdk设为30以上时Golang会出网络异常 https://github.com/ipfs-shipyard/gomobile-ipfs/issues/68 .
func GetIp() (string, error) {
	addrs, e := net.InterfaceAddrs()
	if e != nil {
		return "", e
	}

	// 提取可用IP
	type ipInfo struct {
		Ip    string
		Value int
	}
	var ipInfoArray []ipInfo
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() || ipnet.IP.IsLinkLocalUnicast() {
			continue
		}

		ipValue := 10
		_, block1, _ := net.ParseCIDR("192.168.0.0/16")
		_, block2, _ := net.ParseCIDR("172.16.0.0/12")
		_, block3, _ := net.ParseCIDR("10.0.0.0/8")

		if block1.Contains(ipnet.IP) {
			ipValue = 1
		} else if block2.Contains(ipnet.IP) {
			ipValue = 2
		} else if block3.Contains(ipnet.IP) {
			ipValue = 3
		}

		ipInfoArray = append(ipInfoArray, ipInfo{Ip: ipnet.IP.String(), Value: ipValue})
	}

	// 排序
	sort.SliceStable(ipInfoArray, func(i, j int) bool {
		return ipInfoArray[i].Value < ipInfoArray[j].Value
	})

	log.Println("可用IP", ipInfoArray)

	if len(ipInfoArray) == 0 {
		return "", fmt.Errorf("没有可用IP")
	}

	return ipInfoArray[0].Ip, nil
}

// GetUsedIp 获取在用IP(需要互联网)
func GetUsedIp() (string, error) {
	conn, e := net.DialTimeout("udp", "119.29.29.29:80", time.Second)
	if e != nil {
		return "", fmt.Errorf("无法访问: %s", e.Error())
	}
	defer conn.Close()

	addr := conn.LocalAddr().(*net.UDPAddr)
	return addr.IP.String(), nil
}
