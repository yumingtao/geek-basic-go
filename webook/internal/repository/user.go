package repository

import (
	"context"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository/cache"
	"geek-basic-go/webook/internal/repository/dao"
	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

// UserRepository
// Repository 负责操作数据，当然包括操作数据库也包括缓存
type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
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
	du, err := repo.cache.Get(ctx, id)
	if err != nil {
		return du, err
	}
	u, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	du = repo.toDomain(u)
	err = repo.cache.Set(ctx, du)
	if err != nil {
		// 网络崩了，redis崩了
		return domain.User{}, err
	}
	// 可以使用goroutine异步些缓存
	/*go func() {
		err = repo.cache.Set(ctx, du)
		if err != nil {
			log.Println(err)
		}
	}()*/
	return du, nil
}

func (repo *UserRepository) FindByIdV1(ctx context.Context, id int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, id)

	switch err {
	case nil:
		return du, err
	case cache.ErrKeyInexist:
		// key不存在去查询数据库
		u, err := repo.dao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		du = repo.toDomain(u)
		err = repo.cache.Set(ctx, du)
		if err != nil {
			// 网络崩了，redis崩了
			return domain.User{}, err
		}
		// 可以使用goroutine异步些缓存
		/*go func() {
			err = repo.cache.Set(ctx, du)
			if err != nil {
				log.Println(err)
			}
		}()*/
		return du, nil
	default:
		// 接近降级的写法
		return domain.User{}, err
	}
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
