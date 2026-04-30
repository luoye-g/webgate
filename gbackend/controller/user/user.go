package user

import (
	"github.com/gin-gonic/gin"
	"github.com/luoye-g/webgate/errors"
	pctx "github.com/luoye-g/webgate/pkg/ctx"
	"github.com/luoye-g/webgate/pkg/log"
	"github.com/luoye-g/webgate/services/user"
)

type UserInfoResponse struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func GetUserInfo(ctx *gin.Context) {
	userID := pctx.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(200, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}

	userInfo, err := user.GetUserService().GetUserInfo(userID)
	if err != nil {
		if err.Code() == errors.NotFoundError {
			log.Warn("get_user_info_not_found",
				"trace_id", pctx.GetTraceID(ctx),
				"user_id", userID,
				"err", err.Error(),
			)
			ctx.JSON(200, gin.H{
				"code": 404,
				"msg":  err.Error(),
			})
			return
		}
		log.Error("get_user_info_failed",
			"trace_id", pctx.GetTraceID(ctx),
			"user_id", userID,
			"err", err.Error(),
		)
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "服务器错误",
		})
		return
	}

	response := UserInfoResponse{
		UserName: userInfo.UserName,
		Email:    userInfo.Email,
		Phone:    userInfo.Phone,
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": response,
	})
}

type UpdateUserInfoRequest struct {
	UserName string `json:"user_name" binding:"required"`
}

func UpdateUserInfo(ctx *gin.Context) {
	userID := pctx.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(200, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}

	var req UpdateUserInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn("update_user_info_invalid_param",
			"trace_id", pctx.GetTraceID(ctx),
			"user_id", userID,
			"err", err.Error(),
		)
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数错误：昵称不能为空",
		})
		return
	}

	if appErr := user.GetUserService().UpdateUserName(userID, req.UserName); appErr != nil {
		if appErr.Code() == errors.ValidationError {
			log.Warn("update_user_name_validation_failed",
				"trace_id", pctx.GetTraceID(ctx),
				"user_id", userID,
				"user_name", req.UserName,
				"err", appErr.Error(),
			)
			ctx.JSON(200, gin.H{
				"code": 400,
				"msg":  appErr.Error(),
			})
			return
		}
		log.Error("update_user_name_failed",
			"trace_id", pctx.GetTraceID(ctx),
			"user_id", userID,
			"user_name", req.UserName,
			"err", appErr.Error(),
		)
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "服务器错误",
		})
		return
	}

	log.Info("update_user_name_success",
		"trace_id", pctx.GetTraceID(ctx),
		"user_id", userID,
		"user_name", req.UserName,
	)

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}
