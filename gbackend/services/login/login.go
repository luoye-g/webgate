package login

import (
	"context"
	"fmt"
	"time"

	"github.com/luoye-g/webgate/pkg/endecrypt"
	"github.com/luoye-g/webgate/pkg/redis"
	"github.com/luoye-g/webgate/repository/user"
)

type LoginService struct {
	userRepo user.User
}

var loginService *LoginService

func GetLoginService() *LoginService {
	if loginService == nil {
		loginService = &LoginService{userRepo: user.GetUserRepo()}
	}
	return loginService
}

func (loginService *LoginService) Login(username, password string) (userSession string, timeout int, err error) {
	user, err := loginService.userRepo.GetUserByUserNameAndPass(username, password)
	if err != nil {
		return "", 0, err
	}
	if user == nil {
		return "", 0, err
	}

	timeout = 360 * 24 // 360 * 24 = 1 day
	userSession, err = endecrypt.EncryptAES([]byte(fmt.Sprintf("%v", user.ID)), endecrypt.UserSessionkey)
	if err != nil {
		return "", 0, err
	}
	if userSession == "" {
		return "", 0, err
	}

	// 设置缓存
	if err = redis.GetRedisCli().Set(context.Background(), userSession, user.ID,
		time.Duration(timeout)*time.Second).Err(); err != nil {
		return "", 0, err
	}
	return userSession, timeout, nil
}

func (loginService *LoginService) Logout(userSession string) (err error) {
	return redis.GetRedisCli().Del(context.Background(), userSession).Err()
}
