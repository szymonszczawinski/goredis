package main

import "log"

func main() {
	cfg := Config{
		ListenAddress: ":3000",
	}
	server := NewServer(cfg)
	log.Fatal(server.Start())
}
