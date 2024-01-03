//go:build wireinject

package main

import (
	"geek-basic-go/webook/internal/repository"
	"geek-basic-go/webook/internal/repository/cache"
	"geek-basic-go/webook/internal/repository/dao"
	"geek-basic-go/webook/internal/service"
	"geek-basic-go/webook/internal/web"
	ijwt "geek-basic-go/webook/internal/web/jwt"
	"geek-basic-go/webook/ioc"
	"geek-basic-go/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitDB, ioc.InitRedis, ioc.InitLogger, logger.NewErrLogger,
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
		ijwt.NewRedisJwtHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
