package user

import (
	"github.com/gin-gonic/gin"
	"github.com/luoye-g/webgate/controller/response"
	"github.com/luoye-g/webgate/services/user"
)

type RegisterParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func Register(ctx *gin.Context) {
	var param RegisterParam
	if err := ctx.BindJSON(&param); err != nil {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if param.Username == "" || param.Password == "" {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	user, err := user.GetUserService().Register(param.Username, param.Password, param.Phone, param.Email)
	if err != nil {
		response.Failed(ctx, err)
		return
	}
	response.Success(ctx, user)
}
