package user

import (
	"context"
	"fmt"
	"time"

	"github.com/luoye-g/webgate/errors"
	"github.com/luoye-g/webgate/model"
	"github.com/luoye-g/webgate/pkg/endecrypt"
	"github.com/luoye-g/webgate/pkg/redis"
)

func (loginService *UserService) Login(identifier, identifierType, password string) (userSession string, timeout int64, err *errors.AppError) {
	var user *model.User
	var e error
	if identifierType == model.IdentifierTypeEmail && validateEmail(identifier) {
		user, e = loginService.userRepo.GetByEmail(identifier)
	} else if identifierType == model.IdentifierTypePhone && validatePhone(identifier) {
		user, e = loginService.userRepo.GetByPhone(identifier)
	} else {
		return "", 0, errors.NewError(errors.ValidationError, fmt.Errorf("手机号或邮箱格式错误"))
	}
	if e != nil {
		return "", 0, errors.NewError(errors.InternalError, e)
	}
	if user == nil {
		return "", 0, errors.NewError(errors.NotFoundError, fmt.Errorf("用户不存在"))
	}
	if user.PassWord != password {
		return "", 0, errors.NewError(errors.ValidationError, fmt.Errorf("密码错误"))
	}

	timeout = 360 * 24 // 360 * 24 = 1 day
	userSession, e = endecrypt.EncryptAES([]byte(fmt.Sprintf("%v", user.ID)), endecrypt.UserSessionkey)
	if e != nil {
		return "", 0, errors.NewError(errors.InternalError, e)
	}
	if userSession == "" {
		return "", 0, errors.NewError(errors.InternalError, fmt.Errorf("用户session生成失败"))
	}

	// 设置缓存
	if e = redis.GetRedisCli().Set(context.Background(), userSession, user.ID,
		time.Duration(timeout)*time.Second).Err(); e != nil {
		return "", 0, errors.NewError(errors.InternalError, e)
	}

	return userSession, time.Now().Add(time.Duration(timeout * int64(time.Second))).Unix(), nil
}

func (loginService *UserService) Logout(userSession string) (err *errors.AppError) {
	e := redis.GetRedisCli().Del(context.Background(), userSession).Err()
	return errors.NewError(errors.InternalError, e)
}
