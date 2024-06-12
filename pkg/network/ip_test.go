package network_test

import (
	"testing"

	"github.com/bobacgo/kit/pkg/network"
)

func TestGetOutBoundIP(t *testing.T) {
	ip, err := network.OutBoundIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip)
}
