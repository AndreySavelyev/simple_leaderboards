package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"exmpl.com/leaders/config"
	"exmpl.com/leaders/engine"
	"exmpl.com/leaders/repository"
)

func ConsumeEvents(ctx context.Context, config *config.Cfg) {
	var rdb = config.RedisClient
	var bets_channel = config.BetsChannel

	pubsub := rdb.Subscribe(ctx, bets_channel)
	defer pubsub.Close()
	log.Println("Subscribing to channel:", bets_channel)

	_, err := pubsub.Receive(ctx)
	if err != nil {
		fmt.Println("Error subscribing to channel:", err)
		return
	}

	for {
		if config.Shutdown {
			log.Println("Shutting down the consumer")
			break
		}

		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
			// or log & skip?
		}

		event := repository.Event{}

		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Println("Error unmarshalling a message", err)
			continue
		}

		// log.Println(msg.Channel, event)
		// log.Printf("Received ========================: %+v\n", event)
		engine.ProcessEvent(&event)
	}
}
