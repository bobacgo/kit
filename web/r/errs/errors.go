package errs

import (
	"github.com/bobacgo/kit/web/r/codes"
	"github.com/bobacgo/kit/web/r/status"
)

var (
	BadRequest    = status.New(codes.BadRequest, "请求参数错误")
	InternalError = status.New(codes.InternalServerError, "服务器繁忙")
)
