package soso

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	actions = map[string]string{
		"create":   "CREATED",
		"retrieve": "RETRIEVED",
		"update":   "UPDATED",
		"delete":   "DELETED",
		"flush":    "FLUSHED",
	}
)

type Msg struct {
	DataType   string
	ActionStr  string
	LogMap     Log
	RequestMap map[string]interface{}
	TransMap   map[string]interface{}

	Response *Response

	// Client socket session, public for testing convinience
	Session Session

}

func (c *Msg) Send() {
	c.Response.Log(log_code_by_action_type(c.ActionStr), LevelDebug, "")
	c.sendJSON(c.Response)
}

func (c *Msg) Error(code int, level Level, err error) {
	c.Response.Log(code, level, err.Error())
	c.sendJSON(c.Response)
}

func (c *Msg) Success(ResponseMap interface{}) {
	c.Response.ResponseMap = ResponseMap
	c.Response.Log(log_code_by_action_type(c.ActionStr), LevelDebug, "")

	c.sendJSON(c.Response)
}

func (c *Msg) sendJSON(data interface{}) {
	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	err2 := c.Session.Send(string(json_data))
	if (err2 != nil) {
		fmt.Println("Error send msg:", err2)
	}
}

func newMsgFromRequest(req *Request, session Session) *Msg {
	msg := &Msg{
		DataType:   req.DataType,
		ActionStr:  req.ActionStr,
		Session:    session,
		RequestMap: req.RequestMap,
		TransMap:   req.TransMap,
		LogMap:     req.LogMap,
	}
	msg.Response = NewResponse(msg)

	//Sessions.Push(session, 500) // need for users with token

	return msg
}

func SendMsg(dataType, action string, session Session, response map[string]interface{}) {
	msg := &Msg{
		DataType:  dataType,
		ActionStr: action,
		TransMap:   map[string]interface{}{},
	}
	msg.Session = session
	msg.Response = NewResponse(msg)
	msg.Response.ResponseMap = response

	msg.Send()
}

func reverse_action_type(action_str string) string {
	act, ok := actions[action_str]
	if !ok {
		act = strings.ToUpper(action_str)
	}
	return act
}

func log_code_by_action_type(action_str string) int {
	if action_str == "create" {
		return 201
	}
	return 200
}
