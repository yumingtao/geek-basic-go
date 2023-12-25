package ioc

import (
	"geek-basic-go/webook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// 设置默认值
	//viper.Set("db.dsn", "localhost:3306")
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	// 通过接口体设置默认值
	/*var cfg =  Config{
		DSN: "localhost:3306",
	}*/
	var cfg Config
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN))

	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
