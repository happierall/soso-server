package soso

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

var WebSocketReadBufSize = 4096
var WebSocketWriteBufSize = 4096

var nextId int64 = 1

// Simple ID scheme
// For anonymity or overflow reason something more elaborate could be needed
func getId() string {
	nextId += 1
	return strconv.FormatInt(nextId, 10)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  WebSocketReadBufSize,
	WriteBufferSize: WebSocketWriteBufSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SosoWebsocketReceiver(rw http.ResponseWriter, req *http.Request, engine *Engine) {
	conn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	sessID := getId()

	session := newSession(req, conn, sessID)

	engine.RunReceiver(session)
}

type websocketSession struct {
	sync.RWMutex
	id       string
	req      *http.Request
	conn     *websocket.Conn
	isClosed bool
}

func newSession(req *http.Request, conn *websocket.Conn, sessionID string) *websocketSession {
	s := &websocketSession{
		id:       sessionID,
		req:      req,
		conn:     conn,
		isClosed: false,
	}
	return s
}

func (s *websocketSession) ID() string {
	return s.id
}

func (s *websocketSession) Recv() ([]byte, error) {

	mt, msg, err := s.conn.ReadMessage()
	if mt != websocket.TextMessage && mt != -1 {
		Loger.Warnf("only text can be sent. MsgType = %d. Msg = %s\n", mt, msg)
	}
	if mt == -1 {
		s.Close(1, "client was closed")
	}

	return msg, err
}

func (s *websocketSession) Send(msg string) error {
	s.Lock()
	defer s.Unlock()
	return s.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (s *websocketSession) Close(status uint32, reason string) error {
	s.isClosed = true
	s.conn.Close()
	return nil
}

func (s *websocketSession) IsClosed() bool {
	return s.isClosed
}
