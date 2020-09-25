package main

import (
	"fmt"
	"redis.cli/listen"
)

func main() {
	redisListenUpdate := listen.RedisListenUpdate{}
	fmt.Println("Start Listen Redis ... ")
	redisListenUpdate.StartListen()
}
