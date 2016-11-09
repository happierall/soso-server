package soso

import (
	"fmt"
	"net/http"
)

const (
	Version   string = "3.0.0"
	logPrefix string = "[SoSo]"
)

/*
	ToDo
	1. Add protobuf support
	2. Add save session to soso_test
	3. Auth and save soso_test configs in my server for other users

*/

type Engine struct {
	Router

	OnOpen  func(Session)
	OnClose func(Session)
}

func (s *Engine) RunReceiver(session Session) {
	if s.OnOpen != nil {
		s.OnOpen(session)
	}

	// Process incoming messages
	for {
		if msg, err := session.Recv(); err == nil {
			go s.processIncomingMsg(session, msg)
			continue
		}
		break
	}

	if s.OnClose != nil {
		fmt.Println("close onClose")
		s.OnClose(session)
	}
}

func (s *Engine) Receiver(w http.ResponseWriter, r *http.Request) {
	SosoWebsocketReceiver(w, r, s)
}

func (s *Engine) Run(port int) error {
	http.HandleFunc("/soso", s.Receiver)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func New() *Engine {
	soso := Engine{}
	soso.Router = Router{}
	return &soso
}

func Default() *Engine {
	return New()
}
