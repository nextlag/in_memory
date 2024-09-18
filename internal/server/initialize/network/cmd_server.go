package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/nextlag/in_memory/internal"
	"github.com/nextlag/in_memory/pkg/logger/l"
)

type CMDServer struct {
	log *l.Logger
}

func NewCMDServer(log *l.Logger) (*CMDServer, error) {
	return &CMDServer{
		log: log,
	}, nil
}

func (s *CMDServer) LaunchServer(ctx context.Context, handler Handler) error {
	var response []byte
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != nil {
				return errors.New("error context closing")
			}
			return nil
		default:
		}
		fmt.Print("> ")
		request, err := reader.ReadString('\n')
		if err != nil {
			s.log.Error("error read string", "err", l.ErrAttr(err))
			continue
		}
		request = strings.TrimSpace(request)
		if len(request) == 0 {
			continue
		}

		response = handler(ctx, []byte(request))
		switch string(response) {
		case internal.ResponseOk:
			color.Green(string(response))
		default:
			color.Red(string(response))
		}
	}
}

// Close closes the server and its associated resources.
func (s *CMDServer) Close() {
	s.log.Info("CMDServer is shutting down")
	return
}
