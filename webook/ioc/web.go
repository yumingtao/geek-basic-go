package ioc

import (
	"geek-basic-go/webook/internal/web"
	"geek-basic-go/webook/internal/web/middlewares/login"
	"geek-basic-go/webook/pkg/ginx/middleware/ratelimit"
	"geek-basic-go/webook/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitWebServerV1(mdls []gin.HandlerFunc, hdls []web.Handler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	// wire做不到通过接口查找实现的能力
	for _, hdl := range hdls {
		hdl.RegisterRoutes(server)
	}
	//userHdl.RegisterRoutes(server)
	return server
}

// InitGinMiddlewares
// 这一部分需要手动添加进来，go和wire目前不能自动发现gin.HandlerFunc实例后自动组装进来
// wire是通过抽象语法树来发现依赖，并注入的
func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			AllowCredentials: true,
			// 不要用*
			AllowOrigins: []string{"http://localhost:3000"},
			//AllowMethods: []string{},
			// 加上Authorization头部
			AllowHeaders: []string{"Content-Type", "Authorization"},
			// 允许前端访问后端响应中带的头部
			ExposeHeaders: []string{"X-Jwt-Token", "X-Refresh-Token"},
			AllowOriginFunc: func(origin string) bool {
				// if strings.Contains(origin, "http://localhost") {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "your_company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		func(ctx *gin.Context) {
			println("这个另一个middleware")
		},
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 100)).Build(),
		(&login.JwtMiddlewareBuilder{}).CheckLogin(),
	}
}
