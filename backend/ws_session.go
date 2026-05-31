package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsReadWait  = 10 * time.Minute
	wsWriteWait = 30 * time.Second
	wsPingEvery = 25 * time.Second
)

type wsSession struct {
	conn   *websocket.Conn
	textCh chan string
	binCh  chan []byte
	errCh  chan error
	done   chan struct{}
	once   sync.Once
}

func newWSSession(c *websocket.Conn) *wsSession {
	s := &wsSession{
		conn:   c,
		textCh: make(chan string, 8),
		binCh:  make(chan []byte, 4),
		errCh:  make(chan error, 1),
		done:   make(chan struct{}),
	}
	c.SetPongHandler(func(string) error {
		return c.SetReadDeadline(time.Now().Add(wsReadWait))
	})
	go s.readLoop()
	go s.pingLoop()
	return s
}

func (s *wsSession) close() {
	s.once.Do(func() { close(s.done) })
}

func (s *wsSession) readLoop() {
	defer s.close()
	for {
		if err := s.conn.SetReadDeadline(time.Now().Add(wsReadWait)); err != nil {
			s.errCh <- err
			return
		}
		msgType, data, err := s.conn.ReadMessage()
		if err != nil {
			s.errCh <- err
			return
		}
		switch msgType {
		case websocket.TextMessage:
			msg := string(data)
			if msg == "KEEPALIVE" {
				continue
			}
			select {
			case s.textCh <- msg:
			case <-s.done:
				return
			}
		case websocket.BinaryMessage:
			select {
			case s.binCh <- data:
			case <-s.done:
				return
			}
		}
	}
}

func (s *wsSession) pingLoop() {
	ticker := time.NewTicker(wsPingEvery)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			if err := s.conn.SetWriteDeadline(time.Now().Add(wsWriteWait)); err != nil {
				return
			}
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *wsSession) sendText(msg string) error {
	if err := s.conn.SetWriteDeadline(time.Now().Add(wsWriteWait)); err != nil {
		return err
	}
	return s.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (s *wsSession) recvText() (string, error) {
	select {
	case msg := <-s.textCh:
		return msg, nil
	case err := <-s.errCh:
		return "", err
	case <-s.done:
		return "", fmt.Errorf("connection closed")
	}
}
