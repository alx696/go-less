package lilu_net_test

import (
	"testing"

	"github.com/alx696/go-less/lilu_net"
)

func TestGetFreePort(t *testing.T) {
	t.Log(lilu_net.GetFreePort())
}

func TestCheckPortFree(t *testing.T) {
	t.Log(lilu_net.CheckPortFree(83))
}

func TestGetIp(t *testing.T) {
	t.Log(lilu_net.GetIp())
}

func TestGetUsedIp(t *testing.T) {
	t.Log(lilu_net.GetUsedIp())
}
