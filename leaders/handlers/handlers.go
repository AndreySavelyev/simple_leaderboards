package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	st, _ := strconv.Atoi(r.PostFormValue("start_at"))
	en, _ := strconv.Atoi(r.PostFormValue("end_at"))
	rules := r.PostFormValue("rules")

	sqlite.InsertCompetition(st, en, rules)

	fmt.Fprintf(w, "Competition created\n")
}

func CompetitionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetCompetitions(w, r)
	} else {
		PostCompetition(w, r)
	}
}
