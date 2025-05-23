package engine

import (
	"database/sql"
	"fmt"
	"testing"

	"exmpl.com/leaders/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetAllCompetitions(db *sql.DB) ([]repository.Competition, error) {
	return []repository.Competition{
		{
			Id:            1,
			StartAt:       1,
			EndAt:         1999999999,
			Rules:         "event_type == 'bet' ? amount : 0",
			Compiles:      true,
			CompiledRules: nil,
		},
	}, nil
}

func (m *MockRepo) GetCompetitionById(db *sql.DB, id int64) (repository.Competition, error) {
	return repository.Competition{
		Id:            id,
		StartAt:       1672531199,
		EndAt:         1672617599,
		Rules:         "event_type == 'bet' ? amount : 0",
		Compiles:      true,
		CompiledRules: nil,
	}, nil
}

func (m *MockRepo) CreateCompetition(db *sql.DB, start, end int, rules string) {
}

func (m *MockRepo) CreateUser(db *sql.DB, user_id int) {
}

func (m *MockRepo) CreateBet(db *sql.DB, event *repository.Event, comp_id int64) {
}

func (m *MockRepo) CreateEvent(db *sql.DB, event *repository.Event) {
}

func (m *MockRepo) GetLeaderboardByCompetitionId(db *sql.DB, compId, limit int) (*repository.Leaderboard, error) {
	// TODO: add smth
	return &repository.Leaderboard{}, nil
}

func TestInitengine(t *testing.T) {
	sqliteRepoMock := &MockRepo{}
	db := sql.DB{}

	pers := repository.NewPersistenceservice(sqliteRepoMock, &db)

	InitEngine(pers)
	assert.Equal(t, len(Competitions), 1, "has one loaded competition")
}

func TestRewardForPosition(t *testing.T) {
	var ranks = [][]int{
		{1, 1000},
		{2, 500},
		{3, 250},
		{4, 50},
		{5, 50},
		{50, 50},
		{51, 0},
	}
	for _, rank := range ranks {
		testname := fmt.Sprintf("%d,%d", rank[0], rank[1])
		t.Run(testname, func(t *testing.T) {
			rew := RewardForPosition(rank[0])
			if rew != rank[1] {
				t.Errorf("got %d, expected %d", rew, rank[1])
			}
		})
	}
}

func Test_processEvent(t *testing.T) {
	// var _event = repository.Event{
	// 	EventType:    "bet",
	// 	UserId:       3,
	// 	Amount:       100,
	// 	Currency:     "USD",
	// 	ExchangeRate: 1.0,
	// 	Game:         "game1",
	// 	Distributor:  "distributor1",
	// 	Studio:       "studio1",
	// 	Timestamp:    "2023-10-01T00:00:00Z",
	// }

	// TODO:  mock repo or pers service and check all methods were called
}
