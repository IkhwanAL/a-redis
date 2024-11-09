package src

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var MaxBuffer int = 4096

type Server struct {
	Port int

	Database Database
	Config   map[string]interface{}
}

func loadServerCertAndKey() (string, string) {
	serverKeyPath := filepath.Join("..", "..", "..", "server.key")
	serverCert := filepath.Join("..", "..", "..", "server.crt")

	return serverCert, serverKeyPath
}

// Future Problem How can i tell client doing a request with / without tls
// Local Server Certificate Not Client Certificate
func (s *Server) Run(ctx context.Context) error {
	certPath, privateKeyPath := loadServerCertAndKey()

	cert, err := tls.LoadX509KeyPair(certPath, privateKeyPath)

	if err != nil {
		return fmt.Errorf("failed to Load Certificate: %s", err)
	}

	listener, err := tls.Listen("tcp", ":6379", &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	})

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
	case "config":
		resp = s.GetConfig(messages[2:])
	case "info":
		resp = ParseGenerateRESP("PONG")
	default:
		return fmt.Errorf("unknow command")
	}

	_, err := conn.Write([]byte(resp))

	return err
}

func (s *Server) Set(args []string) string {
	s.Database.Set(args[0], args[1])

	if len(args) > 2 && args[2] == "px" {
		ttl, err := strconv.Atoi(args[3])

		if err != nil {
			log.Fatal(err)
		}

		t := time.Now()

		t = t.Add(time.Duration(ttl) * time.Millisecond)

		s.Database.SetSetting(args[0], "EXPIRETM", t.Format(time.RFC3339))
		s.Database.SetSetting(args[0], "TTL", ttl)

		go func() {
			<-time.After(time.Duration(ttl) * time.Millisecond)
			s.Database.Unset(args[0])
		}()
	}

	return ParseGenerateRESP("OK")
}

func (s *Server) Get(args []string) string {
	value := s.Database.Get(args[0])

	return ParseGenerateRESP(value)
}

func (s *Server) GetConfig(args []string) string {
	key := args[0]

	value, ok := s.Config[key]

	if !ok {
		return ParseGenerateRESPError("-1")
	}

	return ParseGenerateRESP(*value.(*string))
}
