package main

import (
	"geek-basic-go/webook/config"
	"geek-basic-go/webook/internal/repository"
	"geek-basic-go/webook/internal/repository/dao"
	"geek-basic-go/webook/internal/service"
	"geek-basic-go/webook/internal/web"
	"geek-basic-go/webook/internal/web/middlewares/login"
	"geek-basic-go/webook/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	ginredis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	initUserHdl(db, server)
	//server := gin.Default()
	server.GET("/hello", func(context *gin.Context) {
		// context核心职责：处理请求，返回响应
		context.String(http.StatusOK, "Hello, World!")
	})
	err := server.Run(":8081")
	if err != nil {
		return
	}
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)

	//hdl := &user.UserHandler{}
	hdl := web.NewUserHandler(us)
	// 分散注册
	// 优点：比较有条理 缺点：找路由的时候不好找
	hdl.RegisterRoutes(server)

	// 集中注册
	// 优点：在一个文件中能够看到全部路由 缺点：路由太多找起来费劲
	// registerRoutes(server, hdl)
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// 使用Use方法接入middleware
	server.Use(cors.New(cors.Config{
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
	}), func(ctx *gin.Context) {
		println("这个另一个middleware")
	})

	redisClient := redis.NewClient(&redis.Options{
		//Addr: "localhost:6379",
		Addr: config.Config.Redis.Addr,
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	//useSession(server)
	useJwt(server)
	return server
}

func useJwt(server *gin.Engine) {
	loginMiddleware := &login.JwtMiddlewareBuilder{}
	server.Use(loginMiddleware.CheckLogin())
}

func useSession(server *gin.Engine) {
	// store用于存储数据
	//store := cookie.NewStore([]byte("secret"))
	// 基于内存的实现，第一个参数：身份认证，authentication key，最好是32或64
	// 第二个参数：数据加密 encryption key
	// 数据安全的三个核心概念：身份认证，数据加密，授权（权限控制）
	//store := memstore.NewStore([]byte("ef9a6efa89E711Ee91Bb1A5958B90E3A"), []byte("99c5468490C311Ee91Bb1A5958B90E3A"))
	// 基于redis的实现
	store, err := ginredis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("ef9a6efa89E711Ee91Bb1A5958B90E3A"), []byte("99c5468490C311Ee91Bb1A5958B90E3A"))
	if err != nil {
		panic(err)
	}
	loginMiddleware := &login.MiddlewareBuilder{}
	server.Use(sessions.Sessions("ssid", store), loginMiddleware.CheckLogin())
}

func initDB() *gorm.DB {
	//dsn := "root:root@tcp(127.0.0.1:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn))
	// 采用配置文件中的配置
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func registerRoutes(server *gin.Engine, hdl *web.UserHandler) {
	server.POST("/users", hdl.SignUp)
	//server.POST("/users/login", hdl.Login)
	server.POST("/users/login", hdl.LoginWithJwt)
	server.GET("/users/:id", hdl.Profile)
	server.PUT("/users/:id", hdl.Edit)
}
