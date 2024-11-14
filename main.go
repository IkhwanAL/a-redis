package main

import (
	"context"
	"flag"
	"fmt"

	Srv "github.com/IkhwanAL/a-redis/src"
)

func main() {

	dir := flag.String("dir", "", "the directory where RDB Store")
	dbfilename := flag.String("dbfilename", "", "the filename for RDB Store")
	replicaOf := flag.String("replicaof", "", "set replica")
	port := flag.Int("port", 6379, "port to connect")

	flag.Parse()

	isAReplica := false

	fmt.Printf("Run Server %d\n", *port)

	if *replicaOf != "" {
		replica := Srv.NewReplication(*replicaOf)

		replica.Run(context.Background())
	}

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
