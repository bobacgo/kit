package util

import (
	"fmt"
	"testing"
)

func TestRandSeq(t *testing.T) {
	seq := RandSeqID(32)
	for i := 0; i < 10; i++ {
		fmt.Println(seq())
	}
}

func TestUUID(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(UUID())
	}
}

func TestRandNumber(t *testing.T) {
	r := RandNumber(100, 1000)
	for i := 0; i < 10; i++ {
		fmt.Println(r())
	}
}
