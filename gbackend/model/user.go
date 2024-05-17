package model

import (
	"time"

	"github.com/luoye-g/webgate/pkg/mysql"
)

const (
	StatusInit    = 0
	StatusValid   = 1
	StatusInvalid = 2

	UserTableName = "users"
)

type User struct {
	ID          uint64    `json:"id" gorm:"id"`
	UserName    string    `json:"user_name" gorm:"user_name"`
	PassWord    string    `json:"pass_word" gorm:"pass_word"`
	Phone       string    `json:"phone" gorm:"phone"`
	Email       string    `json:"email" gorm:"email"`
	Status      uint      `json:"status" gorm:"status"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time"`
}

func (u *User) TableName() string {
	return UserTableName
}

type UserDao struct {
}

func NewUserDAO() *UserDao {
	return &UserDao{}
}

func (dao *UserDao) Create(userName, pass, phone, email string) (*User, error) {
	now := time.Now()
	user := &User{
		UserName:    userName,
		PassWord:    pass,
		Phone:       phone,
		Email:       email,
		Status:      StatusValid,
		CreatedTime: now,
		UpdatedTime: now,
	}
	err := mysql.GetDB().Create(user).Error
	return user, err
}

func (dao *UserDao) GetUserByUserNameAndPass(userName, passWord string) (*User, error) {
	user := &User{}
	err := mysql.GetDB().Table(user.TableName()).Where("user_name = ? and pass_word = ? and status = ?",
		userName, passWord, StatusValid).First(user).Error
	return user, err
}
