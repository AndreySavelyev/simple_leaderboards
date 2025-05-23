package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"exmpl.com/leaders/config"
	"exmpl.com/leaders/engine"
	"exmpl.com/leaders/repository"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	var tmplFile = "index.html"
	tmpl, _ := template.ParseFiles(tmplFile)
	tmpl.Execute(w, nil)
}

func GetCompetitions(w http.ResponseWriter, r *http.Request) {
	rows, err := config.AppConfig.Db.Query("SELECT * FROM competitions")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var competitions []repository.Competition

	for rows.Next() {
		var cm repository.Competition
		if err := rows.Scan(&cm.Id, &cm.StartAt, &cm.EndAt, &cm.Rules); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		competitions = append(competitions, cm)
	}

	w.Header().Set("Content-Type", "text/html")
	funcMap := template.FuncMap{
		"running": func(comp repository.Competition) string { return strconv.FormatBool(comp.IsRunningNow()) },
	}
	var tmplFile = "competitions.html"
	tmpl, _ := template.New(tmplFile).Funcs(funcMap).ParseFiles(tmplFile)
	tmpl.Execute(w, competitions)
}

func PostCompetition(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	st, _ := strconv.Atoi(r.PostFormValue("start_at"))
	en, _ := strconv.Atoi(r.PostFormValue("end_at"))
	rules := r.PostFormValue("rules")

	config.AppConfig.PersistenceService.AddCompetition(st, en, rules)

	log.Println("Competition created:", st, en, rules)
	fmt.Fprintf(w, "Competition created\n")
}

func CompetitionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetCompetitions(w, r)
	} else {
		PostCompetition(w, r)
	}
}

func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	var path, err = url.Parse(r.URL.String())
	if err != nil {
		log.Println("Error parsing URL:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
	}

	compId, err := strconv.Atoi(path.Query().Get("competition_id"))
	if err != nil { // bad conversion
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if compId == 0 {
		http.Error(w, "competition_id is required", http.StatusBadRequest)
	}
	if err != nil { // bad conversion
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var limit int
	limit, err = strconv.Atoi(path.Query().Get("limit"))
	if err != nil { // bad conversion
		limit = 10 // fallback to default
	}
	if limit == 0 {
		limit = 10
	}
	lb, err := config.AppConfig.PersistenceService.GetLeaderboardByCompetitionId(compId, limit)
	if err != nil {
		log.Println("Error getting leaderboard:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	funcMap := template.FuncMap{
		"pretty": func(n float64) string { return strconv.FormatFloat(float64(n), 'f', 4, 64) },
		"reward": func(rank int) int { return engine.RewardForPosition(rank) },
	}
	var tmplFile = "leaderboard.html"
	tmpl, _ := template.New(tmplFile).Funcs(funcMap).ParseFiles(tmplFile)
	tmpl.Execute(w, lb)
}

// func GetLeaderboardJson(w http.ResponseWriter, r *http.Request) {
// 	var path, err = url.Parse(r.URL.String())
// 	if err != nil {
// 		log.Println("Error parsing URL:", err)
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		log.Fatal(err)
// 	}
// 	compId, err := strconv.Atoi(path.Query().Get("competition_id"))
// 	if err != nil { // bad conversion
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	if compId == 0 {
// 		http.Error(w, "competition_id is required", http.StatusBadRequest)
// 	}
// 	if err != nil { // bad conversion
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	var limit int
// 	limit, err = strconv.Atoi(path.Query().Get("limit"))
// 	if err != nil { // bad conversion
// 		limit = 10 // fallback to default
// 	}
// 	if limit == 0 {
// 		limit = 10
// 	}
// 	lb, err := sqlite.GetLeaderboardByCompetitionId(compId, limit)
// 	if err != nil {
// 		log.Println("Error getting leaderboard:", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(lb)
// }

// func GetCompetitions2(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Competitions list\n")
// 	rows, err := config.AppConfig.Db.Query("SELECT * FROM competitions")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()
// 	var competitions []repository.Competition
// 	for rows.Next() {
// 		var cm repository.Competition
// 		if err := rows.Scan(&cm.Id, &cm.StartAt, &cm.EndAt, &cm.Rules); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		competitions = append(competitions, cm)
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(competitions)
// }
