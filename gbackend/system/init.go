package system

import (
	"os"

	"github.com/luoye-g/webgate/pkg/mysql"
	"github.com/luoye-g/webgate/pkg/redis"
)

type InitConfig struct {
	MySQLHost string `json:"mysql_host"`
	MySQLPort string `json:"mysql_port"`
	MySQLUser string `json:"mysql_user"`
	MySQLPass string `json:"mysql_pass"`

	RedisHost string `json:"redis_host"`
	RedisPort string `json:"redis_port"`
	RedisPass string `json:"redis_pass"`
}

var config *InitConfig

func configRead() {
	config = &InitConfig{
		MySQLHost: os.Getenv("MYSQL_HOST"),
		MySQLPort: os.Getenv("MYSQL_PORT"),
		MySQLUser: os.Getenv("MYSQL_USER"),
		MySQLPass: os.Getenv("MYSQL_PASS"),

		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),
		RedisPass: os.Getenv("REDIS_PASS"),
	}
}

func SystemInit() {
	// load config
	configRead()

	// init mysql
	mysql.InitDB(config.MySQLHost, config.MySQLPort, config.MySQLUser, config.MySQLPass)

	// init redis
	redis.InitRedis(config.RedisHost, config.RedisPort, config.RedisPass)
}
