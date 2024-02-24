package ioc

import (
	"geek-basic-go/webook/internal/repository/dao"
	"geek-basic-go/webook/pkg/gormx"
	"geek-basic-go/webook/pkg/logger"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
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
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询
			SlowThreshold: 0,
			LogLevel:      glogger.Info,
		}),
	})

	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{
					"thread_running",
				},
			},
		},
	}))
	if err != nil {
		panic(err)
	}
	cb := gormx.NewCallbacks(prometheus2.SummaryOpts{
		Namespace: "geektime_yumingtao",
		Subsystem: "webook",
		Name:      "gorm_db",
		Help:      "这是一个统计GORM的数据库查询",
		ConstLabels: map[string]string{
			"instance_id": "my_instance",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	err = db.Use(cb)
	if err != nil {
		return nil
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

// 这个叫做函数衍生类型实现接口
type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(s string, i ...interface{}) {
	g(s, logger.Field{
		Key: "args",
		Val: i,
	})
}
