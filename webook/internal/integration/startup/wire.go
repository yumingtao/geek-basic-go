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
var userSvcProvider = wire.NewSet(
	dao.NewUserDao,
	cache.NewUserCache,
	repository.NewCachedUserRepository,
	service.NewUserService,
)
var articleSvcProvider = wire.NewSet(
	repository.NewArticleRepository,
	cache.NewArticleRedisCache,
	dao.NewGormDBArticleDao,
	service.NewArticleService,
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
		thirdPartySet,
		userSvcProvider,
		articleSvcProvider,
		interactiveSvcSet,
		// Cache
		cache.NewRedisCodeCache,
		// repository
		repository.NewCachedCodeRepository,
		// service
		ioc.InitSmsService, service.NewCodeService,
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

func InitArticleHandler(dao dao.ArticleDao) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		interactiveSvcSet,
		cache.NewArticleRedisCache,
		repository.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}

func InitInteractiveService() service.InteractiveService {
	wire.Build(thirdPartySet, interactiveSvcSet)
	return service.NewInteractiveServiceImpl(nil)
}
