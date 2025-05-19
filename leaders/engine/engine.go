package engine

import (
	"log"

	"exmpl.com/leaders/config"
	"exmpl.com/leaders/sqlite"
	"github.com/expr-lang/expr"
)

var Competitions = make([]sqlite.Competition, 0)

// read from channel into a slice
// program, err := expr.Compile(code, expr.Env(Env{}))

func InitEngine() {
	// TODO: make loading only for the comps that are relevant for current time
	comps, err := sqlite.GetAllCompetitions()
	if err != nil {
		log.Fatal("Error loading competitions:", err)
	}
	log.Printf("Loaded %d competitions\n", len(comps))
	for _, comp := range comps {
		program, err := expr.Compile(comp.Rules, expr.Env(sqlite.Event{}))
		if err != nil {
			log.Printf("Error compiling rules: %s for competition %d. Marking as invalid", err, comp.Id)
			comp.Compiles = false
		}
		comp.Compiles = true
		comp.CompiledRules = program
		Competitions = append(Competitions, comp)

	}
}

func ProcessEvent(event *sqlite.Event) {
	select {
	case newCompId := <-config.AppConfig.CompsChannel:
		newComp, err := sqlite.GetCompetitionById(newCompId)
		if err != nil {
			log.Println("Error getting competition by ID:", err) // shouldn't really happen
			return
		}
		log.Println("New competition received:", newComp)
		program, err := expr.Compile(newComp.Rules, expr.Env(sqlite.Event{}))
		if err != nil {
			log.Printf("Error compiling rules: %s for competition %d. Marking as invalid", err, newComp.Id)
			newComp.Compiles = false
		}
		newComp.Compiles = true
		newComp.CompiledRules = program
		Competitions = append(Competitions, newComp)
		log.Println("New competitions count:", len(Competitions))
		processEvent(event)
	default:
		log.Println("No new competitions, processing event")
		processEvent(event)
	}
}

func processEvent(event *sqlite.Event) {
	sqlite.CreateUser(event.UserId)
	sqlite.CreateEvent(event)
	for _, comp := range Competitions {
		if comp.IsRunningNow() && comp.Compiles {
			output, err := expr.Run(comp.CompiledRules, event)
			if err != nil {
				panic(err)
			}
			if output != 0 {
				sqlite.CreateBet(event, comp.Id)
				log.Println("Event processed successfully for comp: ", comp.Id)
			} else {
				log.Printf("No processing for this comp.")
				// log.Printf("No processing for this comp. Rules %s, evt: %+v \n", comp.Rules, event)
			}
		}
	}
}
