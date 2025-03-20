package r

import (
	"github.com/bobacgo/kit/web/r/errs"
	"github.com/gin-gonic/gin"

	"github.com/bobacgo/kit/app/validator"

	"net/http"

	"github.com/bobacgo/kit/web/r/codes"
	"github.com/bobacgo/kit/web/r/status"
	pkgvalidator "github.com/go-playground/validator/v10"
)

type Response[T any] struct {
	Code codes.Code `json:"code"`
	Data T          `json:"data"`
	Msg  string     `json:"message"`
	Err  any        `json:"err,omitempty"`
}

func Reply(c *gin.Context, data any) {
	httpCode := http.StatusOK
	resp := Response[any]{Code: codes.OK, Data: struct{}{}}
	switch v := data.(type) {
	case nil:
	case *status.Status:
		resp.Code = v.GetCode()
		resp.Msg = v.GetMessage()
	case pkgvalidator.ValidationErrors:
		resp.Code = errs.BadRequest.Code
		resp.Msg = errs.BadRequest.Message
		resp.Err = validator.TransErrCtx(c, v).Error()
	case error:
		resp.Code = errs.InternalError.Code
		resp.Msg = errs.InternalError.Message
	default:
		resp.Data = data
	}
	c.JSON(httpCode, resp)
}
