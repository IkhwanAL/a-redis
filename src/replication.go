package src

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	MASTER = "master"
	SLAVE  = "slave"
)

type Replication struct {
	Host      string
	Port      string
	Role      string
	ReplicaId string
	offset    int
}

func (r *Replication) Run(ctx context.Context) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", r.Host, r.Port))

	if err != nil {
		log.Fatal(err)
	}

	r.handleConnection(conn)

	return nil
}

func (r *Replication) handleConnection(conn net.Conn) error {
	defer conn.Close()

	r.Handshake(conn)

	buf := make([]byte, 1024*3)

	_, err := conn.Read(buf)

	if err == io.EOF {
		return nil
	}

	return nil
}

func (r *Replication) Handshake(conn net.Conn) {
	r.handshakePing(conn)
	r.handshakeREPLCONE(conn)
	r.handshakePSYNC(conn)
}

func (r *Replication) handshakePing(conn net.Conn) {
	resp := ParseGenerateArrayValueRESP("PING")
	conn.Write([]byte(resp))
	buf := make([]byte, 10)
	_, err := conn.Read(buf)

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	value := ParseReadRESP(buf)

	if value[0] != "PONG" {
		log.Fatal("incorrect response")
	}
}

func (r *Replication) handshakeREPLCONE(conn net.Conn) {
	request := []string{
		"REPLCONE",
		"listening-port",
		r.Port,
	}

	generatedRequest := ParseGenerateArrayValueRESP(request...)
	conn.Write([]byte(generatedRequest))
	buf := make([]byte, 15)
	_, err := conn.Read(buf)

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	value := ParseReadRESP(buf)

	if value[0] != "OK" {
		log.Fatal("incorrect response")
	}

	request = []string{
		"REPLCONE",
		"capa",
		"psync2",
	}

	generatedRequest = ParseGenerateArrayValueRESP(request...)
	conn.Write([]byte(generatedRequest))

	_, err = conn.Read(buf)

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	value = ParseReadRESP(buf)

	if value[0] != "OK" {
		log.Fatal("incorrect response")
	}
}

func (r *Replication) handshakePSYNC(conn net.Conn) {
	request := []string{
		"PSYNC",
		"?",
		"-1",
	}

	generatedRequest := ParseGenerateArrayValueRESP(request...)

	conn.Write([]byte(generatedRequest))

	buf := make([]byte, 56)

	_, err := conn.Read(buf)

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	value := ParseReadRESP(buf)

	sync := splitString(value[0])

	if sync[0] != "FULLRESYNC" {
		log.Fatal("incorrect response")
	}

	offset, err := strconv.Atoi(sync[2])

	if err != nil {
		log.Fatal("failed to convert string into int in handshakePYSNC")
	}

	r.ReplicaId = sync[1]
	r.offset = offset
}

func splitString(value string) []string {
	return strings.Split(value, " ")
}

func NewReplication(replicaId string, address string, offset int) Replication {
	splitedAddress := strings.Split(address, ":")

	if splitedAddress[0] == MASTER {
		return Replication{
			ReplicaId: replicaId,
			offset:    offset,
			Role:      MASTER,
		}
	}

	return Replication{
		Host:      splitedAddress[0],
		Port:      splitedAddress[1],
		Role:      SLAVE,
		ReplicaId: replicaId,
		offset:    offset,
	}
}
