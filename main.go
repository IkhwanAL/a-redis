package main

import (
	"context"
	"flag"
	"fmt"

	Srv "github.com/IkhwanAL/a-redis/src"
)

func main() {

	dir := flag.String("dir", "/tmp/redis-files", "the directory where RDB Store")
	dbfilename := flag.String("dbfilename", "redis.rdb", "the filename for RDB")
	replicaOf := flag.String("replicaof", "master", "set replica")
	port := flag.Int("port", 6379, "port to connect")

	flag.Parse()

	isAReplica := false

	fmt.Printf("Run Server %d\n", *port)

	replica := Srv.NewReplication("8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb", *replicaOf, 0)
	if replica.Role != Srv.MASTER {
		replica.Run(context.Background())
	}
	fmt.Print(replica.Role)

	s := Srv.Server{
		Port: *port,
		Database: Srv.Database{
			Data: make(map[string]map[string]interface{}),
		},
		Config: map[string]interface{}{
			"dir":        dir,
			"dbfilename": dbfilename,
		},
	}

	if !isAReplica {
		Srv.Retreive(&s)

		Srv.StoreRDB(&s)
	}

	s.Run(context.Background())
}
