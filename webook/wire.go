//go:build wireinject

package main

import (
	"geek-basic-go/webook/internal/repository"
	"geek-basic-go/webook/internal/repository/cache"
	"geek-basic-go/webook/internal/repository/dao"
	"geek-basic-go/webook/internal/service"
	"geek-basic-go/webook/internal/web"
	"geek-basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitDB, ioc.InitRedis,
		// Dao
		dao.NewUserDao,
		// Cache
		cache.NewUserCache /*cache.NewRedisCodeCache,*/, cache.NewGoCacheCodeCache,
		// repository
		repository.NewCachedUserRepository, repository.NewCachedCodeRepository,
		// service
		ioc.InitSmsService, service.NewUserService, service.NewCodeService,
		ioc.InitWechatService,
		// handler
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
