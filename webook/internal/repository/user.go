package repository

import (
	"context"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository/dao"
	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	err := repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

	return err
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:              u.Id,
		Email:           u.Email,
		Password:        u.Password,
		NickName:        u.NickName,
		BirthDate:       u.BirthDate,
		PersonalProfile: u.PersonalProfile,
	}
}

func (repo *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) Update(ctx *gin.Context, u domain.User) error {
	err := repo.dao.Update(ctx, dao.User{
		Id:              u.Id,
		NickName:        u.NickName,
		BirthDate:       u.BirthDate,
		PersonalProfile: u.PersonalProfile,
	})
	return err
}
