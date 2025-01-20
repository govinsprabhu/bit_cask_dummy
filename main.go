package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	kvstore "github.com/govinsprabhu/kv_store/kv_store"
)

var ctx = context.Background()

func ListenToRedis(redisAddr, channel string) {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	pubsub := rdb.Subscribe(ctx, channel)
	defer pubsub.Close()
	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Println("Received message:", msg.Payload)
	}
}

func main() {
	kvstore.Init_kvstore("kv_store.txt")
	redisAddr := "localhost:6379"
	channel := "kv_store"
	go ListenToRedis(redisAddr, channel)
	http.HandleFunc("/get", kvstore.GetHandler)
	http.HandleFunc("/put", kvstore.PutHandler)
	http.HandleFunc("/delete", kvstore.DeleteHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
