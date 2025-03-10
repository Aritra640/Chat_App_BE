package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var RS RoomStore

//Initialise ws room-store and add rooms to it
func WSinit() {
  RS = *NewRoomStore()
  RS.AddRoom("Red")
  RS.AddRoom("Violet")
}

func ConnectionHandler(c echo.Context) error {
  ws,err := upgrader.Upgrade(c.Response() , c.Request() , c.Response().Header()); if err != nil {
    log.Println("Error: cannot upgrader websocket : " , err)
    return c.JSON(http.StatusInternalServerError , map[string]string{
      "message": "something went wrong1",
    })
  }
  defer ws.Close()

  for{

    _,msg,err := ws.ReadMessage(); if err != nil {
      RS.DeleteAllClients(ws)
      log.Println("Error: client disconnected: " , err)
      break
    }

    //now deal with the message json-tag
    res,err := get_data_resp(msg); if err != nil {
      ws.WriteMessage(websocket.TextMessage , []byte("request format is not correct (System Message)"))
      continue
    }
    
    if res.Type == "join" {
      RS.AddConnectionByID(res.Payload.Object , ws)
    }else {
      RS.SendChatMessage(res.Payload.Object , ws)
    }
  }

  return c.JSON(http.StatusOK, map[string]string{
    "message": "Webscket server stopped",
  })
}
