package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

const BurstFactor = 5
const RatePerSec = 10
const BurstProbability = 0.2
const BetsChannel = "bets"

type Cur struct {
	name    string
	exRate float64
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
	{"BTC", 103092.7835051546},
	{"ETH", 2564.1025641026},
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

var UserCount = 10

func main() {
	numUsers := flag.Int("users", 10, "number of users to generate events for")
	flag.Parse()
	UserCount = *numUsers

	limiter := rate.NewLimiter(RatePerSec, RatePerSec*BurstFactor)
	// control := make(chan bool)
	fmt.Printf("Starting with token balance: %f \n", limiter.Tokens())

	var redisClient = getRedisClient()

	for {
		betCount := RatePerSec
		if rand.Float32() > BurstProbability {
			betCount = RatePerSec * (rand.Intn(BurstFactor) + 1)
		}

		r := limiter.ReserveN(time.Now(), betCount)
		if !r.OK() {
			fmt.Println("cannoooooot")
		}
		fmt.Printf("Tokens balance: %f, sleeping for %d seconds(%d ns) for bets count: %d \n", limiter.Tokens(), r.Delay()/1000000000, r.Delay(), betCount)

		time.Sleep(r.Delay())
		genEvents(betCount, redisClient)
		fmt.Println()
	}

}

func genEvents(num int, redisClient *redis.Client) {
	fmt.Printf("generated %d events for %d users. Time: %s \n", num, rand.Intn(UserCount)+1, time.Now().Format(time.TimeOnly))

	for i := 0; i <= num; i++ {
		userId := rand.Intn(UserCount) + 1
		var bet = buildBet(userId)
		var betJson, _ = json.Marshal(bet)
		err := redisClient.Publish(ctx, BetsChannel, betJson).Err()
		if err != nil {
			panic(err)
		}
	}
}

func buildBet(userId int) Event {
	var betAmount = float64(rand.Intn(100))

	var eventType = randEventType()
	var eventCurrency = randEventCurrency()

	// NOTE: BTC & ETH are making too big numbers in the leaderboards
	// so we adjust them a little to make them more realistic
	if eventCurrency.name == "BTC" {
		betAmount = betAmount / 100000.0
	}
	if eventCurrency.name == "ETH" {
		betAmount = betAmount / 2000.0
	}

	var event = Event{
		EventType:    eventType,
		UserId:       userId,
		Amount:       betAmount,
		Currency:     eventCurrency.name,
		ExchangeRate: eventCurrency.exRate,
		Game:         randEventGame(),
		Distributor:  randEventDistributor(),
		Studio:       randEventStudio(),
		Timestamp:    time.Now().Format(time.RFC3339),
	}
	log.Printf("User: %d | betAmount: %f | currency %s | eventType: %s ", userId, betAmount, eventCurrency.name, eventType)
	return event

}

type Event struct {
	EventType    string  `json:"eventType"` // bet, win, loss
	UserId       int     `json:"userId"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	ExchangeRate float64 `json:"exchangeRate"`
	Game         string  `json:"game"`
	Distributor  string  `json:"distributor"`
	Studio       string  `json:"studio"`
	Timestamp    string  `json:"timestamp"` // make this a Time type?
}

func randEventType() string {
	eventTypes := [3]string{"bet", "win", "loss"}
	return eventTypes[rand.Intn(3)]
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
// 	"eventType": "bet",
// 	"userId": "123",
// 	"amount": 0.03,
// 	"currency": "BTC",
//  "exchangeRate": 0.00001058,
// 	"game": "Poker",
// 	"distributor": "DistributorX",
// 	"studio": "StudioY",
// 	"timestamp": "2025-02-04T12:00:00Z"
//   }
