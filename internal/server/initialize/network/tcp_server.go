package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nextlag/in_memory/pkg/logger/l"
)

type TCPServer struct {
	listener net.Listener

	closing        atomic.Bool
	idleTimeout    time.Duration
	bufferSize     int
	maxConnections int

	log *l.Logger
}

func NewTCPServer(address string, log *l.Logger, options ...TCPServerOption) (Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	server := &TCPServer{
		listener: listener,
		log:      log,
	}

	for _, opt := range options {
		opt(server)
	}

	if server.bufferSize == 0 {
		server.bufferSize = 4 << 10
	}

	return server, nil
}

func (s *TCPServer) LaunchServer(ctx context.Context, handler Handler) error {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			connection, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				s.log.Error("failed to accept", l.ErrAttr(err))
				continue
			}

			wg.Add(1)
			go func(connection net.Conn) {
				defer wg.Done()
				s.handleConnection(ctx, connection, handler)
			}(connection)
		}
	}()

	go func() {
		defer wg.Done()

		<-ctx.Done()
		s.listener.Close()
	}()

	wg.Wait()
	return nil
}

func (s *TCPServer) handleConnection(ctx context.Context, connection net.Conn, handler Handler) {
	defer func() {
		if v := recover(); v != nil {
			s.log.Error("captured panic", "panic", v)
		}

		if err := connection.Close(); err != nil {
			s.log.Warn("failed to close connection", l.ErrAttr(err))
		}
	}()

	request := make([]byte, s.bufferSize)

	for {
		if s.isClosing() {
			break
		}
		if s.idleTimeout != 0 {
			if err := connection.SetDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.log.Warn("failed to set read deadline", l.ErrAttr(err))
				break
			}
		}

		count, err := connection.Read(request)
		if err != nil && err != io.EOF {
			s.log.Warn("failed to read", l.ErrAttr(err))
			break
		} else if count == s.bufferSize {
			s.log.Warn("small buffer size")
			break
		}

		response := handler(ctx, request[:count])
		if _, err = connection.Write(response); err != nil {
			s.log.Warn("failed to write", l.ErrAttr(err))
			break
		}
	}
}

func (s *TCPServer) isClosing() bool {
	return s.closing.Load()
}

// Close graceful shutdown server
func (s *TCPServer) Close() {
	if !s.closing.CompareAndSwap(false, true) {
		return
	}

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			s.log.Warn("failed to close listener", l.ErrAttr(err))
		}
		s.log.Info("server shutdown competed")
	}

	return
}
