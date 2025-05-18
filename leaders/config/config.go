package config

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Cfg struct {
	Db           *sql.DB
	RedisClient  *redis.Client
	BetsChannel  string
	Shutdown     bool
	CompsChannel chan int64
}

var AppConfig = Cfg{Shutdown: false}
