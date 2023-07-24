package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/HardDie/pow_ddos/internal/config"
	"github.com/HardDie/pow_ddos/internal/logger"
	"github.com/HardDie/pow_ddos/internal/quote"
	"github.com/HardDie/pow_ddos/internal/server"
	"github.com/HardDie/pow_ddos/internal/tcp"
)

func Run() error {
	cfg := config.GetServer()
	ctx := context.Background()

	q, err := quote.New(cfg.Quote)
	if err != nil {
		logger.Error.Println("error init quote:", err.Error())
		return err
	}

	s := server.New(cfg.POW, q)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	srv, err := tcp.NewServer(cfg.Server.Port, s)
	if err != nil {
		return err
	}

	go func() {
		<-done
		err := srv.Close()
		if err != nil {
			logger.Error.Println("error close server after signal:", err.Error())
		}
	}()

	logger.Info.Println("Server started on", cfg.Server.Port)
	return srv.Serve(ctx)
}

func main() {
	if err := Run(); err != nil {
		logger.Error.Fatal("server error:", err.Error())
	}
}
