package security

import (
	"fmt"
	"log/slog"
	"testing"
)

func TestPhoneNo_LogValue(t *testing.T) {
	slog.Info("phoneNo", "phoneNo", PhoneNo("13800000000"))
	fmt.Println(PhoneNo("13800000000").LogValue())
}