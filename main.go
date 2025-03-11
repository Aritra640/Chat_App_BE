package main

import (
	"flag"
	"log"

	ws "github.com/Aritra640/ChatAppBE/WS"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

  mode := flag.Bool("mode", false , "test json-response")
  testRoom := flag.Bool("testRoom", false , "test room id exists")
  flag.Parse()
  if *mode {
    str := `{"type":"join","payload":{"object":"red"}}`
    res,err := ws.Get_data_resp([]byte(str));if err != nil {
      log.Println("Incorrect format")
      return 
    }

    log.Println("Type: ", res.Type)
    log.Println("Obj: ", res.Payload.Object)

    return 
  }
  
  if *testRoom {
    ws.WSinit()

    ok := ws.RS.CheckRoomIDs(); if ok {
      log.Println("tested successfuly")
      return 
    }else {
      log.Println("test failed")
      return 
    }
  }

  ws.WSinit()
  go ws.RS.Run()

	e := echo.New()
  e.Use(middleware.CORS())

  e.Any("/connection" , ws.ConnectionHandler)

  e.Logger.Fatal(e.Start(":8080"))
}
