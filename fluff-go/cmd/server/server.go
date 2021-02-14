package main

import (
	"github.com/NoSoundLeR/fluff/fluff-go/api"
)

func main() {
	server := api.NewServer("127.0.0.1:8080")
	server.Run()
}
