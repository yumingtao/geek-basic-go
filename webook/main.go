package main

import (
	"geek-basic-go/webook/internal/repository"
	"geek-basic-go/webook/internal/repository/dao"
	"geek-basic-go/webook/internal/service"
	"geek-basic-go/webook/internal/web"
	"geek-basic-go/webook/internal/web/middlewares/login"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	initUserHdl(db, server)

	err := server.Run(":8080")
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
		AllowHeaders:  []string{"Content-Type", "Authorization"},
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

	loginMiddleware := &login.MiddlewareBuilder{}
	// store用于存储数据
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("ssid", store), loginMiddleware.CheckLogin())
	return server
}

func initDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
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
	server.POST("/users/login", hdl.Login)
	server.GET("/users/:id", hdl.Profile)
	server.PUT("/users/:id", hdl.Edit)
}
