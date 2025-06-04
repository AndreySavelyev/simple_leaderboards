package engine

import (
	"log"
	"reflect"

	// "exmpl.com/leaders/config"
	"exmpl.com/leaders/config"
	"exmpl.com/leaders/repository"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
)

const AmountIdentifier = "amount"

var Competitions = make([]repository.Competition, 0)
var Persistence *repository.PersistenceService

type AmountPatcher struct{}

var floatType = reflect.TypeOf(float64(0))

func (AmountPatcher) Visit(node *ast.Node) {
	if n, ok := (*node).(*ast.IdentifierNode); ok && n.Value == AmountIdentifier {
		cNode := &ast.CallNode{
			Callee:    &ast.IdentifierNode{Value: "BaseAmount"},
			Arguments: []ast.Node{&ast.IdentifierNode{Value: "currency"}, &ast.IdentifierNode{Value: "amount"}},
		}
		ast.Patch(node, cNode)
		(*node).SetType(floatType)
	}
}

func InitEngine(persistence *repository.PersistenceService) {
	Persistence = persistence
	// TODO: make loading only for the comps that are relevant for current time
	comps, err := Persistence.GetCompetitions()
	if err != nil {
		log.Fatal("Error loading competitions:", err)
	}
	log.Printf("Loaded %d competitions\n", len(comps))
	for _, comp := range comps {
		program, err := expr.Compile(comp.Rules, expr.AsFloat64(), expr.Env(repository.Event{}), expr.Patch(AmountPatcher{}))
		if err != nil {
			log.Printf("Error compiling rules: %s for competition %d. Marking as invalid", err, comp.Id)
			comp.Compiles = false
		} else {
			comp.Compiles = true
			comp.CompiledRules = program
		}
		Competitions = append(Competitions, comp)
	}
}

// TODO: add a test
func ProcessEvent(event *repository.Event) {
	select {
	case newCompId := <-config.AppConfig.CompsChannel:
		newComp, err := Persistence.GetCompetitionById(newCompId)
		if err != nil {
			log.Println("Error getting competition by ID:", err) // shouldn't really happen
			return
		}
		log.Println("New competition received:", newComp)
		program, err := expr.Compile(newComp.Rules, expr.AsFloat64(), expr.Env(repository.Event{}), expr.Patch(AmountPatcher{}))
		if err != nil {
			log.Printf("Error compiling rules: %s for competition %d. Marking as invalid", err, newComp.Id)
			newComp.Compiles = false
		} else {
			newComp.Compiles = true
			newComp.CompiledRules = program
		}
		Competitions = append(Competitions, newComp)
		log.Println("New competitions count:", len(Competitions))
		processEvent(event)
	default:
		log.Println("No new competitions, processing event")
		processEvent(event)
	}
}

func processEvent(event *repository.Event) {
	Persistence.CreateUser(event.UserId)
	Persistence.CreateEvent(event)
	for _, comp := range Competitions {
		if comp.IsRunningNow() && comp.Compiles {
			output, err := expr.Run(comp.CompiledRules, event)
			if err != nil {
				log.Println("Error running rules: %s against an event %v. Error: ", comp.Rules, event, err)
				// TODO: skip iteration
			}
			if output != 0 {
				Persistence.CreateBet(event, comp.Id, output.(float64))
				log.Println("Event processed successfully for comp: ", comp.Id)
			} else {
				log.Printf("No processing for this comp. Rules %s, evt: %+v \n", comp.Rules, event)
				log.Println("")
			}
		}
	}
}

func RewardForPosition(rank int) int {
	switch rank {
	case 1:
		return 1000
	case 2:
		return 500
	case 3:
		return 250
	default:
		if rank > 3 && rank <= 50 {
			return 50
		}
		return 0
	}
}
