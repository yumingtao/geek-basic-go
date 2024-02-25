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

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitDB, InitRedis, ioc.InitLogger,
		// Dao
		dao.NewUserDao,
		// Cache
		cache.NewUserCache, cache.NewRedisCodeCache,
		// repository
		repository.NewCachedUserRepository, repository.NewCachedCodeRepository,
		// service
		ioc.InitSmsService, service.NewUserService, service.NewCodeService,
		// handler
		ijwt.NewRedisJwtHandler,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,

		/*logger.NewErrLogger,*/

		ioc.InitWechatService,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
