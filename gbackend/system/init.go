package system

import (
	"encoding/json"
	"io"
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
	config = &InitConfig{}
	file, err := os.Open("env.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	contents, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(contents, config); err != nil {
		panic(err)
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
