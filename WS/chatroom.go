package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Owner *websocket.Conn
	Chat  string
}

type Room struct {
	MessageCh chan Message
	Clients   map[*websocket.Conn]bool
	RoomMU    sync.Mutex
	Stop      chan struct{}
  WriteMu   sync.Mutex
}

func NewRoom() *Room {
	return &Room{
		MessageCh: make(chan Message),
		Clients:   make(map[*websocket.Conn]bool),
		Stop:      make(chan struct{}),
	}
}

func (r *Room) DeleteClient(socket *websocket.Conn) error {

	r.RoomMU.Lock()
	_, ok := r.Clients[socket]
	if ok {
		delete(r.Clients, socket)
	}
	r.RoomMU.Unlock()
	return nil
}

func (r *Room) AddClient(socket *websocket.Conn) error {

	r.RoomMU.Lock()
	r.Clients[socket] = true
	r.RoomMU.Unlock()

	return nil
}

func (r *Room) Run() {
	for {
		select {
		case msg := <-r.MessageCh:
			r.RoomMU.Lock()

			for client := range r.Clients {
        r.WriteMu.Lock()
				err := client.WriteMessage(websocket.TextMessage, []byte(msg.Chat))
        r.WriteMu.Unlock()
				if err != nil {
					log.Println("Error: ", err)
					client.Close()
					r.DeleteClient(client)
				}
			}

			r.RoomMU.Unlock()

		case <-r.Stop:
			log.Println("Stopping...")
			return
		}
	}
}
