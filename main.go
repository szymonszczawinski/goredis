package main

import (
	"context"
	"fmt"
	"goredis/client"
	"log"
	"log/slog"
	"time"
)

func main() {
	cfg := Config{
		ListenAddress: ":3000",
	}
	server := NewServer(cfg)
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)
	client := client.NewClinet("localhost:3000")
	for i := 0; i < 5; i++ {
		client.Set(context.TODO(), fmt.Sprintf("hello %v", i), "world")
		// time.Sleep(time.Millisecond)
	}

	val, err := client.Get(context.TODO(), "hello 0")
	if err != nil {
		slog.Error("error client read value", "err", err)
	}
	slog.Info("clinet read value", "val", val)

	time.Sleep(5 * time.Second)
	fmt.Println(server.storage.data)
}
