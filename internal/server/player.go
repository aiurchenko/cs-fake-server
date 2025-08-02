package server

import "time"

type Player struct {
	Name        string
	Score       int32
	ConnectedAt time.Time
	Skill       float32
}

func (s *FakeServer) AddPlayer(name string, score int32, secondsOnServer float64, skill float32) {
	if skill < 0.0 {
		skill = 0.0
	} else if skill > 1.0 {
		skill = 1.0
	}

	player := Player{
		Name:        name,
		Score:       score,
		ConnectedAt: time.Now().Add(-time.Duration(secondsOnServer * float64(time.Second))),
		Skill:       skill,
	}
	s.mu.Lock()
	s.players = append(s.players, player)
	s.mu.Unlock()
}
