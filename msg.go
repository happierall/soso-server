package soso

import (
	"encoding/json"
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
	Request  *Request
	Response *Response

	User    User
	Context map[string]string

	// Client socket session
	Session Session

	Router *Engine
}

func (c *Msg) Send() {
	if c.Response.Log.LogID == "0" {
		c.Response.NewLog(logCodeByActionType(c.Request.Action), LevelInfo, "")
	}
	c.sendJSON(c.Response)
}

func (c *Msg) Error(code int, level Level, err error) {
	c.Response.NewLog(code, level, err.Error())
	c.sendJSON(c.Response)
}

func (c *Msg) Success(Data interface{}) {
	c.Response.Data = Data
	c.Response.NewLog(logCodeByActionType(c.Request.Action), LevelInfo, "")

	c.sendJSON(c.Response)
}

func (m *Msg) Log(code_key int, lvl_str Level, user_msg string) {
	m.Response.NewLog(code_key, lvl_str, user_msg)
}

func (m *Msg) ReadData(object interface{}) error {
	err := json.Unmarshal(*m.Request.Data, object)
	return err
}

func (m *Msg) ReadOther(object interface{}) error {
	err := json.Unmarshal(*m.Request.Other, object)
	return err
}

func (c *Msg) sendJSON(data interface{}) {
	json_data, err := json.Marshal(data)
	if err != nil {
		Loger.Error(err)
		return
	}

	err2 := c.Session.Send(string(json_data))
	if err2 != nil {
		Loger.Errorf("Error send msg:", err2)
		return
	}
}

func newMsgFromRequest(req *Request, session Session) *Msg {
	msg := &Msg{
		Request: req,
		Session: session,
		Context: make(map[string]string),
	}
	msg.Response = NewResponse(msg)

	//Sessions.Push(session, 500) // need for users with token

	return msg
}

func SendMsg(model, action string, session Session, data map[string]interface{}) {
	msg := &Msg{
		Request: &Request{
			Model:  model,
			Action: action,
			Other:  nil,
		},
		Context: make(map[string]string),
		Session: session,
	}

	msg.Response = NewResponse(msg)
	msg.Response.Other = map[string]interface{}{}
	msg.Response.Data = data

	msg.Send()
}

func reverseActionType(action string) string {
	act, ok := actions[action]
	if !ok {
		act = strings.ToUpper(action)
	}
	return act
}

func logCodeByActionType(action string) int {
	if action == "create" {
		return 201
	}
	return 200
}
