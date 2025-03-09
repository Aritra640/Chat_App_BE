package WS

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

type Room struct {
	RoomID  string
	MsgChan chan Message
	Clients map[*websocket.Conn]bool
	RoomMu  sync.Mutex
	Stop    chan struct{}
}

type Message struct {
	owner *websocket.Conn
	chat  string
}

func NewRoom(room string) *Room {
	return &Room{
		RoomID:  room,
		MsgChan: make(chan Message),
		Clients: make(map[*websocket.Conn]bool),
		Stop:    make(chan struct{}),
	}
}

var Rooms = map[string]*Room{
	"red":    NewRoom("red"),
	"violet": NewRoom("violet"),
}

func (r *Room) Run() {
	for {
		select {
		case msg := <-r.MsgChan:
			r.RoomMu.Lock()

			for client := range r.Clients {
        
        if client == msg.owner {
          continue
        }

				err := client.WriteMessage(websocket.TextMessage, []byte(msg.chat))
				if err != nil {

					log.Println("Error: ", err)
					client.Close()
					delete(r.Clients, client)
				}
			}
      
      r.RoomMu.Unlock()

		case <-r.Stop:
      log.Println("Stopping.... " , r.RoomID)
      return 
		}
	}
}
