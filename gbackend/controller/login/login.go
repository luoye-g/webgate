package login

import (
	"github.com/gin-gonic/gin"
	"github.com/luoye-g/webgate/repository/user"
)

type LoginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context) {
	var param LoginParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
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

	userRepo := user.GetUserRepo()
	user, err := userRepo.GetUserByUserNameAndPass(param.Username, param.Password)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "服务器错误",
		})
		return
	}

	if user == nil {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

}
