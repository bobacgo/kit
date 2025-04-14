package handler

import (
	"log/slog"

	v1 "github.com/bobacgo/kit/examples/api/admin/v1"
	"github.com/bobacgo/kit/web/r"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

type UserHandler struct {
	// svc biz.SysUser
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// PageList 获取用户分页列表
// @Summary 用户管理
// @Description 获取用户分页列表
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param language header string false "language（可选）"
// @Param req body v1.UserPageListReq true "请求参数"
// @Success 200 {object} []v1.UserPageListResp "查询成功"
// @Router /v1/user/pageList [post]
func (u *UserHandler) PageList(c *gin.Context) {
	lang := c.GetHeader("language")
	c.Set("language", lang)
	req := &v1.UserPageListReq{}
	if err := c.ShouldBind(req); err != nil {
		r.Reply(c, err)
		return
	}
	r.Reply(c, nil)
}

func (u *UserHandler) Get(c *gin.Context) {
	_, span := otel.Tracer("examples-service").Start(c, "GetUserById")
	defer span.End()
	slog.InfoContext(c, "GetUserById", slog.String("id", c.Param("id")))
	r.Reply(c, nil)
}

func (u *UserHandler) Create(c *gin.Context) {
}

func (u *UserHandler) Update(c *gin.Context) {

}

func (u *UserHandler) Delete(c *gin.Context) {

}
