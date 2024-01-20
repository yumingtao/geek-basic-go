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
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var interactiveSvcSet = wire.NewSet(
	dao.NewGormInteractiveDao,
	cache.NewInteractiveRedisCache,
	repository.NewCachedInteractiveRepository,
	service.NewInteractiveServiceImpl,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitDB, ioc.InitRedis, ioc.InitLogger,
		// Dao
		dao.NewUserDao,
		dao.NewGormDBArticleDao,

		interactiveSvcSet,

		// Cache
		cache.NewUserCache /*cache.NewRedisCodeCache,*/, cache.NewGoCacheCodeCache, cache.NewArticleRedisCache,
		// repository
		repository.NewCachedUserRepository, repository.NewCachedCodeRepository, repository.NewArticleRepository,
		// service
		ioc.InitSmsService, service.NewUserService, service.NewCodeService, service.NewArticleService,
		ioc.InitWechatService,
		// handler
		web.NewUserHandler,
		ijwt.NewRedisJwtHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
