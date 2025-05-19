package sqlite

// maybe better to name this to smth like persistence but it's ok for now

import (
	"database/sql"
	"log"
	"time"

	"exmpl.com/leaders/config"
	"github.com/expr-lang/expr/vm"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: should this be in a config?
var Currencies = map[string]float64{
	"KWD": 3.2597402597402594,
	"BHD": 2.662337662337662,
	"OMR": 2.61038961038961,
	"JOD": 1.4155844155844157,
	"GBP": 1.2987012987012987,
	"KYD": 1.2077922077922079,
	"GIP": 1.2987012987012987,
	"CHF": 1.12987012987013,
	"EUR": 1.0909090909090908,
	"USD": 1.0,
	"BTC": 103092.7835051546,
	"ETH": 2564.1025641026,
}

func InitSqlite() *sql.DB {
	db, err := sql.Open("sqlite3", "./leaderboards.db")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	createCompetitionsTable(db)
	createUsersTable(db)
	createBetsTable(db)
	createEventsTable(db)

	return db
}

func createCompetitionsTable(db *sql.DB) {
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS competitions (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		start_at INTEGER NOT NULL,
		end_at INTEGER NOT NULL,
		rules TEXT
    );
  `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table 'competitions' created successfully")
}

func createBetsTable(db *sql.DB) {
	// REAL should really be DECIMAL
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS bets (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			amount REAL NOT NULL,
			competition_id INTEGER NOT NULL
    );
  `
	var _, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	indSqlStmt := `
		CREATE INDEX IF NOT EXISTS bets_idx on bets (competition_id);
	`
	_, err = db.Exec(indSqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table 'bets' created successfully")
}

func createEventsTable(db *sql.DB) {
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			event_type TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			amount REAL NOT NULL,
			currency TEXT NOT NULL,
			exchange_rate REAL NOT NULL,
			game TEXT NOT NULL,
			distributor TEXT NOT NULL,
			studio TEXT NOT NULL,
			timestamp INTEGER NOT NULL
		);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table 'events' created successfully")
}

func createUsersTable(db *sql.DB) {
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS users (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL
		);
	`
	var _, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	indSqlStmt := `
		CREATE UNIQUE INDEX IF NOT EXISTS users_idx on users  (user_id);
	`
	_, err = db.Exec(indSqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table 'users' created successfully")
}

func InsertCompetition(start, end int, rules string) {
	res, err := config.AppConfig.Db.Exec(`INSERT INTO competitions (start_at, end_at, rules) VALUES (?, ?, ?)`, start, end, rules)
	if err != nil {
		log.Fatal(err)
		return
	}
	newId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return
	}
	config.AppConfig.CompsChannel <- newId
}

func GetCompetitionById(Id int64) (Competition, error) {
	var cm Competition
	err := config.AppConfig.Db.QueryRow("SELECT * FROM competitions WHERE id = ?", Id).Scan(&cm.Id, &cm.StartAt, &cm.EndAt, &cm.Rules)
	if err != nil {
		return cm, err
	}
	return cm, nil
}

func GetAllCompetitions() ([]Competition, error) {
	rows, err := config.AppConfig.Db.Query("SELECT * FROM competitions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var competitions []Competition

	for rows.Next() {
		var cm Competition
		if err := rows.Scan(&cm.Id, &cm.StartAt, &cm.EndAt, &cm.Rules); err != nil {
			return nil, err
		}
		competitions = append(competitions, cm)
	}

	return competitions, nil
}

func CreateUser(user_id int) {
	_, err := config.AppConfig.Db.Exec("INSERT OR IGNORE INTO users (user_id) VALUES (?)", user_id)
	if err != nil {
		log.Println("Error creating user:", err)
	}
}

func CreateBet(event *Event, comp_id int) {
	_, err := config.AppConfig.Db.Exec("INSERT INTO bets (user_id, amount, competition_id) VALUES (?, ?, ?)", event.UserId, event.base_amount(), comp_id)
	if err != nil {
		log.Println("Error creating bet:", err)
	}
}

func CreateEvent(event *Event) {
	_, err := config.AppConfig.Db.Exec("INSERT INTO events (event_type, user_id, amount, currency, exchange_rate, game, distributor, studio, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		event.EventType,
		event.UserId,
		event.Amount,
		event.Currency,
		event.ExchangeRate,
		event.Game,
		event.Distributor,
		event.Studio,
		event.Timestamp)
	if err != nil {
		log.Println("Error creating event:", err)
	}
}

type Competition struct {
	Id            int    `json:"id"`
	StartAt       int64  `json:"start_at"`
	EndAt         int64  `json:"end_at"`
	Rules         string `json:"rules"`
	Compiles      bool
	CompiledRules *vm.Program
}

func (c *Competition) IsRunningNow() bool {
	t := time.Now().Unix()
	return c.StartAt <= t && c.EndAt >= t
}

type Event struct {
	Id           int     `json:"id"`
	EventType    string  `json:"event_type" expr:"event_type"` // bet, win, loss
	UserId       int     `json:"user_id" expr:"user_id"`
	Amount       float64 `json:"amount" expr:"amount"`
	Currency     string  `json:"currency"`
	ExchangeRate float64 `json:"exchange_rate"`
	Game         string  `json:"game" expr:"game"`
	Distributor  string  `json:"distributor" expr:"distributor"`
	Studio       string  `json:"studio" expr:"studio"`
	Timestamp    string  `json:"timestamp" expr:"timestamp"` // make this a Time type
}

func (e *Event) base_amount() float64 {
	return Currencies[e.Currency] * e.Amount
}

type Player struct {
	Id     int     `json:"id"`
	Amount float64 `json:"amount"`
}

type Leaderboard struct {
	CompetitionId int      `json:"competition_id"`
	Players       []Player `json:"players"`
}

func GetLeaderboardByCompetitionId(comp_id int, limit int) (Leaderboard, error) {
	var lb Leaderboard
	rows, err := config.AppConfig.Db.Query("SELECT user_id, sum(amount) FROM bets WHERE competition_id = ? group by user_id order by sum(amount) desc limit ?", comp_id, limit)
	if err != nil {
		log.Fatal(err)
		return lb, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.Id, &p.Amount); err != nil {
			log.Fatal(err)
			return lb, err
		}
		lb.Players = append(lb.Players, p)
	}

	lb.CompetitionId = comp_id
	return lb, nil
}
