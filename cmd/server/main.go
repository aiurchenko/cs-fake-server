package main

import (
	"math/rand"
	"time"

	"github.com/aiurchenko/cs-fake-server/internal/config"
	"github.com/aiurchenko/cs-fake-server/internal/server"
)

func main() {
	cfg := config.LoadFromEnv()
	s := server.New(cfg)
	rand.Seed(time.Now().UnixNano())
	if err := s.Start(); err != nil {
		panic(err)
	}
}
