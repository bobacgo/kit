package errs

import (
	"github.com/bobacgo/kit/web/r/codes"
	"github.com/bobacgo/kit/web/r/status"
)

var (
	BadRequest    = status.New(codes.BadRequest, "请求参数错误")
	InternalError = status.New(codes.InternalServerError, "服务器繁忙")
	RecordRepeat  = status.New(codes.RecordRepeat, "数据已经存在")
	DateBusy      = status.New(codes.DateBusy, "数据在使用")
)
