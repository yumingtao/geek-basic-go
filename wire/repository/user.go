package repository

import "geek-basic-go/wire/repository/dao"

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(d *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: d,
	}
}
