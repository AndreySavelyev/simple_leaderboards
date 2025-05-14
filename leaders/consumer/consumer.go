package consumer

import (
	"context"
	"fmt"
	"log"

	"exmpl.com/leaders/config"
)

func ConsumeEvents2(ctx context.Context, config *config.Cfg) {
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
		}
		fmt.Println(msg.Channel, msg.Payload)
		log.Println(msg.Channel, msg.Payload)
	}
}

// func ConsumeEvents(ctx context.Context, config *config.Cfg, ch chan string) {
// 	var rdb = config.RedisClient
// 	var bets_channel = config.BetsChannel

// 	pubsub := rdb.Subscribe(ctx, bets_channel)
// 	defer pubsub.Close()

// 	_, err := pubsub.Receive(ctx)
// 	if err != nil {
// 		fmt.Println("Error subscribing to channel:", err)
// 		return
// 	}

// 	for {
// 		msg, err := pubsub.ReceiveMessage(ctx)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(msg.Channel, msg.Payload)
// 		if !config.Shutdown {
// 			close(ch)
// 		} else {
// 			ch <- msg.Payload
// 		}
// 	}
// }
