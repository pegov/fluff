package main

import (
	"os"
	"strings"

	"github.com/NoSoundLeR/fluff/fluff-go/api"
)

func main() {
	bindAddr, ok := os.LookupEnv("BASE_URL")
	if !ok {
		bindAddr = "0.0.0.0:8000"
	}
	redisAddr, ok := os.LookupEnv("REDIS_URL")
	redisAddr = strings.ReplaceAll(redisAddr, "redis://", "")
	redisAddr = strings.ReplaceAll(redisAddr, "/", "")
	if !ok {
		redisAddr = "127.0.0.1:6379"
	}
	server := api.NewServer(bindAddr, redisAddr)
	server.Run()
}
