package src

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
)

var MaxBuffer int = 4096

type Server struct {
	Port int

	Database Database
	Config   map[string]interface{}
}

func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		return fmt.Errorf("failed to bind with port: %d", s.Port)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			conn, err := listener.Accept()

			if err != nil {
				return fmt.Errorf("error accepting connection: %w", err)
			}

			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) error {
	defer conn.Close()

	buf := make([]byte, MaxBuffer)

	for {
		_, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				return nil
			}

			fmt.Printf("Error Reading from Client: %s", err)
			return nil
		}

		err = s.runMessage(conn, buf)

		if err != nil {
			return fmt.Errorf("error running a command: %w", err)
		}

	}
}

func (s *Server) runMessage(conn net.Conn, requests []byte) error {
	messages := ParseReadRESP(requests)

	cmd := strings.ToLower(messages[0])

	var resp string

	switch cmd {
	case "ping":
		resp = ParseGenerateRESP("PONG")
	case "set":
		resp = s.Set(messages[1:])
	case "get":
		resp = s.Get(messages[1:])
	default:
		return fmt.Errorf("unknow command")
	}

	_, err := conn.Write([]byte(resp))

	return err
}

func (s *Server) Set(args []string) string {
	s.Database.Set(args[0], args[1])

	return ParseGenerateRESP("OK")
}

func (s *Server) Get(args []string) string {
	value := s.Database.Get(args[0])

	return ParseGenerateRESP(value)
}
