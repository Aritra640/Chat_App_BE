package main

import (
	ws "github.com/Aritra640/ChatAppBE/WS"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
  e.Use(middleware.CORS())

  e.Any("/connection" , ws.ConnectionHandler)

  e.Logger.Fatal(e.Start(":8080"))
}
