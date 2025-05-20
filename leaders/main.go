package main

import (
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

	// TODO: think how to implement re-init of the engine
	// when a new competition is added but no events are received
	config.AppConfig.CompsChannel = make(chan int64, 100)
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

	go consumer.ConsumeEvents(ctx, &config.AppConfig)

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/competitions", handlers.CompetitionsHandler)
	http.HandleFunc("/leaderboard", handlers.GetLeaderboard)
	server := http.Server{
		Addr: ":8080",
	}

	go func(ch chan bool, s *http.Server) {
		<-ch
		log.Println("Shutting down the server in 2s")
		config.AppConfig.Shutdown = true
		time.Sleep(2 * time.Second)
		s.Shutdown(ctx)
	}(shutdown_ch, &server)

	server.ListenAndServe()

	defer config.AppConfig.Db.Close()
}
