package  sqlite

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"exmpl.com/leaders/config"
)

func InitSqlite() *sql.DB {
	db, err := sql.Open("sqlite3", "./leaderboards.db")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS competitions (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		start_at INTEGER NOT NULL,
		end_at INTEGER NOT NULL,
		rules TEXT
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table 'competitions' created successfully")

	return db
}

func InsertCompetition(start, end int, rules string) {
	_, err := config.AppConfig.Db.Exec(`INSERT INTO competitions (start_at, end_at, rules) VALUES (?, ?, ?)`, start, end, rules)
	if err != nil {
		log.Fatal(err)
		return
	}
}

type Competition struct {
	Id int  `json:"id"`
	StartAt int  `json:"start_at"`
	EndAt int  `json:"end_at"`
	Rules string  `json:"rules"`
}
