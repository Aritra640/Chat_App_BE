package ws

import (
	"encoding/json"
	"log"
)

type Response struct {
  Type string `json:"type"`
  Payload Payload `json:"payload"`
}

type Payload struct {
  Object string `json:"object"`
}

func get_data_resp(str []byte) (Response , error) {
  
  var res Response
  err := json.Unmarshal(str, &res); if err != nil {
    log.Println("Error: cannot transform data in response")
    return Response{} , err
  }

  return res , nil
}
