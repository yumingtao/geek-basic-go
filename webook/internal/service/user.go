package service

import (
	"context"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
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

func (scv *UserService) Edit(ctx *gin.Context, u domain.User) (domain.User, error) {
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

func (scv *UserService) Profile(ctx *gin.Context, id int64) (domain.User, error) {
	u, err := scv.repo.FindById(ctx, id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
