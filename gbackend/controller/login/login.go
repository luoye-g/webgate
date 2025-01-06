package login

import (
	"github.com/gin-gonic/gin"
	"github.com/luoye-g/webgate/services/login"
)

type LoginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context) {
	param := LoginParam{
		Username: ctx.PostForm("username"),
		Password: ctx.PostForm("password"),
	}

	if param.Username == "" || param.Password == "" {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	session, timeout, err := login.GetLoginService().Login(param.Username, param.Password)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "服务器错误",
		})
		return
	}

	if session == "" {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

	ctx.SetCookie("user_session", session, timeout, "/", "", false, true)
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功",
	})
}
