//go:build wireinject

package startup

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

var thirdPartySet = wire.NewSet(InitDB, InitRedis, InitLogger)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		thirdPartySet,
		// Dao
		dao.NewUserDao,
		dao.NewArticleDao,
		// Cache
		cache.NewUserCache, cache.NewRedisCodeCache,
		// repository
		repository.NewCachedUserRepository, repository.NewCachedCodeRepository, repository.NewArticleRepository,
		// service
		ioc.InitSmsService, service.NewUserService, service.NewCodeService, service.NewArticleService,
		InitWechatService,
		// handler
		web.NewUserHandler,
		ioc.InitGinMiddlewares,
		web.NewArticleHandler,
		ijwt.NewRedisJwtHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitWebServer,
	)
	return gin.Default()
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		dao.NewArticleDao,
		repository.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}
