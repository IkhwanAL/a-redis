package src

import (
	"context"
	"fmt"
	"log"
	"net"
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
	return nil
	// buf := make([]byte, 1024*4)
	// TODO: Create A Handshake Between Primary And Slave
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
