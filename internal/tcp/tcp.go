package tcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/HardDie/pow_ddos/internal/entity"
	"github.com/HardDie/pow_ddos/internal/logger"
)

type IHandler interface {
	Handle(ctx context.Context, conn net.Conn)
}

type Server struct {
	l net.Listener
	h IHandler

	sync.Mutex
}

func NewServer(addr string, h IHandler) (*Server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		l: l,
		h: h,
	}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	defer func() {
		if err := s.Close(); err != nil {
			logger.Error.Println("error closing tcp socket:", err.Error())
		}
	}()

	for {
		conn, err := s.l.Accept()
		if err != nil {
			return err
		}
		go s.h.Handle(ctx, conn)
	}
}

func (s *Server) Close() error {
	s.Lock()
	defer s.Unlock()

	if s.l != nil {
		err := s.l.Close()
		s.l = nil
		return err
	}
	return nil
}

func SendMessage(_ context.Context, conn net.Conn, msgType int, msg any) error {
	sendMsg := &entity.MessageSend{
		Type: msgType,
		Data: msg,
	}
	data, err := json.Marshal(sendMsg)
	if err != nil {
		return fmt.Errorf("error marshal msg: %w", err)
	}
	n, err := conn.Write(data)
	if err != nil {
		return fmt.Errorf("error sending data: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("error send bytes %d; want %d", n, len(data))
	}
	return nil
}

func ReceiveMessage[T any](_ context.Context, conn net.Conn, msgType int) (*T, error) {
	var msg entity.Message
	var resp T

	// Decode message
	err := json.NewDecoder(conn).Decode(&msg)
	if err != nil {
		return nil, fmt.Errorf("error decoding msg: %w", err)
	}
	// Validate message type
	if msg.Type != msgType {
		return nil, fmt.Errorf("invalid message type %d; want %d", msg.Type, msgType)
	}
	// Decode data
	err = json.Unmarshal(msg.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("error decoding data: %w", err)
	}

	return &resp, nil
}
