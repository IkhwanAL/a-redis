package src

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	SLAVE = "SLAVE"
)

type Replication struct {
	Host string
	Port string
	Role string
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

	buf := make([]byte, 1024*4)
	// TODO: Create A Handshake Between Primary And Slave
	for {
		conn.Read(buf)
	}
}

func NewReplication(address string) Replication {
	splitedAddress := strings.Split(address, ":")

	return Replication{
		Host: splitedAddress[0],
		Port: splitedAddress[1],
		Role: SLAVE,
	}
}
