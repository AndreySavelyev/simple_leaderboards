package repository

import (
	"database/sql"
	"time"

	"github.com/expr-lang/expr/vm"
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

type Competition struct {
	Id            int64  `json:"id"`
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
	Timestamp    string  `json:"timestamp" expr:"timestamp"` // make this a Time type?
}

func (e *Event) BaseAmount() float64 {
	return Currencies[e.Currency] * e.Amount
}

type Player struct {
	Id     int     `json:"id"`
	Amount float64 `json:"amount"`
	Rank   int     `json:"rank"`
}

type Leaderboard struct {
	CompetitionId int      `json:"competition_id"`
	Players       []Player `json:"players"`
}

type Repository interface {
	GetAllCompetitions(db *sql.DB) ([]Competition, error)
	GetCompetitionById(db *sql.DB, id int64) (Competition, error)
	CreateCompetition(db *sql.DB, start, end int, rules string)
}

type PersistenceService struct {
	repository Repository
	db         *sql.DB
}

func NewPersistenceservice(repository Repository, db *sql.DB) *PersistenceService {
	return &PersistenceService{
		repository: repository,
		db:         db,
	}
}

func (s *PersistenceService) GetCompetitions() ([]Competition, error) {
	comps, _ := s.repository.GetAllCompetitions(s.db)
	return comps, nil
}

func (s *PersistenceService) GetCompetitionById(id int64) (Competition, error) {
	comp, err := s.repository.GetCompetitionById(s.db, id)
	if err != nil {
		return comp, err
	}
	return comp, nil
}

func (s *PersistenceService) AddCompetition(start, end int, rules string) {
	s.repository.CreateCompetition(s.db, start, end, rules)
}
