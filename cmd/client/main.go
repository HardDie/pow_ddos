package main

import (
	"context"
	"net"
	"runtime"

	"github.com/HardDie/pow_ddos/internal/config"
	"github.com/HardDie/pow_ddos/internal/entity"
	"github.com/HardDie/pow_ddos/internal/logger"
	"github.com/HardDie/pow_ddos/internal/pow"
	"github.com/HardDie/pow_ddos/internal/tcp"
)

func Run() error {
	ctx := context.Background()
	cfg := config.GetClient()

	// 1. Request service
	conn, err := net.Dial("tcp", cfg.Client.Host)
	if err != nil {
		logger.Error.Println("error dial server:", err.Error())
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error.Println("error close tcp connection:", err.Error())
		}
	}()

	// 3. Challenge (receive)
	logger.Debug.Println("Start receiving challenge...")
	block, err := tcp.ReceiveMessage[pow.Block](ctx, conn, entity.MessageTypeChallenge)
	if err != nil {
		logger.Error.Println("error receiving challenge:", err.Error())
		return err
	}

	// 4. Solve
	logger.Debug.Println("Started POW calculation...")
	err = block.POW(runtime.NumCPU())
	if err != nil {
		logger.Error.Println("error calculate POW:", err.Error())
		return err
	}

	// 5. Response (send)
	err = tcp.SendMessage(ctx, conn, entity.MessageTypeResponse, block)
	if err != nil {
		logger.Error.Println("error sending response:", err.Error())
		return err
	}

	// Receive quote
	logger.Debug.Println("Start receiving quote...")
	quote, err := tcp.ReceiveMessage[entity.QuoteMessage](ctx, conn, entity.MessageTypeQuote)
	if err != nil {
		logger.Error.Println("error receiving quote:", err.Error())
		return err
	}
	logger.Info.Println("Received quote:", quote.Quote)
	return nil
}

func main() {
	if err := Run(); err != nil {
		logger.Error.Fatal("client error:", err.Error())
	}
}
