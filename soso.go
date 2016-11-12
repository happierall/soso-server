package soso

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

var (
	Loger = log.New()
)

const (
	Version   string = "3.1.0"
	logPrefix string = "[SoSo]"
)

func init() {
	Loger.Level = log.InfoLevel
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
	Loger.Level = log.DebugLevel
}

func DisableDebug() {
	Loger.Level = log.InfoLevel
}
