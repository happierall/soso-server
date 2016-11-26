package soso

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

   resp := Response{
      Action: "RETRIEVED",
      model:  "person",
      Log:    log,
      Data:   {},
      Other:  {}
   }

    Other:

    Other: map[string]interface{}{
      AuthToken: "c76aa3577f8b5a60206f9d041c76034a",
      TransId:   "eb99ec08-7e90-400d-9585-62a1385ec158"
    }
*/

// direct and indirect responses
type Response struct {
	Model  string      `json:"model"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
	Log    Log         `json:"log"`
	Other  interface{} `json:"other"`
}

func NewResponse(msg *Msg) *Response {
	return &Response{
		Action: reverseActionType(msg.Request.Action),
		Model:  msg.Request.Model,
		Other:  msg.Request.Other,
	}
}

func (r *Response) NewLog(code_key int, lvl_str Level, user_msg string) *Response {
	r.Log = NewLog(code_key, lvl_str, user_msg)
	return r
}

func (r Response) Result() Response { return r }
