package login

import (
	"crypto/md5"

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

func sessionGen(userName, passWord string) (string, error) {
	hash := md5.New()
	hash.Write([]byte(userName + passWord))
	return string(hash.Sum(nil)), nil
}
