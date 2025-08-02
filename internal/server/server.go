package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/aiurchenko/cs-fake-server/internal/config"
)

var challengeNumber = int32(123456789)

type FakeServer struct {
	cfg     *config.ServerConfig
	players []Player
	mu      sync.RWMutex
}

func New(cfg *config.ServerConfig) *FakeServer {
	server := &FakeServer{
		cfg:     cfg,
		players: []Player{},
	}

	server.AddPlayer("ProPlayer", 5, 1234.5, 0.9) // почти всегда получает
	server.AddPlayer("CasualJoe", 3, 456.7, 0.5)  // иногда
	server.AddPlayer("Bot_Dumb", 1, 789.0, 0.1)   // редко

	return server
}

func (s *FakeServer) Start() error {
	s.StartScoreUpdater()

	addr := net.UDPAddr{
		Port: s.cfg.Port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return fmt.Errorf("ошибка запуска UDP сервера: %w", err)
	}
	defer conn.Close()

	fmt.Printf("🛰️  Сервер запущен на порту %d\n", s.cfg.Port)

	for {
		buf := make([]byte, 1400)
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			continue
		}
		if n < 5 {
			continue
		}
		s.handleRequest(buf[:n], clientAddr, conn)
	}
}

func (s *FakeServer) StartScoreUpdater() {
	go func() {
		for {
			wait := time.Duration(rand.Intn(26)+5) * time.Second
			log.Printf("⏳ Следующее обновление через %v...\n", wait)
			time.Sleep(wait)

			s.mu.Lock()
			if len(s.players) == 0 {
				log.Println("⚠️ Нет игроков для обновления")
				s.mu.Unlock()
				continue
			}

			// 1. Суммируем все веса (Skill)
			var totalWeight float32
			for _, p := range s.players {
				// Можно дать минимальный вес, чтобы 0.0 не были абсолютно исключены
				if p.Skill <= 0.0 {
					totalWeight += 0.01
				} else {
					totalWeight += p.Skill
				}
			}

			// 2. Выбираем случайное значение от 0 до totalWeight
			target := rand.Float32() * totalWeight

			// 3. Пробегаем по игрокам и ищем, кто попал
			var chosen *Player
			acc := float32(0)
			for i := range s.players {
				weight := s.players[i].Skill
				if weight <= 0 {
					weight = 0.01 // минимальный шанс
				}
				acc += weight
				if target <= acc {
					chosen = &s.players[i]
					break
				}
			}

			if chosen != nil {
				chosen.Score += 1
				log.Printf("⭐ Выбран игрок %q (skill %.2f), теперь score: %d", chosen.Name, chosen.Skill, chosen.Score)
			} else {
				log.Println("⚠️ Никто не выбран (что-то пошло не так)")
			}

			s.mu.Unlock()
		}
	}()
}
