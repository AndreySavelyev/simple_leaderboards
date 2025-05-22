package engine

import (
	"database/sql"
	"fmt"
	"testing"

	"exmpl.com/leaders/repository"
	"github.com/stretchr/testify/assert"
)

type MockRepo struct {
}

func (m *MockRepo) GetAllCompetitions(db *sql.DB) ([]repository.Competition, error) {
	return []repository.Competition{
		{
			Id:            1,
			StartAt:       1672531199,
			EndAt:         1672617599,
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

func TestInitengine(t *testing.T) {
	sqliteRepo := MockRepo{}
	db := sql.DB{}
	pers := repository.NewPersistenceservice(&sqliteRepo, &db)

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
