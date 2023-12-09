package service

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对！")
	ErrUserNotFound          = errors.New("用户不存在! ")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (scv *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return scv.repo.Create(ctx, u)
}

func (scv *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := scv.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (scv *UserService) Edit(ctx context.Context, u domain.User) (domain.User, error) {
	_, err := scv.repo.FindById(ctx, u.Id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	err = scv.repo.Update(ctx, u)
	if err != nil {
		return domain.User{}, err
	}
	nu, err := scv.repo.FindById(ctx, u.Id)
	if err != nil {
		return domain.User{}, err
	}
	return nu, err
}

func (scv *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := scv.repo.FindById(ctx, id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (scv *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 认为大部分用户是已存在用户
	u, err := scv.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// err == nil, 找到用户
		// err != nil, 系统错误
		return u, err
	}
	// 用户没找到，注册用户
	err = scv.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 有两种可能，1. err是phone唯一索引冲突 2.err是系统错误
	if err != nil && !errors.Is(err, repository.ErrDuplicateUser) {
		return domain.User{}, err
	}

	// err == nil 或 ErrDuplicateUser
	// 可能存在主从延迟，理论上应该强制查询主库
	return scv.repo.FindByPhone(ctx, phone)
}
