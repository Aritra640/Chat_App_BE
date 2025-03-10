package ws

import (
	"errors"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type RoomStore struct {
	Rooms map[string]*Room
	Mu    sync.RWMutex //Read Write Mutex for safe access
}

// get a new room store
func NewRoomStore() *RoomStore {
	return &RoomStore{
		Rooms: make(map[string]*Room),
	}
}

// delete a client socket form all rooms (if exists)
func (rs *RoomStore) DeleteAllClients(socket *websocket.Conn) error {
	rs.Mu.RLock()
	defer rs.Mu.RUnlock()

	for roomID := range rs.Rooms {
		err := rs.Rooms[roomID].DeleteClient(socket)
		if err != nil {
			log.Println("Error: cannot delete socket connection in roomID: ", roomID)
			return err
		}
	}

	return nil
}

// add a socket (join) to a roomID (is exists)
func (rs *RoomStore) AddConnectionByID(roomID string, socket *websocket.Conn) error {

	rs.Mu.RLock()
	room, ok := rs.Rooms[roomID]
	rs.Mu.RUnlock()

	if !ok {
		log.Println("cannot find roomID: ", roomID)
		return errors.New("Room not found!")
	}
	err := room.AddClient(socket)
	if err != nil {
		log.Println("Error: cannot add client in room: ", err)
		return err
	}
	return nil
}

//Run add rooms concurrently
func (rs *RoomStore) Run() {
  rs.Mu.RLock()
  defer rs.Mu.RUnlock() 

  for _,room := range rs.Rooms {
    go room.Run()
  }
}

//send a message to all subscribed rooms 
func (rs *RoomStore) SendChatMessage(chat string, socket *websocket.Conn) error {
  rs.Mu.RLock()
  defer rs.Mu.RUnlock() 

  for _,room := range rs.Rooms {
    room.RoomMU.Lock()
    _,ok := room.Clients[socket]
    room.RoomMU.Unlock()
    if ok {
      room.MessageCh <- Message{
        Owner: socket,
        Chat: chat,
      }
    }
  }

  return nil
}


//Add rooms in a rs 
func (rs *RoomStore) AddRoom(roomID string) {
  rs.Rooms[roomID] = NewRoom()
}
