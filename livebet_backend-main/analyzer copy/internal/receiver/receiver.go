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

// Сервер для приема данных от парсеров. Теперь он не принимает sourceName.
func (rv *Receiver) StartParserServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[ERROR] Ошибка подключения: %v", err)
			return
		}
		defer conn.Close()

		log.Printf("[DEBUG] Новое подключение от парсера")

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[ERROR] Ошибка чтения сообщения: %v", err)
				break
			}
			rv.receiveChan <- msg
		}
	})

	log.Printf("[DEBUG] Сервер для парсеров запущен на порту %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Printf("[ERROR] Ошибка запуска сервера для парсеров: %v", err)
	}
}
