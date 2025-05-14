package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

const BurstFactor = 5
const RatePerSec = 10
const UserCount = 3
const BurstProbability = 0.2
const BetsChannel = "bets"

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

		err := redis_client.Publish(ctx, BetsChannel, bet).Err()
		if err != nil {
			panic(err)
		}
	}

	// for i := 0; i < num; i++ {
	// 	// Simulate event generation
	// 	fmt.Println("Event generated for user:", rand.Intn(UserCount))
	// }
}

func build_bet(user_id int) string {
	var bet = rand.Intn(100)
	return fmt.Sprintf("User: %d | bet: %d", user_id, bet)
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
