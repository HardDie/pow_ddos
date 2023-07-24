package config

import (
	"os"
	"strconv"

	"github.com/HardDie/pow_ddos/internal/logger"
)

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		logger.Error.Fatalf("env %q value not found", key)
	}
	return value
}

func getEnvAsInt(key string) int {
	value := getEnv(key)
	v, e := strconv.Atoi(value)
	if e != nil {
		logger.Error.Fatalf("env %q value invalid int", key)
	}
	return v
}
