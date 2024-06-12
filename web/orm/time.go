package orm

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/bobacgo/kit/g"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
	TimeLocalZone     = "Asia/Shanghai"
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

	b := make([]byte, 0, len(t.Format())+2)
	b = append(b, '"')
	b = tlt.AppendFormat(b, t.Format())
	b = append(b, '"')
	return b, nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"`+t.Format()+`"`, string(data), time.Local)
	*t = LocalTime(now)
	return
}

func (t LocalTime) String() string {
	return time.Time(t).Local().Format(t.Format())
}

// Format 从配置文件获取时间格式
func (t LocalTime) Format() string {
	tf := DefaultTimeFormat
	if g.Conf != nil && g.Conf.Logger.TimeFormat != "" {
		tf = g.Conf.Logger.TimeFormat
	}
	return tf
}
