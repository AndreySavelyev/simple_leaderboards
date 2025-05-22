package repository

import (
	"database/sql"
	"time"

	"github.com/expr-lang/expr/vm"
)

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

type Repository interface {
	GetAllCompetitions(db *sql.DB) ([]Competition, error)
	GetCompetitionById(db *sql.DB, id int64) (Competition, error)
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
