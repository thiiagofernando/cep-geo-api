package config

import (
	"os"
	"strconv"
)

// Config é a estrutura de configuração da aplicação
type Config struct {
	Port int
}

// LoadConfig carrega as configurações da aplicação
func LoadConfig() *Config {
	port := 8070

	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	return &Config{
		Port: port,
	}
}
