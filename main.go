package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	Srv "github.com/IkhwanAL/a-redis/src"
)

func main() {

	dir := flag.String("dir", "/tmp/redis-files", "the directory where RDB Store")
	dbfilename := flag.String("dbfilename", "redis.rdb", "the filename for RDB")
	replicaOf := flag.String("replicaof", "master", "set replica")
	port := flag.Int("port", 6379, "port to connect")

	flag.Parse()

	i, _ := strconv.ParseInt("007", 10, 32)

	fmt.Printf("%v\n", i)

	isAReplica := false

	randomSeed := Srv.RandomReplciateSeedId()

	fmt.Printf("Run Server %d\n", *port)

	replica := Srv.NewReplication(randomSeed, *replicaOf, 0)

	if replica.Role != Srv.MASTER {
		replica.Run(context.Background())
	}

	s := Srv.Server{
		Port:          *port,
		HashTableInfo: make(map[string]int),
		Database: Srv.Database{
			Data: make(map[string]map[string]interface{}),
		},
		Config: map[string]interface{}{
			"dir":        dir,
			"dbfilename": dbfilename,
		},
	}

	if !isAReplica {
		// Srv.Retreive(&s)

		Srv.StoreRDBFormat(&s)
	}

	// s.Run(context.Background(), &replica)
}
