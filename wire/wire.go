//go:build wireinject

package wire

import (
	"geek-basic-go/wire/repository"
	"geek-basic-go/wire/repository/dao"
	"github.com/google/wire"
)

func InitUserRepository() *repository.UserRepository {
	// 传入的顺序没有关系
	// 如果传入了没有用上会报错，如传入InitRedis
	// 如果漏传也会报错，如没有传入InitDB
	// wire本身不支持单例模式，多次初始化会生成多个实例
	// 依赖注入是控制反转的一种现实形式，还有一种是依赖发现
	wire.Build(repository.NewUserRepository, dao.NewUserDao, InitDB)
	return &repository.UserRepository{}
}
