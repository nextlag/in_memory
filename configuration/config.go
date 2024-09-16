package configuration

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/nextlag/in_memory/pkg/cleanenv"
	"log"
	"log/slog"
	"sync"
)

var (
	cfg  Config
	once sync.Once
	Env  = ".env"
	Yaml = "configuration/config.yaml"
)

type (
	Config struct {
		CfgYAML string   `yaml:"config_yaml" env:"CONFIG_YAML"`
		Server  *Server  `yaml:"server"`
		Engine  *Engine  `yaml:"engine"`
		Logging *Logging `yaml:"logging"`
	}

	Server struct {
		SrvType string `yaml:"server_type"` // cmd or tcp
	}

	Engine struct {
		Type             string `yaml:"type"`
		PartitionsNumber uint   `yaml:"partitions_number"`
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

		flag.Var(&LogLevelValue{Value: &cfg.Logging.Level}, "level", "Log level (debug, info, warn, error)")

		if err = env.Parse(&cfg); err != nil {
			log.Fatalf("error parsing .env variables: %v", err)
		}

		flag.Parse()
	})
	return &cfg
}
