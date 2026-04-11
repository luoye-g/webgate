package user

import (
	"fmt"
	"regexp"

	"github.com/luoye-g/webgate/errors"
)

const (
	passwordRegx = `^[A-Za-z0-9_]+$`
	phoneRegx    = `^(\+?86)?1[3-9]\d{9}$`
	emailRegx    = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)

func validatePassword(password string) bool {
	// 密码长度及字符集校验：不超过20个字符，且仅由大小写字母、数字、下划线组成
	if len(password) > 20 {
		return false
	}
	passwordRe := regexp.MustCompile(passwordRegx)
	if !passwordRe.MatchString(password) {
		return false
	}
	return true
}

func validatePhone(phone string) bool {
	// 电话号码校验（中国区手机号，支持可选前缀+86或86）
	phoneRe := regexp.MustCompile(phoneRegx)
	if !phoneRe.MatchString(phone) {
		return false
	}
	return true
}

func validateEmail(email string) bool {
	emailRe := regexp.MustCompile(emailRegx)
	if !emailRe.MatchString(email) {
		return false
	}
	return true
}

type UserInfoDTO struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func (loginService *UserService) Register(userName, pass, phone, email string) (*UserInfoDTO, *errors.AppError) {

	// 用户名长度校验：不超过32个字符
	if len(userName) > 32 {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("用户名不能超过32个字符"))
	}

	if !validatePassword(pass) {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("密码格式不合法"))
	}

	if !validatePhone(phone) {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("电话号码格式不合法"))
	}

	// 邮箱格式校验（简化的常用邮箱正则）
	if !validateEmail(email) {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("邮箱格式不合法"))
	}

	// 检查电话号码或邮箱是否已被注册
	existingUser, err := loginService.userRepo.GetByPhoneOrEmail(phone, email)
	if err != nil {
		return nil, errors.NewError(errors.InternalError, err)
	}
	if existingUser != nil {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("电话号码或邮箱已被注册"))
	}

	user, err := loginService.userRepo.Create(userName, pass, phone, email)
	if err != nil {
		return nil, errors.NewError(errors.InternalError, err)
	}
	return &UserInfoDTO{
		UserID:   user.ID,
		UserName: user.UserName,
		Phone:    user.Phone,
		Email:    user.Email,
	}, nil
}
