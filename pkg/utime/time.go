package utime

import "time"

// ZeroHour 后几天零点
func ZeroHour(day int) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day()+day, 0, 0, 0, 0, now.Location())
}
