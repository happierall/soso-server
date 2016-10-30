package soso

import (
	"net/http"
	"fmt"
)

const (
	Version   string = "2.4.0"
	logPrefix string = "[SoSo]"
)

/*
	ToDo
	1. Add protobuf support
	2. Add save session to soso_test
	3. Auth and save soso_test configs in my server for other users

 */

/*
	Simple Use:
		Router := soso.Default()
		Router.CREATE("message", ChatSendMessage)
		Router.Run(4000)


	Add routes as list:
		var Routes = soso.Routes{}
		Routes.Add("create", "message", ChatSendMessage)

		Router := soso.Default()
		Router.HandleRoutes = Routes
		Router.Run(4000)


	Custom listener:
		Router := soso.Default()
		Router.CREATE("message", ChatSendMessage)
		http.HandleFunc("/soso", Router.receiver)
		http.ListenAndServe("localhost:4000", nil)


	Handler:
		func ChatSendMessage(m *soso.Msg) {

			m.Success(map[string]interface{}{
				"message": "message hi",
				"id": m.RequestMap["id"],
			})

		}


	Send direct message:
		soso.SendMsg("message", "created", session,
			map[string]interface{}{
				"id": "1",
			},
		)

 */

/*
	Client request (javascript):
		var sock = new WebSocket("ws://localhost:4000/soso")
		var data = {
	        data_type: "message",
			action_str: "create",
	        log_map: {},
	        request_map: {msg: "hello world"},
	        trans_map: {}
		}
		sock.send( JSON.stringify( data ) )

 */

type Engine struct {
	Router

	OnOpen func(Session)
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

func(s *Engine) receiver(w http.ResponseWriter, r *http.Request) {
	SosoWebsocketReceiver(w, r, s)
}

func (s *Engine) Run(port int) error {
	http.HandleFunc("/soso", s.receiver)

	return http.ListenAndServe(fmt.Sprintf(":%d", 4000), nil)
}


func New() *Engine {
	soso := Engine{}
	soso.Router = Router{}
	return &soso
}

func Default() *Engine {
	return New()
}
