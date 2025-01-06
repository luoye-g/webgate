package user

import (
	"github.com/luoye-g/webgate/model"
	"gorm.io/gorm"
)

type User interface {
	Create(userName, pass, phone, email string) (*model.User, error)
	GetUserByUserNameAndPass(username, password string) (*model.User, error)
	GetByUserID(id uint64) (*model.User, error)
	// Update(user *model.User) error
	// Delete(id string) error
}

func GetUserRepo() User {
	if userRepo == nil {
		userRepo = &user{
			UserDao: model.NewUserDAO(),
		}
	}
	return userRepo
}

var userRepo *user

type user struct {
	UserDao *model.UserDao
}

func (u *user) Create(userName, pass, phone, email string) (*model.User, error) {
	return u.UserDao.Create(userName, pass, phone, email)
}

func (u *user) GetUserByUserNameAndPass(username, password string) (*model.User, error) {
	user, err := u.UserDao.GetUserByUserNameAndPass(username, password)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func (u *user) GetByUserID(id uint64) (*model.User, error) {
	user, err := u.UserDao.GetByUserID(id)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}
