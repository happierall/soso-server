package soso

import (
	"encoding/json"
)

/*
   Example:

   log := Log{
       CodeKey: "400",
       CodeStr: "Bad request",

       LevelInt: int(Level(3)), // or int(LevelError)
       LevelStr: Level(3), // or obj of LevelError

       LogID: "1096",
       UserMsg: "action required"
   }

   resp := Request{
      Action: "get",
      Model: "person",
      Log: log,
      Data: {},
      Other: {}
   }


    Other: {
      AuthToken: "c76aa3577f8b5a60206f9d041c76034a...",
      TransID: "eb99ec08-7e90-400d-9585-62a1385ec158"
    }

*/

// Request
type Request struct {
	Action string           `json:"action"`
	Model  string           `json:"model"`
	Log    Log              `json:"log"`
	Data   *json.RawMessage `json:"data"`
	Other  *json.RawMessage `json:"other"`
}

func NewRequest(msg []byte) (*Request, error) {
	var req *Request
	err := json.Unmarshal(msg, &req)
	return req, err
}
