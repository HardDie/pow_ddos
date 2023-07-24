package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/HardDie/pow_ddos/internal/logger"
)

type ServerConfig struct {
	POW    POW
	Server Server
	Quote  Quote
}

func GetServer() ServerConfig {
	if err := godotenv.Load(); err != nil {
		if check := os.IsNotExist(err); !check {
			logger.Error.Fatalf("failed to load env vars: %s", err)
		}
	}

	cfg := ServerConfig{
		POW:    powConfig(),
		Server: serverConfig(),
		Quote:  quoteConfig(),
	}
	return cfg
}

type ClientConfig struct {
	Client Client
}

func GetClient() ClientConfig {
	if err := godotenv.Load(); err != nil {
		if check := os.IsNotExist(err); !check {
			logger.Error.Fatalf("failed to load env vars: %s", err)
		}
	}

	cfg := ClientConfig{
		Client: clientConfig(),
	}
	return cfg
}
