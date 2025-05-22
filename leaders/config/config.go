package config

import (
	"database/sql"

	"exmpl.com/leaders/repository"
	"github.com/redis/go-redis/v9"
)

type Cfg struct {
	Db                 *sql.DB
	RedisClient        *redis.Client
	BetsChannel        string
	Shutdown           bool
	CompsChannel       chan int64
	PersistenceService *repository.PersistenceService
}

var AppConfig = Cfg{Shutdown: false}
