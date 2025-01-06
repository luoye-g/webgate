package middlewares

import (
	"strconv"

	"github.com/gin-gonic/gin"
	pctx "github.com/luoye-g/webgate/pkg/ctx"
	"github.com/luoye-g/webgate/pkg/endecrypt"
	"github.com/luoye-g/webgate/pkg/redis"
	"github.com/luoye-g/webgate/repository/user"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userSession, err := ctx.Cookie("user_session")
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
			ctx.Abort()
			return
		}

		if userSession == "" {
			ctx.JSON(200, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
			ctx.Abort()
			return
		}

		_, err = redis.GetRedisCli().Get(ctx, userSession).Result()
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
			ctx.Abort()
			return
		}

		userIDStr, err := endecrypt.DecryptAES(userSession, endecrypt.UserSessionkey)
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
		}
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
		}
		userInfo, err := user.GetUserRepo().GetByUserID(userID)
		if err != nil || userInfo == nil {
			ctx.JSON(200, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
			ctx.Abort()
			return
		}
		pctx.SetUserID(ctx, userInfo.ID)
	}
}
