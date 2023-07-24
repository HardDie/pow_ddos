package server

import (
	"bytes"
	"context"
	"net"

	"github.com/HardDie/pow_ddos/internal/config"
	"github.com/HardDie/pow_ddos/internal/entity"
	"github.com/HardDie/pow_ddos/internal/logger"
	"github.com/HardDie/pow_ddos/internal/pow"
	"github.com/HardDie/pow_ddos/internal/tcp"
)

type IQuote interface {
	GetRandomQuote() string
}

type Server struct {
	config config.POW
	quote  IQuote
}

func New(config config.POW, quote IQuote) *Server {
	return &Server{
		config: config,
		quote:  quote,
	}
}

func (s *Server) Handle(ctx context.Context, conn net.Conn) {
	// ctx - We should be able to close the application at any time,
	// but read/write operations with ctx.Done() support look terrible,
	// so I decided not to implement them

	// 2. Choose
	reqMsg, err := pow.NewBlock(s.config.MsgSize, s.config.Difficulty)
	if err != nil {
		logger.Error.Println("error create pow block:", err.Error())
		return
	}

	// 3. Challenge (send)
	err = tcp.SendMessage(ctx, conn, entity.MessageTypeChallenge, reqMsg)
	if err != nil {
		logger.Error.Println("error sending challenge message:", err.Error())
		return
	}

	// 5. Response (receive)
	respMsg, err := tcp.ReceiveMessage[pow.Block](ctx, conn, entity.MessageTypeResponse)
	if err != nil {
		logger.Error.Println("error receiving response:", err.Error())
		return
	}

	// 6. Verify
	reqMsg.Nonce = respMsg.Nonce
	if !bytes.Equal(reqMsg.CalcHash(), respMsg.CalcHash()) {
		logger.Debug.Println("Invalid POW")
		return
	}

	// 7. Grant service
	quote := entity.QuoteMessage{
		Quote: s.quote.GetRandomQuote(),
	}
	err = tcp.SendMessage(ctx, conn, entity.MessageTypeQuote, quote)
	if err != nil {
		logger.Error.Println("error sending quote:", err.Error())
		return
	}
}
