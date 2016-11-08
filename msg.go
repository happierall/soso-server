package soso

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	actions = map[string]string{
		"get":    "GOT",
		"create": "CREATED",
		"search": "SEARCHED",
		"update": "UPDATED",
		"delete": "DELETED",
		"flush":  "FLUSHED",
	}
)

type Msg struct {
	Model  string
	Action string
	Data   *json.RawMessage
	Log    Log
	Other  *json.RawMessage

	Response *Response

	// Client socket session, public for testing convinience
	Session Session
}

func (c *Msg) Send() {
	c.Response.Log(log_code_by_action_type(c.Action), LevelDebug, "")
	c.sendJSON(c.Response)
}

func (c *Msg) Error(code int, level Level, err error) {
	c.Response.Log(code, level, err.Error())
	c.sendJSON(c.Response)
}

func (c *Msg) Success(Data interface{}) {
	c.Response.Data = Data
	c.Response.Log(log_code_by_action_type(c.Action), LevelDebug, "")

	c.sendJSON(c.Response)
}

func (m *Msg) ReadData(object interface{}) error {
	err := json.Unmarshal(*m.Data, object)
	return err
}

func (m *Msg) ReadOther(object interface{}) error {
	err := json.Unmarshal(*m.Other, object)
	return err
}

func (c *Msg) sendJSON(data interface{}) {
	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	err2 := c.Session.Send(string(json_data))
	if err2 != nil {
		fmt.Println("Error send msg:", err2)
	}
}

func newMsgFromRequest(req *Request, session Session) *Msg {
	msg := &Msg{
		Model:   req.Model,
		Action:  req.Action,
		Session: session,
		Data:    req.Data,
		Log:     req.Log,
		Other:   req.Other,
	}
	msg.Response = NewResponse(msg)

	//Sessions.Push(session, 500) // need for users with token

	return msg
}

func SendMsg(mode, action string, session Session, data map[string]interface{}) {
	msg := &Msg{
		Model:  model,
		Action: action,
		Other:  nil,
	}
	msg.Session = session
	msg.Response = NewResponse(msg)
	msg.Response.Data = data

	msg.Send()
}

func reverse_action_type(action string) string {
	act, ok := actions[action]
	if !ok {
		act = strings.ToUpper(action)
	}
	return act
}

func log_code_by_action_type(action string) int {
	if action == "create" {
		return 201
	}
	return 200
}
