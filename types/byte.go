package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 512m 512Mb

type ByteSize string

func (bs ByteSize) Int() int64 {
	size, _ := ParseByteUnit(string(bs))
	return size
}

func (bs ByteSize) Check() error {
	_, err := ParseByteUnit(string(bs))
	return err
}

// 将正则表达式编译为包级别的变量，避免每次调用时重复编译
var numRegexp = regexp.MustCompile("[0-9]+")

func ParseByteUnit(memSize string) (int64, error) {
	loc := numRegexp.FindStringIndex(memSize)
	if len(loc) != 2 {
		return 0, fmt.Errorf("unit parse not exist: %s len(loc) = %d", memSize, len(loc))
	}
	num, _ := strconv.ParseInt(memSize[:loc[1]], 10, 64)
	unit := strings.ToUpper(memSize[loc[1]:])
	switch unit {
	case "B":
		return num, nil
	case "KB", "K":
		return num << 10, nil
	case "MB", "M":
		return num << 20, nil
	case "GB", "G":
		return num << 30, nil
	case "TB", "T":
		return num << 40, nil
	case "PB", "P":
		return num << 50, nil
	default:
		return 0, fmt.Errorf("unit parse not exist: %s -> %s", memSize, memSize[loc[1]:])
	}
}
