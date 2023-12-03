//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-ymt-mysql:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: RedisConfig{
		Addr: "webook-ymt-redis:6380",
	},
}
