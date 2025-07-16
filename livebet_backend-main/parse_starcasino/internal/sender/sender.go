package sender

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"livebets/parse_starcasino/cmd/config"
	"livebets/shared"
	"sync"
	"time"
)

type Sender struct {
	cfg            config.SenderConfig
	conn           *websocket.Conn
	websocketMutex sync.Mutex
}

func New(
	cfg config.SenderConfig,
) *Sender {
	conn := connectToAnalyzer(cfg)
	return &Sender{
		cfg:            cfg,
		conn:           conn,
		websocketMutex: sync.Mutex{},
	}
}

// Функция подключения к анализатору
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

func (s *Sender) SendMessage(msg shared.GameData) error {
	byteMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	s.websocketMutex.Lock()
	err = s.conn.WriteMessage(websocket.TextMessage, byteMsg)
	s.websocketMutex.Unlock()

	if err != nil {
		return err
	}

	return nil
}
