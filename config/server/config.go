package server

import (
	"flag"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/nextlag/in_memory/config"
	"github.com/nextlag/in_memory/pkg/cleanenv"
)

var (
	cfg  Config
	once sync.Once
	Env  = ".env"
	Yaml = "config/server/config.yaml"
)

type (
	Config struct {
		CfgYAML string   `yaml:"config_yaml" env:"CONFIG_SERVER_YAML"`
		Server  *Server  `yaml:"server"`
		Engine  *Engine  `yaml:"engine"`
		Network *Network `yaml:"network"`
		Logging *Logging `yaml:"logging"`
	}

	Server struct {
		SrvType string `yaml:"server_type"` // cmd or tcp
	}

	Engine struct {
		Type             string `yaml:"type"`
		PartitionsNumber int    `yaml:"partitions_number"`
	}

	Network struct {
		TCPSocket      string        `yaml:"tcp_socket"`
		MaxConnections int           `yaml:"max_connections"`
		MaxMessageSize string        `yaml:"max_message_size"`
		IdleTimeout    time.Duration `yaml:"idle_timeout"`
	}

	Logging struct {
		Level       slog.Level `yaml:"level"`
		ProjectPath string     `yaml:"project_path" env:"PROJECT_PATH"`
		LogToFile   bool       `yaml:"log_to_file"`
		LogPath     string     `yaml:"log_path"`
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

		flag.StringVar(&cfg.Network.TCPSocket, "addr", cfg.Network.TCPSocket, "Host TCP server")
		flag.Var(&config.LogLevelValue{Value: &cfg.Logging.Level}, "level", "Log level (debug, info, warn, error)")

		if err = env.Parse(&cfg); err != nil {
			log.Fatalf("error parsing .env variables: %v", err)
		}

		flag.Parse()
	})
	return &cfg
}
