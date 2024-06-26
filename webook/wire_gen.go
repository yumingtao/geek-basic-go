// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"geek-basic-go/webook/internal/events/article"
	"geek-basic-go/webook/internal/repository"
	"geek-basic-go/webook/internal/repository/cache"
	"geek-basic-go/webook/internal/repository/dao"
	"geek-basic-go/webook/internal/service"
	"geek-basic-go/webook/internal/web"
	"geek-basic-go/webook/internal/web/jwt"
	"geek-basic-go/webook/ioc"
	"github.com/google/wire"
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

func InitWebServer() *App {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJwtHandler(cmdable)
	loggerV1 := ioc.InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, handler, loggerV1)
	db := ioc.InitDB(loggerV1)
	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDao, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewGoCacheCodeCache()
	codeRepository := repository.NewCachedCodeRepository(codeCache)
	smsService := ioc.InitSmsService()
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, handler, loggerV1)
	wechatService := ioc.InitWechatService(loggerV1)
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, handler)
	articleDao := dao.NewGormDBArticleDao(db)
	articleCache := cache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewArticleRepository(articleDao, userRepository, articleCache)
	client := ioc.InitSaramaClient()
	syncProducer := ioc.InitSyncProducer(client)
	producer := article.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, producer)
	interactiveDao := dao.NewGormInteractiveDao(db)
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDao, loggerV1, interactiveCache)
	interactiveService := service.NewInteractiveServiceImpl(interactiveRepository)
	articleHandler := web.NewArticleHandler(articleService, interactiveService, loggerV1)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler, articleHandler)
	interactiveReadEventConsumer := article.NewInteractiveReadEventConsumer(interactiveRepository, client, loggerV1)
	v2 := ioc.InitConsumers(interactiveReadEventConsumer)
	app := &App{
		server:    engine,
		consumers: v2,
	}
	return app
}

// wire.go:

var interactiveSvcSet = wire.NewSet(dao.NewGormInteractiveDao, cache.NewInteractiveRedisCache, repository.NewCachedInteractiveRepository, service.NewInteractiveServiceImpl)
