package user

import (
	"github.com/luoye-g/webgate/model"
	"gorm.io/gorm"
)

type User interface {
	Create(userName, pass, phone, email string) (*model.User, error)
	GetUserByUserNameAndPass(username, password string) (*model.User, error)
	GetUserByPhoneOrEmailAndPass(phone, email, password string) (*model.User, error)
	GetByUserID(id uint64) (*model.User, error)
	GetByPhone(phone string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByPhoneOrEmail(phone, email string) (*model.User, error)
	UpdateUserName(id uint64, userName string) error
}

// GetUserRepo returns the singleton instance of the user repository.
// If the repository hasn't been initialized, it creates a new instance with a UserDAO.
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

func (u *user) GetByPhone(phone string) (*model.User, error) {
	user, err := u.UserDao.FindByPhone(phone)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func (u *user) GetByEmail(email string) (*model.User, error) {
	user, err := u.UserDao.FindByEmail(email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func (u *user) GetUserByPhoneOrEmailAndPass(phone, email, password string) (*model.User, error) {
	user, err := u.UserDao.GetUserByPhoneOrEmailAndPass(phone, email, password)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func (u *user) GetByPhoneOrEmail(phone, email string) (*model.User, error) {
	user, err := u.UserDao.GetByPhoneOrEmail(phone, email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func (u *user) UpdateUserName(id uint64, userName string) error {
	return u.UserDao.UpdateUserName(id, userName)
}
