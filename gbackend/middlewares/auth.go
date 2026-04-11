package middlewares

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luoye-g/webgate/controller/response"
	pctx "github.com/luoye-g/webgate/pkg/ctx"
	"github.com/luoye-g/webgate/pkg/endecrypt"
	"github.com/luoye-g/webgate/pkg/redis"
	"github.com/luoye-g/webgate/repository/user"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userSession, err := ctx.Cookie("user_session")
		if err != nil {
			response.InvalidParam(ctx, "unlogin")
			ctx.Abort()
			return
		}

		if userSession == "" {
			response.InvalidParam(ctx, "unlogin")
			ctx.Abort()
			return
		}

		_, err = redis.GetRedisCli().Get(ctx, userSession).Result()
		if err != nil {
			response.InvalidParam(ctx, "unlogin")
			ctx.Abort()
			return
		}

		userIDStr, err := endecrypt.DecryptAES(userSession, endecrypt.UserSessionkey)
		if err != nil {
			response.InvalidParam(ctx, "unlogin")
		}
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			response.InvalidParam(ctx, "unlogin")
		}
		userInfo, err := user.GetUserRepo().GetByUserID(userID)
		if err != nil || userInfo == nil {
			response.InvalidParam(ctx, "unlogin")
			ctx.Abort()
			return
		}
		pctx.SetUserID(ctx, userInfo.ID)
	}
}
