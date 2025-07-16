package receiver

import (
	"fmt"
	"livebets/analazer/internal/entity"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

type Receiver struct {
	receiveChan chan<- entity.ReceivedMsg
}

func NewReceiver(receiveChan chan<- entity.ReceivedMsg) *Receiver {
	return &Receiver{receiveChan: receiveChan}
}

// Сервер для приема данных от парсеров
func (rv *Receiver) StartParserServer(port int, sourceName string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[ERROR] Ошибка подключения от %s: %v", sourceName, err)
			return
		}
		defer conn.Close()

		log.Printf("[DEBUG] Новое подключение от парсера %s", sourceName)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[ERROR] Ошибка чтения сообщения от %s: %v", sourceName, err)
				break
			}
			// log.Printf("[DEBUG] Получены данные от %s: %s", sourceName, string(msg))

			// CANNEL
			rv.receiveChan <- msg
		}
	})

	log.Printf("[DEBUG] Сервер для парсера %s запущен на порту %d", sourceName, port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Printf("[ERROR] Ошибка запуска сервера для парсера %s: %v", sourceName, err)
	}
}
