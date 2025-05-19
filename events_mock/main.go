package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

const BurstFactor = 5
const RatePerSec = 50
const UserCount = 20
const BurstProbability = 0.2
const BetsChannel = "bets"

type Cur struct {
	name    string
	ex_rate float64
}

var Currencies = [12]Cur{
	{"KWD", 3.2597402597402594},
	{"BHD", 2.662337662337662},
	{"OMR", 2.61038961038961},
	{"JOD", 1.4155844155844157},
	{"GBP", 1.2987012987012987},
	{"KYD", 1.2077922077922079},
	{"GIP", 1.2987012987012987},
	{"CHF", 1.12987012987013},
	{"EUR", 1.0909090909090908},
	{"USD", 1.0},
	{"BTC", 0.0000097},
	{"ETH", 0.00039},
}

var Games = [5]string{"Poker", "Blackjack", "Roulette", "Baccarat", "Slots"}
var Distributors = [6]string{"DistributorX", "DistributorY", "DistributorZ", "DistributorA", "DistributorB", "DistributorC"}
var Studios = [5]string{"StudioX", "StudioY", "StudioZ", "StudioA", "StudioB"}

var ctx = context.Background()

func getRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}

func main() {
	limiter := rate.NewLimiter(RatePerSec, RatePerSec*BurstFactor)
	// control := make(chan bool)
	fmt.Printf("Starting with token balance: %f \n", limiter.Tokens())

	var redis_client = getRedisClient()

	// go func() {
	for {
		bets_count := RatePerSec
		if rand.Float32() > BurstProbability {
			bets_count = RatePerSec * (rand.Intn(BurstFactor) + 1)
		}

		r := limiter.ReserveN(time.Now(), bets_count)
		if !r.OK() {
			fmt.Println("cannoooooot")
		}
		fmt.Printf("Tokens balance: %f, sleeping for %d seconds(%d ns) for bets count: %d \n", limiter.Tokens(), r.Delay()/1000000000, r.Delay(), bets_count)

		time.Sleep(r.Delay())
		gen_events(bets_count, redis_client)
		fmt.Println()
	}
	// }()

	// Wait for a while to let the goroutine finish
	// time.Sleep(10 * time.Second)
	// fmt.Println("All requests processed")
}

func gen_events(num int, redis_client *redis.Client) {
	fmt.Printf("generated %d events for %d users. Time: %s \n", num, rand.Intn(UserCount)+1, time.Now().Format(time.TimeOnly))
	fmt.Println(redis_client)

	for i := 0; i <= num; i++ {
		user_id := rand.Intn(UserCount) + 1
		var bet = build_bet(user_id)
		var bet_json, _ = json.Marshal(bet)
		err := redis_client.Publish(ctx, BetsChannel, bet_json).Err()
		if err != nil {
			panic(err)
		}
	}

	// for i := 0; i < num; i++ {
	// 	// Simulate event generation
	// 	fmt.Println("Event generated for user:", rand.Intn(UserCount))
	// }
}

func build_bet(user_id int) Event {
	var bet = rand.Intn(100)
	var event_type = randEventType()
	var event_currency = randEventCurrency()
	var event = Event{
		EventType:    event_type,
		UserId:       user_id,
		Amount:       bet,
		Currency:     event_currency.name,
		ExchangeRate: event_currency.ex_rate,
		Game:         randEventGame(),
		Distributor:  randEventDistributor(),
		Studio:       randEventStudio(),
		Timestamp:    time.Now().Format(time.RFC3339),
	}
	log.Printf("User: %d | bet: %d", user_id, bet)
	return event

}

type Event struct {
	EventType    string  `json:"event_type"` // bet, win, loss
	UserId       int     `json:"user_id"`
	Amount       int     `json:"amount"`
	Currency     string  `json:"currency"`
	ExchangeRate float64 `json:"exchange_rate"`
	Game         string  `json:"game"`
	Distributor  string  `json:"distributor"`
	Studio       string  `json:"studio"`
	Timestamp    string  `json:"timestamp"` // make this a Time type
}

func randEventType() string {
	event_types := [3]string{"bet", "win", "loss"}
	return event_types[rand.Intn(3)]
}

func randEventCurrency() Cur {
	return Currencies[rand.Intn(len(Currencies))]
}

func randEventGame() string {
	return Games[rand.Intn(len(Games))]
}
func randEventDistributor() string {
	return Distributors[rand.Intn(len(Distributors))]
}
func randEventStudio() string {
	return Studios[rand.Intn(len(Studios))]
}

// {
// 	"event_type": "bet",
// 	"user_id": "123",
// 	"amount": 0.03,
// 	"currency": "BTC",
// 	"exchange_rate": 0.00001058,
// 	"game": "Poker",
// 	"distributor": "DistributorX",
// 	"studio": "StudioY",
// 	"timestamp": "2025-02-04T12:00:00Z"
//   }
