package uid_test

import (
	"testing"

	"github.com/bobacgo/kit/pkg/uid"
)

func TestRandSeqID(t *testing.T) {
	randString := uid.RandSeqID(16)
	for i := 0; i < 10; i++ {
		t.Log(randString())
	}
}
