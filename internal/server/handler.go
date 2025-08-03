package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/aiurchenko/cs-fake-server/pkg/utils"
)

func (s *FakeServer) handleRequest(buffer []byte, clientAddr *net.UDPAddr, conn *net.UDPConn) {
	// Логируем весь запрос в hex-формате
	fmt.Printf("📨 Запрос от %s: % X\n", clientAddr, buffer)

	if len(buffer) < 5 {
		fmt.Printf("⚠️ Слишком короткий пакет от %s: длина %d\n", clientAddr, len(buffer))
		return
	}

	header := buffer[4]
	switch header {
	case 0x54: // A2S_INFO
		if strings.HasPrefix(string(buffer[5:]), "Source Engine Query") {
			fmt.Printf("📩 A2S_INFO от %s\n", clientAddr)
			var response bytes.Buffer
			binary.Write(&response, binary.LittleEndian, int32(-1))
			response.WriteByte(0x49)

			response.WriteByte(0x11)
			utils.WriteString(&response, s.cfg.Name)
			utils.WriteString(&response, s.cfg.Map)
			utils.WriteString(&response, "cstrike")
			utils.WriteString(&response, "Counter-Strike 1.6")
			binary.Write(&response, binary.LittleEndian, int16(6011))

			s.mu.RLock()
			players := make([]Player, len(s.players))
			copy(players, s.players)
			s.mu.RUnlock()

			response.WriteByte(byte(len(players)))
			response.WriteByte(s.cfg.MaxPlayers)
			response.WriteByte(0)   // Bots
			response.WriteByte('d') // Server type
			response.WriteByte('l') // OS
			response.WriteByte(0)   // Visibility
			response.WriteByte(0)   // VAC

			utils.WriteString(&response, "v48")
			response.WriteByte(0x00) // EDF

			conn.WriteToUDP(response.Bytes(), clientAddr)
			fmt.Println("✅ Ответ на A2S_INFO отправлен")
		} else {
			fmt.Printf("⚠️ Невалидный A2S_INFO запрос от %s\n", clientAddr)
		}
	case 0x55: // A2S_PLAYER
		fmt.Printf("📩 A2S_PLAYER от %s\n", clientAddr)
		var response bytes.Buffer
		binary.Write(&response, binary.LittleEndian, int32(-1))

		if len(buffer) >= 9 {
			recvChallenge := int32(binary.LittleEndian.Uint32(buffer[5:9]))
			if recvChallenge == -1 {
				response.WriteByte(0x41)
				binary.Write(&response, binary.LittleEndian, challengeNumber)
				conn.WriteToUDP(response.Bytes(), clientAddr)
				fmt.Println("🔑 Отправлен challenge:", challengeNumber)
			} else if recvChallenge == challengeNumber {
				s.mu.RLock()
				players := make([]Player, len(s.players))
				copy(players, s.players)
				s.mu.RUnlock()

				response.WriteByte(0x44)
				response.WriteByte(byte(len(players)))
				for i, p := range players {
					response.WriteByte(byte(i))
					utils.WriteString(&response, p.Name)
					binary.Write(&response, binary.LittleEndian, p.Score)
					duration := time.Since(p.ConnectedAt)
					binary.Write(&response, binary.LittleEndian, float32(duration.Seconds()))
				}
				conn.WriteToUDP(response.Bytes(), clientAddr)
				fmt.Println("✅ Отправлен список игроков")
			} else {
				fmt.Println("⚠️ Неверный challenge от", clientAddr)
			}
		} else {
			fmt.Printf("⚠️ Недостаточно данных в A2S_PLAYER от %s\n", clientAddr)
		}
	default:
		fmt.Printf("❓ Неизвестный тип запроса (0x%X) от %s\n", header, clientAddr)
	}
}
