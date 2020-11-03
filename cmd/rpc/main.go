package main

import (
	"github.com/captain-corgi/rabbit-mq-example/cmd/rpc/client"
	"github.com/captain-corgi/rabbit-mq-example/cmd/rpc/server"
)

func main() {
	go func() {
		client.Run()
	}()

	go func() {
		server.Run()
	}()

	select {}
}
