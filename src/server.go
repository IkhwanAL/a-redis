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
	Port          int
	HashTableInfo map[string]int

	Database Database
	Config   map[string]interface{}
}

func (s *Server) Run(ctx context.Context, replica *Replication) error {

	if s.Port == 6380 {
		return s.runServerSecureConnection(ctx, replica)
	}

	return s.runUnSecureConnection(ctx, replica)
}

func (s *Server) runUnSecureConnection(ctx context.Context, replica *Replication) error {
	port := fmt.Sprintf(":%d", s.Port)

	listener, err := net.Listen("tcp", port)

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

			go s.handleConnection(conn, replica)
		}
	}
}

func loadServerCertAndKey() (string, string) {
	// Need to accept cli option to read server and server certificate
	serverKeyPath := filepath.Join("..", "..", "..", "server.key")
	serverCert := filepath.Join("..", "..", "..", "server.crt")

	return serverCert, serverKeyPath
}

func (s *Server) runServerSecureConnection(ctx context.Context, replica *Replication) error {
	certPath, privateKeyPath := loadServerCertAndKey()

	cert, err := tls.LoadX509KeyPair(certPath, privateKeyPath)

	if err != nil {
		return fmt.Errorf("failed to Load Certificate: %s", err)
	}

	port := fmt.Sprintf(":%d", s.Port)

	listener, err := tls.Listen("tcp", port, &tls.Config{
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

			go s.handleConnection(conn, replica)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn, replica *Replication) error {
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

		err = s.runMessage(conn, buf, replica)

		if err != nil {
			return fmt.Errorf("error running a command: %w", err)
		}

	}
}

func (s *Server) runMessage(conn net.Conn, requests []byte, replica *Replication) error {
	messages := ParseReadRESP(requests)

	cmd := strings.ToLower(messages[0])

	var resp string

	switch cmd {
	case "ping":
		resp = ParseGenerateRESP("PONG")
	case "replcone":
		resp = ParseGenerateRESP("+OK")
	case "set":
		resp = s.Set(messages[1:])
	case "get":
		resp = s.Get(messages[1:])
	case "config":
		resp = s.GetConfig(messages[2:])
	case "info":
		role := fmt.Sprintf("role:%s", replica.Role)
		replicaId := fmt.Sprintf("master_replid:%s", replica.ReplicaId)
		offset := fmt.Sprintf("master_repl_offset:%d", replica.offset)

		resp = ParseGenerateMultipleValue(role, replicaId, offset)
	case "psync":
		replicaId := RandomReplciateSeedId()
		offset := 0

		resp = fmt.Sprintf("+FULLRESYNC %s %s\r\n", replicaId, strconv.Itoa(offset))
	default:
		return fmt.Errorf("unknow command")
	}

	_, err := conn.Write([]byte(resp))

	return err
}

func (s *Server) Set(args []string) string {
	s.Database.Set(args[0], args[1])
	s.HashTableInfo["keyValues"] += 1

	if len(args) > 2 && args[2] == "px" {
		ttl, err := strconv.Atoi(args[3])

		if err != nil {
			log.Fatal(err)
		}

		t := time.Now()
		s.HashTableInfo["withPx"] += 1

		t = t.Add(time.Duration(ttl) * time.Millisecond)

		s.Database.SetSetting(args[0], "TTL", t.UnixMilli())

		go func() {
			<-time.After(time.Duration(ttl) * time.Millisecond)
			s.Database.Unset(args[0])
		}()
	}

	return ParseGenerateRESP("OK")
}

func (s *Server) Get(args []string) string {
	if s.Config["isAReplica"].(bool) {
		return ParseGenerateRESPError("-1")
	}

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
