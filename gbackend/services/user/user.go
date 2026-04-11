package user

import (
	"fmt"

	"github.com/luoye-g/webgate/errors"
	"github.com/luoye-g/webgate/model"
	"github.com/luoye-g/webgate/repository/user"
)

type UserService struct {
	userRepo user.User
}

var userService *UserService

func GetUserService() *UserService {
	if userService == nil {
		userService = &UserService{userRepo: user.GetUserRepo()}
	}
	return userService
}

func (service *UserService) GetUserInfo(userID uint64) (*model.User, *errors.AppError) {
	if userID == 0 {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("用户ID不能为空"))
	}

	userInfo, err := service.userRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.NewError(errors.InternalError, err)
	}

	if userInfo == nil {
		return nil, errors.NewError(errors.NotFoundError, fmt.Errorf("用户不存在"))
	}

	return userInfo, nil
}

func (service *UserService) UpdateUserName(userID uint64, userName string) *errors.AppError {
	if userID == 0 {
		return errors.NewError(errors.ValidationError, fmt.Errorf("用户ID不能为空"))
	}

	if userName == "" {
		return errors.NewError(errors.ValidationError, fmt.Errorf("昵称不能为空"))
	}

	if len([]rune(userName)) > 20 {
		return errors.NewError(errors.ValidationError, fmt.Errorf("昵称长度不能超过20个字符"))
	}

	err := service.userRepo.UpdateUserName(userID, userName)
	if err != nil {
		return errors.NewError(errors.InternalError, err)
	}

	return nil
}
