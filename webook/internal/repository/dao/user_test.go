package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name      string
		mock      func(t *testing.T) *sql.DB
		ctx       context.Context
		u         User
		wantedErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 注意这里不能加defer，否则*errors.errorString(&errors.errorString{s:"sql: database is closed"})
				/*defer func(db *sql.DB) {
					_ = db.Close()
				}(db)*/
				// 注意这里不能加ExpectBegin和ExpectClose，否则
				//*errors.errorString(&errors.errorString{s:"call to ExecQuery 'INSERT INTO `users`
				//(`email`,`password`,`nick_name`,`birth_date`,`personal_profile`,`phone`,`created_at`,`u_at`)
				//VALUES (?,?,?,?,?,?,?,?)' with args [{Name: Ordinal:1 Value:<nil>} {Name: Ordinal:2 Value:}
				//{Name: Ordinal:3 Value:Tom} {Name: Ordinal:4 Value:} {Name: Ordinal:5 Value:}
				//{Name: Ordinal:6 Value:<nil>} {Name: Ordinal:7 Value:1702807978777}
				//{Name: Ordinal:8 Value:1702807978777}], was not expected,
				//next expectation is: ExpectedBegin => expecting database transaction Begin"})
				//mock.ExpectBegin()
				//mock.ExpectClose()
				mockRes := sqlmock.NewResult(123, 1)
				// 这里要求传入的是sql的正则表达式
				mock.ExpectExec("INSERT INTO").WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			u: User{
				NickName: "Tom",
			},
			wantedErr: nil,
		},

		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 这里要求传入的是sql的正则表达式
				mock.ExpectExec("INSERT INTO").WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				return db
			},
			ctx: context.Background(),
			u: User{
				NickName: "Tom",
			},
			wantedErr: ErrDuplicateEmail,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 这里要求传入的是sql的正则表达式
				mock.ExpectExec("INSERT INTO").WillReturnError(errors.New("数据库错误"))
				return db
			},
			ctx: context.Background(),
			u: User{
				NickName: "Tom",
			},
			wantedErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 没有用ctrl，所以不用gomock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			sqlDB := tc.mock(t)
			// 下边这三个设置，保证gorm不会在初始化过程中发起额外的调用
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDB,
				// 设置成true，跳过检查数据库版本;如果设置为false，gorm在初始化的时候，会调用show version
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// 设置为true，不允许ping数据库
				DisableAutomaticPing: true,
				// 设置为false的时候，即使一个单一的增删改语句，gorm也会开启事务
				SkipDefaultTransaction: true,
			})
			// 初始化DB不能出错，所以这里要断言必须为nil
			assert.NoError(t, err)
			// 注意这里是直接调用的new方法，而不是调用的daomocks下的方法
			// dao := daomocks.NewMockUserDao(ctrl)
			dao := NewUserDao(db)
			err = dao.Insert(tc.ctx, tc.u)
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}
