
package main

import (
	// "database/sql"
	// "fmt"
	"log"

	// _ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
	"net/http"
	"exmpl.com/leaders/config"
	"exmpl.com/leaders/sqlite"
	"exmpl.com/leaders/handlers"
)

var REDIS_CLIENT *redis.Client

func initDB() {
	REDIS_CLIENT = getRedisClient()
	config.AppConfig.Db = sqlite.InitSqlite()
	// init redis client
	// ...
}


func main() {
	initDB()
	log.Default().Println("Starting server on :8080")
	http.HandleFunc("/", handlers.RootHandler)
	// http.HandleFunc("/leaderboards", get_leaderboards)
	http.HandleFunc("/competitions", handlers.CompetitionsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	defer config.AppConfig.Db.Close()
}


func getRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

type Event struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	UserId    int    `json:"user_id"`
	BetAmount int    `json:"bet_amount"`
}

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

