package client

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/nextlag/in_memory/pkg/cleanenv"
)

var (
	cfg  Config
	once sync.Once
	Env  = ".env"
	Yaml = "config/client/config.yaml"
)

type (
	Config struct {
		CfgYAML string   `yaml:"config_yaml" env:"CONFIG_CLIENT_YAML"`
		Network *Network `yaml:"network"`
	}

	Network struct {
		ServerAddress  string        `yaml:"server_address"`
		MaxConnections int           `yaml:"max_connections"`
		MaxMessageSize string        `yaml:"max_message_size"`
		IdleTimeout    time.Duration `yaml:"idle_timeout"`
	}
)

func Load() *Config {
	once.Do(func() {
		err := godotenv.Load(Env)
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		if err = cleanenv.ReadConfig(Yaml, &cfg); err != nil {
			log.Fatalf("error load config.yaml: %v", err)
		}

		flag.StringVar(&cfg.Network.ServerAddress, "addr", cfg.Network.ServerAddress, "Address TCP server")

		if err = env.Parse(&cfg); err != nil {
			log.Fatalf("error parsing .env variables: %v", err)
		}

		flag.Parse()
	})
	return &cfg
}
