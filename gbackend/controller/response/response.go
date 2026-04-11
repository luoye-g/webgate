package response

import (
	"github.com/gin-gonic/gin"
	"github.com/luoye-g/webgate/errors"
)

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

func Failed(ctx *gin.Context, err *errors.AppError) {
	if err != nil {
		if err.Code() == errors.ValidationError {
			ctx.JSON(200, gin.H{
				"code": 400,
				"msg":  err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "internal error",
		})
		return
	}
}

func InvalidParam(ctx *gin.Context, msg string) {
	ctx.JSON(200, gin.H{
		"code": 400,
		"msg":  msg,
	})
}
func RequireLogin(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code": 401,
		"msg":  "require login",
	})
}
