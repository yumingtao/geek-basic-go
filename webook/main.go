package main

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func main() {
	initViperV1()
	//initViperRemote()
	//initViperWatch()
	initLogger()
	server := InitWebServer()
	server.GET("/hello", func(context *gin.Context) {
		// context核心职责：处理请求，返回响应
		context.String(http.StatusOK, "Hello, World!")
	})
	err := server.Run(":8080")
	if err != nil {
		return
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func initViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	//当前工作目录的config子目录
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	log.Println(viper.Get("test.key"))
}

func initViperV1() {
	// 读取命令行配置的方式 go run . --config=config/dev.yaml
	cFile := pflag.String("config", "config/dev.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cFile)

	// 设置默认值
	//viper.Set("db.dsn", "localhost:3306")
	// viper.SetConfigType("yaml")
	//viper.SetConfigFile("config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	log.Println(viper.Get("test.key"))
}

func initViperV2() {
	cfg := `
test:
  key: value1

redis:
  addr: "localhost:6379"

db:
  dsn: "root:root@tcp(127.0.0.1:13306)/webook?charset=utf8mb4&parseTime=True&loc=Local"
`
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		return
	}
}

func initViperRemote() {
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("远程配置中心发生变更")
	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}

	// 这段代码如果放到err = viper.ReadRemoteConfig()之前，去掉time.sleep，会有并发安全问题
	// 监听远程配置中心事件，推荐使用etcd原始api
	go func() {
		for {
			err = viper.WatchRemoteConfig()
			if err != nil {
				panic(err)
			}
			log.Println("watch:", viper.Get("test.key"))
			time.Sleep(time.Second * 3)
		}
	}()
}

func initViperWatch() {
	cfile := pflag.String("config", "config/config.yaml", "配置文件路径")
	// 这一步之后cfile中才有值
	pflag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println(viper.Get("test.key"))
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}
