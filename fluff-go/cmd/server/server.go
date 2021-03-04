package main

import (
	"github.com/NoSoundLeR/fluff/fluff-go/api"
)

func main() {
	server := api.NewServer()
	server.Run()
}
