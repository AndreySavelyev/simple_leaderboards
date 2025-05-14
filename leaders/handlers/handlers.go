package handlers

import (
	"encoding/json"
	"strconv"
	"net/http"
	"fmt"
	"exmpl.com/leaders/config"
	"exmpl.com/leaders/sqlite"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from the leaderboard service\n")
	fmt.Fprintf(w, "To list leaderboards visit /leaderboards\n")
	fmt.Fprintf(w, "To list competitions visit /competitions\n")
}

func GetCompetitions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Competitions list\n")

	rows, err := config.AppConfig.Db.Query("SELECT * FROM competitions")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var competitions []sqlite.Competition

	for rows.Next() {
		var cm sqlite.Competition
		if err := rows.Scan(&cm.Id, &cm.StartAt, &cm.EndAt, &cm.Rules); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		competitions = append(competitions, cm)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(competitions)
}

func PostCompetition(w http.ResponseWriter, r *http.Request) {
	st, _ := strconv.Atoi(r.PostForm.Get("start_at"))
	en, _ := strconv.Atoi(r.PostForm.Get("end_at"))
	var rules = r.PostForm.Get("rules")

	sqlite.InsertCompetition(st, en, rules)

	fmt.Fprintf(w, "Competition created\n")
}

func CompetitionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET"  {
		GetCompetitions(w, r)
	} else {
		PostCompetition(w, r)
	}
}
