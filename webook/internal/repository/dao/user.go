package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("邮箱已被注册，请换个邮箱重新申请")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreatedAt = now
	u.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var me *mysql.MySQLError
	//if me, ok := err.(*mysql.MySQLError); ok {
	if errors.As(err, &me) {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 返回预定义错误
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u = User{Id: id}
	err := dao.db.WithContext(ctx).First(&u).Error
	return u, err
}

func (dao *UserDao) Update(ctx *gin.Context, user User) error {
	// save会更新所有字段，即使字段是零值
	//err := dao.db.Save(&user).Error
	err := dao.db.WithContext(ctx).Model(&user).Updates(User{
		NickName:        user.NickName,
		BirthDate:       user.BirthDate,
		PersonalProfile: user.PersonalProfile,
		//UpdatedAt:       time.Now().UnixMilli(),
	}).Error
	return err
}

type User struct {
	Id              int64  `gorm:"primaryKey,autoincrement"`
	Email           string `gorm:"unique"`
	Password        string
	NickName        string
	BirthDate       string
	PersonalProfile string
	CreatedAt       int64
	UpdatedAt       int64
}
