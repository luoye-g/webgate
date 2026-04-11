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

	IdentifierTypePhone = "phone"
	IdentifierTypeEmail = "email"
)

type User struct {
	ID        uint64    `json:"id" gorm:"id"`
	UserName  string    `json:"user_name" gorm:"user_name"`
	PassWord  string    `json:"pass_word" gorm:"pass_word"`
	Phone     string    `json:"phone" gorm:"phone"`
	Email     string    `json:"email" gorm:"email"`
	Status    uint      `json:"status" gorm:"status"`
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"`
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
		UserName:  userName,
		PassWord:  pass,
		Phone:     phone,
		Email:     email,
		Status:    StatusValid,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := mysql.GetDB().Create(user).Error
	return user, err
}

func (dao *UserDao) GetByUserID(id uint64) (*User, error) {
	user := &User{}
	err := mysql.GetDB().Table(user.TableName()).Where("id = ? and status = ?", id, StatusValid).First(user).Error
	return user, err
}

func (dao *UserDao) GetUserByUserNameAndPass(userName, passWord string) (*User, error) {
	user := &User{}
	err := mysql.GetDB().Table(user.TableName()).Where("user_name = ? and pass_word = ? and status = ?",
		userName, passWord, StatusValid).First(user).Error
	return user, err
}

func (dao *UserDao) GetByPhoneOrEmail(phone, email string) (*User, error) {
	user := &User{}
	err := mysql.GetDB().Table(user.TableName()).Where("(phone = ? or email = ?) and status = ?", phone, email, StatusValid).First(user).Error
	return user, err
}

func (dao *UserDao) FindByEmail(email string) (*User, error) {
	user := &User{}
	err := mysql.GetDB().Table(user.TableName()).Where("email = ? and status = ?", email, StatusValid).First(user).Error
	return user, err
}

func (dao *UserDao) FindByPhone(phone string) (*User, error) {
	user := &User{}
	err := mysql.GetDB().Table(user.TableName()).Where("phone = ? and status = ?", phone, StatusValid).First(user).Error
	return user, err
}

func (dao *UserDao) GetUserByPhoneOrEmailAndPass(phone, email, passWord string) (*User, error) {
	user := &User{}
	err := mysql.GetDB().Table(user.TableName()).Where("(phone = ? or email = ?) and pass_word = ? and status = ?",
		phone, email, passWord, StatusValid).First(user).Error
	return user, err
}

func (dao *UserDao) UpdateUserName(id uint64, userName string) error {
	return mysql.GetDB().Table(UserTableName).Where("id = ? and status = ?", id, StatusValid).
		Update("user_name", userName).Error
}
