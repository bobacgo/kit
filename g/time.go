package g

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	TimeYYYYMMss  = "2006-01-02 15:04:05"
	TimeLocalZone = "Asia/Shanghai"
)

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

	b := make([]byte, 0, len(TimeYYYYMMss)+2)
	b = append(b, '"')
	b = tlt.AppendFormat(b, TimeYYYYMMss)
	b = append(b, '"')
	return b, nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(TimeYYYYMMss, string(data), time.Local)
	*t = LocalTime(now)
	return
}

func (t LocalTime) String() string {
	return time.Time(t).Local().Format(TimeYYYYMMss)
}
