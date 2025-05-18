package main

import (
	// "database/sql"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"exmpl.com/leaders/config"
	"exmpl.com/leaders/consumer"
	"exmpl.com/leaders/engine"
	"exmpl.com/leaders/handlers"
	"exmpl.com/leaders/redis"
	"exmpl.com/leaders/sqlite"
)

func initApp() {
	config.AppConfig.Db = sqlite.InitSqlite()
	config.AppConfig.RedisClient = redis.InitRedis()
	config.AppConfig.BetsChannel = "bets"
	config.AppConfig.CompsChannel = make(chan int64)
	engine.InitEngine()
}

var ctx = context.Background()

func main() {
	initApp()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	shutdown_ch := make(chan bool, 1)

	go func(ch chan bool, config *config.Cfg) {
		sig := <-sigs
		fmt.Println()
		fmt.Println("signal:", sig)
		config.Shutdown = true
		ch <- true
	}(shutdown_ch, &config.AppConfig)

	log.Default().Println("Starting server on :8080")

	go consumer.ConsumeEvents2(ctx, &config.AppConfig)

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/leaderboards", handlers.GetLeaderboards)
	http.HandleFunc("/competitions", handlers.CompetitionsHandler)
	server := http.Server{
		Addr: ":8080",
	}

	go func(ch chan bool, s *http.Server) {
		<-ch
		log.Println("Shutting down the server in 3s")
		config.AppConfig.Shutdown = true
		time.Sleep(3 * time.Second)
		s.Shutdown(ctx)
	}(shutdown_ch, &server)

	server.ListenAndServe()

	// bets_channel := make(chan string)
	// go consumer.ConsumeEvents(ctx, &config.AppConfig, bets_channel)

	defer config.AppConfig.Db.Close()
}

type Event struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	UserId    int    `json:"user_id"`
	BetAmount int    `json:"bet_amount"`
}

// T:
//

// event from redis
// for each competition {
// 	if event.suits?(competition)
// 		// create a bet for this comp
// 	}
// }

// type Event
// type Competition
// type Bet(user, amount, competition_id)

// func eventsListener(redis_client *redis.Client, ch chan string) {
// 	// Listen for events from Redis
// 	pubsub := redisClient.PSubscribe("events")
// 	defer pubsub.Close()

// 	for {
// 		msg, err := pubsub.ReceiveMessage()
// 		if err != nil {
// 			fmt.Println("Error receiving message:", err)
// 			continue
// 		}

// 		fmt.Println("Received message:", msg.Payload)
// 	}

// }
