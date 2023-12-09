package dao

import (
	"context"
	"database/sql"
	"errors"
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
	u.UAt = now
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
	/*var res User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return res, err*/
	var u = User{Id: id}
	err := dao.db.WithContext(ctx).First(&u).Error
	return u, err
}

func (dao *UserDao) Update(ctx context.Context, user User) error {
	// save会更新所有字段，即使字段是零值
	//err := dao.db.Save(&user).Error
	err := dao.db.WithContext(ctx).Model(&user).Updates(User{
		NickName:        user.NickName,
		BirthDate:       user.BirthDate,
		PersonalProfile: user.PersonalProfile,
		//UAt:       time.Now().UnixMilli(),
	}).Error
	return err
}

func (dao *UserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&res).Error
	return res, err
}

type User struct {
	Id int64 `gorm:"primaryKey,autoincrement"`
	// 代表可以为null
	// Email *string 早起没有sql.NullString，使用*string
	Email           sql.NullString `gorm:"unique"`
	Password        string
	NickName        string
	BirthDate       string
	PersonalProfile string
	Phone           sql.NullString `gorm:"unique"`
	CreatedAt       int64
	UAt             int64
}
