package sqlite

// maybe better to name this to smth like persistence but it's ok for now

import (
	"database/sql"
	"log"

	"exmpl.com/leaders/config"
	"exmpl.com/leaders/repository"
	_ "github.com/mattn/go-sqlite3"
)

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

type SqliteRepo struct {
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

// func InsertCompetition(start, end int, rules string) {
// }

func (r *SqliteRepo) CreateCompetition(db *sql.DB, start, end int, rules string) {
	res, err := db.Exec(`INSERT INTO competitions (start_at, end_at, rules) VALUES (?, ?, ?)`, start, end, rules)
	if err != nil {
		log.Fatal(err)
		return
	}
	newId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("New competition created with ID:", newId)
	config.AppConfig.CompsChannel <- newId
}

func (r *SqliteRepo) GetCompetitionById(db *sql.DB, id int64) (repository.Competition, error) {
	var cm repository.Competition
	err := config.AppConfig.Db.QueryRow("SELECT * FROM competitions WHERE id = ?", id).Scan(&cm.Id, &cm.StartAt, &cm.EndAt, &cm.Rules)
	if err != nil {
		return cm, err
	}
	return cm, nil
}

func (r *SqliteRepo) GetAllCompetitions(db *sql.DB) ([]repository.Competition, error) {
	rows, err := db.Query("SELECT * FROM competitions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var competitions []repository.Competition

	for rows.Next() {
		var cm repository.Competition
		if err := rows.Scan(&cm.Id, &cm.StartAt, &cm.EndAt, &cm.Rules); err != nil {
			return nil, err
		}
		competitions = append(competitions, cm)
	}

	return competitions, nil
}
func (r *SqliteRepo) CreateUser(db *sql.DB, user_id int) {
	_, err := db.Exec("INSERT OR IGNORE INTO users (user_id) VALUES (?)", user_id)
	if err != nil {
		log.Println("Error creating user:", err)
	}
}

func (r *SqliteRepo) CreateBet(db *sql.DB, event *repository.Event, comp_id int64, contrib float64) {
	_, err := db.Exec("INSERT INTO bets (user_id, amount, competition_id) VALUES (?, ?, ?)", event.UserId, contrib, comp_id)
	if err != nil {
		log.Println("Error creating bet:", err)
	}
}

func (r *SqliteRepo) CreateEvent(db *sql.DB, event *repository.Event) {
	_, err := db.Exec("INSERT INTO events (event_type, user_id, amount, currency, exchange_rate, game, distributor, studio, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
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

func (r *SqliteRepo) GetLeaderboardByCompetitionId(db *sql.DB, comp_id int, limit int) (*repository.Leaderboard, error) {
	var lb repository.Leaderboard
	rows, err := db.Query("SELECT user_id, sum(amount) FROM bets WHERE competition_id = ? group by user_id order by sum(amount) desc limit ?", comp_id, limit)
	if err != nil {
		log.Fatal(err)
		return &lb, err
	}
	defer rows.Close()
	rank := 1
	for rows.Next() {
		var p repository.Player
		if err := rows.Scan(&p.Id, &p.Amount); err != nil {
			log.Fatal(err)
			return &lb, err
		}
		p.Rank = rank
		rank++
		lb.Players = append(lb.Players, p)
	}

	lb.CompetitionId = comp_id
	return &lb, nil
}
