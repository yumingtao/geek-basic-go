package repository

import (
	"context"
	"database/sql"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/repository/cache"
	"geek-basic-go/webook/internal/repository/dao"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateEmail
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	Update(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

// CachedUserRepository
// Repository 负责操作数据，当然包括操作数据库也包括缓存
type CachedUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDao, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	err := repo.dao.Insert(ctx, repo.toEntity(u))

	return err
}

func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:              u.Id,
		Email:           u.Email.String,
		Password:        u.Password,
		NickName:        u.NickName,
		BirthDate:       u.BirthDate,
		PersonalProfile: u.PersonalProfile,
		Phone:           u.Phone.String,
	}
}

func (repo *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
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

func (repo *CachedUserRepository) FindByIdV1(ctx context.Context, id int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, id)

	switch err {
	case nil:
		return du, err
	case cache.ErrKeyNotExist:
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

func (repo *CachedUserRepository) Update(ctx context.Context, u domain.User) error {
	err := repo.dao.Update(ctx, dao.User{
		Id:              u.Id,
		NickName:        u.NickName,
		BirthDate:       u.BirthDate,
		PersonalProfile: u.PersonalProfile,
	})
	return err
}

func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password:        u.Password,
		BirthDate:       u.BirthDate,
		PersonalProfile: u.PersonalProfile,
		NickName:        u.NickName,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
	}
}
