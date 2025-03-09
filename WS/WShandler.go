package WS

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func WShandler(c echo.Context) error {

  ws,err := upgrader.Upgrade(c.Response() , c.Request() , c.Response().Header())
  if err != nil {
    
    log.Println("Error in upgrading websocket: " , err)
    return c.JSON(http.StatusInternalServerError , "something went wrong")
  }

  defer ws.Close()


}
