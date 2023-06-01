package util

import (
	"fmt"
	"testing"
)

func TestExample(t *testing.T) {
	hash, salt := BcryptHash("gogo")
	fmt.Println(salt, hash)
	verify := BcryptVerify(salt, hash, "gogo")
	t.Logf("verify: %v\n", verify)
}
