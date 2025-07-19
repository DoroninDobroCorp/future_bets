package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"livebets/analazer/internal/entity"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type connFilter struct {
	conn   *websocket.Conn
	filter GeneralFilter
}

type client struct {
	filter    GeneralFilter
	isClosing bool
	mu        sync.Mutex
}

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

type Sender struct {
	clients    map[*websocket.Conn]*client
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	receive    chan connFilter
	broadcast  <-chan []entity.ResponsePair
	logger     *zerolog.Logger
}

func NewSender(broadcast <-chan []entity.ResponsePair, logger *zerolog.Logger) *Sender {
	clients := make(map[*websocket.Conn]*client)
	register := make(chan *websocket.Conn)
	unregister := make(chan *websocket.Conn)
	receive := make(chan connFilter)
	return &Sender{
		clients:    clients,
		register:   register,
		unregister: unregister,
		receive:    receive,
		broadcast:  broadcast,
		logger:     logger,
	}
}

func (s *Sender) RunHub(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case connection := <-s.register:
			s.clients[connection] = &client{}
			s.logger.Info().Msgf("[Sender.RunHub] connection registered - %v", connection.RemoteAddr())
		case connFilter := <-s.receive:
			client, ok := s.clients[connFilter.conn]
			if ok {
				client.filter = connFilter.filter
				s.clients[connFilter.conn] = client
			}
		case data := <-s.broadcast:
			// Send the message to all clients
			for connection, c := range s.clients {
				go func(connection *websocket.Conn, c *client) { // send to each client in parallel so we don't block on a slow client
					c.mu.Lock()
					defer c.mu.Unlock()
					if c.isClosing {
						return
					}

					filterData := s.Filter(data, c.filter)

					respData, err := json.Marshal(filterData)
					if err != nil {
						s.logger.Error().Err(err).Msgf("[Sender.RunHub] marshal error - %v", respData)
					}

					if err := connection.WriteMessage(websocket.TextMessage, respData); err != nil {
						c.isClosing = true
						s.logger.Error().Err(err).Msgf("[Sender.RunHub] write error - %v", connection.RemoteAddr())

						connection.WriteMessage(websocket.CloseMessage, []byte{})
						connection.Close()
						s.unregister <- connection
					}
				}(connection, c)
			}

		case connection := <-s.unregister:
			// Remove the client from the hub
			delete(s.clients, connection)

			s.logger.Info().Msgf("[Sender.RunHub] connection unregistered - %v", connection.RemoteAddr())
		case <-ctx.Done():
			return
		}
	}
}

// Сервер для фронтенда
func (s *Sender) StartServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/output", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.logger.Error().Err(err).Msgf("[Sender.StartServer] connect error")
			return
		}

		s.register <- conn
		defer func() {
			s.unregister <- conn
			conn.Close()
		}()

		for {
			_, filter, err := conn.ReadMessage()
			if err != nil {
				s.logger.Error().Err(err).Msgf("[Sender.StartServer] read message error")
				break
			}

			var connFilter connFilter
			connFilter.conn = conn
			err = json.Unmarshal(filter, &connFilter.filter)
			if err != nil {
				s.logger.Error().Err(err).Msgf("[Sender.StartServer] json filter unmarshall error")
			}

			s.receive <- connFilter
		}
	})

	s.logger.Info().Msgf("[Sender.StartServer] start server port - %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		s.logger.Error().Err(err).Msgf("[Sender.StartServer] start frontend server error")
	}
}
