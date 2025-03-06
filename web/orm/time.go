package orm

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	defaultTimeFormat = "2006-01-02 15:04:05"
	timeLocalZone     = "Asia/Shanghai"
)

func SetTimeFormat(format string) {
	defaultTimeFormat = format
}

func SetTimeLocalZone(zone string) {
	timeLocalZone = zone
}

type LocalTime time.Time

func (t *LocalTime) Scan(v any) error {
	if val, ok := v.(time.Time); ok {
		*t = LocalTime(val)
		return nil
	}
	return fmt.Errorf("can not conver %v to timestamp", v)
}

func (t LocalTime) Value() (driver.Value, error) {
	tlt := time.Time(t)
	if tlt.IsZero() {
		return nil, nil
	}
	return tlt, nil
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tlt := time.Time(t)
	if tlt.IsZero() {
		return []byte("null"), nil
	}

	b := make([]byte, 0, len(defaultTimeFormat)+2)
	b = append(b, '"')
	b = tlt.AppendFormat(b, defaultTimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"`+defaultTimeFormat+`"`, string(data), time.Local)
	*t = LocalTime(now)
	return
}

func (t LocalTime) String() string {
	return time.Time(t).Local().Format(defaultTimeFormat)
}
