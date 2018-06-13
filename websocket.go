package soso

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"errors"
	"time"
)

type envelope struct {
	t   int
	msg []byte
}

var (
	WebSocketReadBufSize  = 1024
	WebSocketWriteBufSize = 1024

	WriteWait = 10 * time.Second // Milliseconds until write times out.
	PongWait  = 60 * time.Second  // Timeout for waiting on pong.

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10 // Milliseconds between pings.

	MaxMessageSize    int64 = 512
	MessageBufferSize       = 256
)

var mux = sync.RWMutex{}
var nextId int64 = 1

// Simple ID scheme
// For anonymity or overflow reason something more elaborate could be needed
func getId() string {
	mux.Lock()
	defer mux.Unlock()

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
	sessID := getId()

	session := newSession(req, conn, sessID)
	defer session.close(0, "")

	go session.writePump()

	engine.RunReceiver(session)
}

type websocketSession struct {
	sync.RWMutex
	id       string
	req      *http.Request
	conn     *websocket.Conn
	isClosed bool
	output   chan *envelope
}

func newSession(req *http.Request, conn *websocket.Conn, sessionID string) *websocketSession {
	s := &websocketSession{
		id:       sessionID,
		req:      req,
		conn:     conn,
		isClosed: false,
		output:   make(chan *envelope, MessageBufferSize),
	}
	return s
}

func (s *websocketSession) ID() string {
	s.RLock()
	defer s.RUnlock()

	return s.id
}

func (s *websocketSession) Recv() ([]byte, error) {
	s.conn.SetReadLimit(MaxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(PongWait))
	s.conn.SetPongHandler(func(string) error { s.conn.SetReadDeadline(time.Now().Add(PongWait)); return nil })

	mt, msg, err := s.conn.ReadMessage()
	if mt != websocket.TextMessage && mt != -1 {
		Loger.Warnf("only text can be sent. MsgType = %d. Msg = %s\n", mt, msg)
		return nil, errors.New("only text can be sent")
	}
	if mt == -1 {
		s.close(1, "client was closed")
	}

	return msg, err
}

func (s *websocketSession) Send(msg string) error {
	if s.IsClosed() {
		return errors.New("session already closed")
	}
	return s.writeMessage(&envelope{websocket.TextMessage, []byte(msg)})
}

func (s *websocketSession) SendBinary(msg []byte) error {
	if s.IsClosed() {
		return errors.New("session already closed")
	}
	return s.writeMessage(&envelope{websocket.TextMessage, msg})
}

func (s *websocketSession) Close(status uint32, reason string) error {
	if s.IsClosed() {
		return errors.New("already closed session")
	}

	s.writeMessage(&envelope{t: websocket.CloseMessage, msg: []byte{}})

	return nil
}

func (s *websocketSession) close(status uint32, reason string) error {
	if s.IsClosed() {
		return errors.New("already closed session")
	}

	s.Lock()
	defer s.Unlock()

	s.isClosed = true
	s.conn.Close()
	close(s.output)
	return nil
}

func (s *websocketSession) IsClosed() bool {
	s.RLock()
	defer s.RUnlock()
	return s.isClosed
}

func (s *websocketSession) writeMessage(message *envelope) error {
	if s.IsClosed() {
		return errors.New("session already closed")
	}

	select {
	case s.output <- message:
	default:
		return errors.New("session message buffer is full")
	}
	return nil
}

func (s *websocketSession) writeRaw(message *envelope) error {
	if s.IsClosed() {
		return errors.New("tried to write to a closed session")
	}

	s.conn.SetWriteDeadline(time.Now().Add(WriteWait))
	err := s.conn.WriteMessage(message.t, message.msg)

	if err != nil {
		return err
	}

	return nil
}

func (s *websocketSession) ping() {
	s.writeRaw(&envelope{t: websocket.PingMessage, msg: []byte{}})
}

func (s *websocketSession) writePump() {
	ticker := time.NewTicker(PingPeriod)
	defer ticker.Stop()

loop:
	for {
		select {
		case msg, ok := <-s.output:
			if !ok {
				Loger.Warn("Close session")
				break loop
			}

			err := s.writeRaw(msg)

			if err != nil {
				Loger.Error(err)
				break loop
			}

			if msg.t == websocket.CloseMessage {
				break loop
			}
		case <-ticker.C:
			s.ping()
		}
	}
}
