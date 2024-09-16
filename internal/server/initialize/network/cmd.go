package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nextlag/in_memory/pkg/logger/l"
)

// CloseFunc - graceful shutdown
type CloseFunc func(context.Context) error

// CMDServer represents a server instance for processing requests.
type CMDServer struct {
	closers []CloseFunc
	log     *l.Logger
}

// NewServer creates and returns a new server instance.
func NewServer(log *l.Logger) (*CMDServer, error) {
	return &CMDServer{
		log: log,
	}, nil
}

// LaunchServer launches the server to receive and process requests from the user.
func (s *CMDServer) LaunchServer(ctx context.Context, handler func(context.Context, []byte, int) string) error {
	var response string
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

		response = handler(ctx, []byte(request), len(request))
		fmt.Println(response)
	}
}

// Close closes the server and its associated resources.
func (s *CMDServer) Close(ctx context.Context) error {
	for _, fn := range s.closers {
		if err := fn(ctx); err != nil {
			s.log.Error("server shutdown error", "err", l.ErrAttr(err))
			return err
		}
	}
	return nil
}
