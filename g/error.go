package g

import (
	"errors"
	"fmt"
	"runtime/debug"
)

var (
	ErrRecordRepeat = errors.New("数据已经存在")
	ErrDateBusy     = errors.New("数据在使用")
)

// Error 自定义错误
type Error struct {
	Inner      error
	Text       string // 业务错误信息，通常是中文提示，返回给客户端
	StackTrace string
	Misc       map[string]any // miscellaneous（各种各样） information
}

func NewError(message string) *Error {
	return NewErrorf(message)
}

func NewErrorf(format string, a ...any) *Error {
	return &Error{
		Inner:      nil,
		Text:       fmt.Sprintf(format, a...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]any),
	}
}

// WrapError new Error
func WrapError(err error, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Inner:      err,
		Text:       message,
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]any),
	}
}

// Unwrap 返回 inner error
func (err *Error) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Inner
}

// Is inner error 断言
func (err *Error) Is(target error) bool {
	return errors.Is(err.Unwrap(), target)
}

func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	msg := err.Text
	if err.Inner != nil {
		if msg != "" {
			msg += ": "
		}
		msg += err.Inner.Error()
	}
	return msg
}
