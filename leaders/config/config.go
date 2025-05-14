package config

import (
	"database/sql"
)

type Cfg struct {
	Db *sql.DB
	Zz int
}

var AppConfig = Cfg{}