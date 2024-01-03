package ioc

import (
	"geek-basic-go/webook/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitLogger() logger.LoggerV1 {
	cfg := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	//l, err := zap.NewDevelopment()
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}

/*func InitErrLogger(l logger.LoggerV1) logger.ErrLogger {
	return logger.NewErrLogger(l)
}*/
