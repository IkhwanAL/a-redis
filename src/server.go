package src

import (
	"context"
	"fmt"
	"io"
	"net"
)

var MaxBuffer int = 4096

type Server struct {
	Port int

	Data   Database
	Config map[string]interface{}
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
				fmt.Printf("%s", "end of file")
				return nil
			}

			fmt.Printf("Error Reading from Client: %s", err)
			return nil
		}

		// Read

		conn.Write([]byte("*2\r\nOK\r\n"))
	}
}
