package soso

import (
	"fmt"
	"net/http"

	"github.com/happierall/l"
)

var (
	Loger = l.New()
)

const (
	Version string = "3.2.0"
)

func init() {
	Loger.Prefix = l.Colorize("SoSo ", l.LightPink)
	Loger.Level = l.LevelInfo
}

/*
	ToDo
	1. Add protobuf support
	2. Add save session to soso_test
	3. Auth and save soso_test configs in my server for other users

*/

type Engine struct {
	Router
}

func (s *Engine) RunReceiver(session Session) {
	Sessions.OnOpenExecute(session)

	// Process incoming messages
	for {
		if msg, err := session.Recv(); err == nil {
			go s.processIncomingMsg(session, msg)
			continue
		}
		break
	}

	Sessions.OnCloseExecute(session)
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

func EnableDebug() {
	Loger.Level = l.LevelDebug
}

func DisableDebug() {
	Loger.Level = l.LevelInfo
}
