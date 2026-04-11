package user

import (
	"github.com/gin-gonic/gin"
	"github.com/luoye-g/webgate/controller/response"
	"github.com/luoye-g/webgate/pkg/ctx"
	"github.com/luoye-g/webgate/services/user"
)

type LoginParam struct {
	Identifier string `json:"identifier"`
	LoginType  string `json:"login_type"`
	Password   string `json:"password"`
}

func Login(ctx *gin.Context) {
	var param LoginParam
	if err := ctx.ShouldBindBodyWithJSON(&param); err != nil {
		response.InvalidParam(ctx, "invalid_param")
		return
	}

	if param.Identifier == "" || param.Password == "" || param.LoginType == "" {
		response.InvalidParam(ctx, "invalid_param")
		return
	}

	session, timeout, err := user.GetUserService().Login(param.Identifier, param.LoginType, param.Password)
	if err != nil {
		response.Failed(ctx, err)
		return
	}

	if session == "" {
		response.InvalidParam(ctx, "invalid_user")
		return
	}

	ctx.SetCookie("user_session", session, int(timeout), "/", "", false, false)
	response.Success(ctx, nil)
}

func Logout(c *gin.Context) {
	err := user.GetUserService().Logout(ctx.GetUserSession(c))
	if err != nil {
		response.Failed(c, err)
		return
	}
	response.Success(c, nil)
}
