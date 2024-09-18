package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Parallel()
	Env = "../../.env"
	Yaml = "config.yaml"
	cfg = *Load()
	cfg.Logging.ProjectPath = "."
	tests := map[string]struct {
		expectedCfg Config
	}{
		"load config": {
			expectedCfg: Config{
				CfgYAML: "config/server/config.yaml",
				Engine: &Engine{
					Type:             "in_memory",
					PartitionsNumber: 8,
				},
				Server: &Server{
					SrvType: "tcp",
				},
				Network: &Network{
					TCPSocket:      ":9080",
					MaxConnections: 100,
					MaxMessageSize: "4Kb",
					IdleTimeout:    time.Minute * 5,
				},
				Logging: &Logging{
					Level:       -4,
					ProjectPath: ".",
					LogToFile:   false,
					LogPath:     "data/logs/out.log",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedCfg, cfg)
		})
	}
}
