package util

import (
	"fmt"
	"testing"
	"time"
)

func TestRandSeq(t *testing.T) {
	seq := RandSeqID(32)
	for i := 0; i < 10; i++ {
		fmt.Println(seq())
	}
}

func TestSnowflakeID(t *testing.T) {
	sf := Snowflake(time.Now(), 2)
	for i := 0; i < 10; i++ {
		fmt.Println(sf.NextID())
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
