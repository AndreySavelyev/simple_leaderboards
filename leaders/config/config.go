package config

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Cfg struct {
	Db          *sql.DB
	RedisClient *redis.Client
	Zz          int
	BetsChannel string
	Shutdown    bool
}

var AppConfig = Cfg{Shutdown: false}
