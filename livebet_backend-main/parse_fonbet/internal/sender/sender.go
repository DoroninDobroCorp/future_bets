package sender

import (
	"context"
	"encoding/json"
	"livebets/parse_fonbet/cmd/config"
	"livebets/parse_fonbet/internal/entity"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Sender struct {
	cfg      config.SenderConfig
	conn     *websocket.Conn
	sendChan <-chan entity.ResponseGame
}

func NewSender(cfg config.SenderConfig, sendChan <-chan entity.ResponseGame) *Sender {
	conn := connectToAnalyzer(cfg)
	return &Sender{
		cfg:      cfg,
		conn:     conn,
		sendChan: sendChan,
	}
}

func connectToAnalyzer(cfg config.SenderConfig) *websocket.Conn {
	var analyzerConnection *websocket.Conn
	var err error
	for {
		analyzerConnection, _, err = websocket.DefaultDialer.Dial(cfg.Url, nil)
		if err != nil {
			log.Printf("[ERROR] Ошибка подключения к анализатору: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	return analyzerConnection
}

func (s *Sender) SendingToAnalyzer(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		select {
		case data := <-s.sendChan:
			log.Printf("{%v}", data)
			byteMsg, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				return err
			}

			if err := s.conn.WriteMessage(websocket.TextMessage, byteMsg); err != nil {
				log.Printf("[ERROR] Ошибка отправки данных клиенту (%v): %v", s.conn.RemoteAddr(), err)
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}
