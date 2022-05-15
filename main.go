package main

import (
	"test/redis/server"
	"test/tcp"
)

func main() {
	config := tcp.NewConfig()
	handers := server.NenHandler()
	tcp.ListenAndServe(config, handers)


}
