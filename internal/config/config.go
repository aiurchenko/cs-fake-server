package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Name       string
	Map        string
	MaxPlayers byte
	Port       int
}

func LoadFromEnv() *ServerConfig {
	_ = godotenv.Load()

	name := os.Getenv("SERVER_NAME")
	if name == "" {
		log.Fatal("SERVER_NAME is required")
	}

	m := os.Getenv("SERVER_MAP")
	if m == "" {
		log.Fatal("SERVER_MAP is required")
	}

	maxPlayers, err := strconv.Atoi(os.Getenv("SERVER_MAX_PLAYERS"))
	if err != nil {
		maxPlayers = 32
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		port = 27015
	}

	return &ServerConfig{
		Name:       name,
		Map:        m,
		MaxPlayers: byte(maxPlayers),
		Port:       port,
	}
}

