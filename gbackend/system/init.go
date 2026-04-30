package system

import (
	"log"
	"os"

	logpkg "github.com/luoye-g/webgate/pkg/log"
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

	LogLevel   string `json:"log_level"`
	LogFormat  string `json:"log_format"`
	LogDir     string `json:"log_dir"`
	LogFile    string `json:"log_file"`
	LogConsole string `json:"log_console"`
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

		LogLevel:   os.Getenv("LOG_LEVEL"),
		LogFormat:  os.Getenv("LOG_FORMAT"),
		LogDir:     os.Getenv("LOG_DIR"),
		LogFile:    os.Getenv("LOG_FILE"),
		LogConsole: os.Getenv("LOG_CONSOLE"),
	}
}

func SystemInit() {
	// load config
	configRead()

	// init logger (do this first so subsequent inits can use it)
	console := config.LogConsole != "false"
	if err := logpkg.Init(logpkg.Config{
		Level:    config.LogLevel,
		Format:   config.LogFormat,
		Dir:      config.LogDir,
		FileName: config.LogFile,
		Console:  console,
	}); err != nil {
		log.Fatalf("init logger failed: %v", err)
	}
	logpkg.Info("logger_initialized",
		"level", config.LogLevel,
		"format", config.LogFormat,
		"dir", config.LogDir,
	)

	// init mysql
	mysql.InitDB(config.MySQLHost, config.MySQLPort, config.MySQLUser, config.MySQLPass)

	// init redis
	redis.InitRedis(config.RedisHost, config.RedisPort, config.RedisPass)
}
