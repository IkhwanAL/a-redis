package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"time"

	Srv "github.com/IkhwanAL/a-redis/src"
)

func randomReplciateSeedId() string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomId := make([]byte, 40)

	for i := 0; i < 40; i++ {
		randomId[i] = charset[random.Intn(len(charset))]
	}

	return string(randomId)
}

func main() {

	dir := flag.String("dir", "/tmp/redis-files", "the directory where RDB Store")
	dbfilename := flag.String("dbfilename", "redis.rdb", "the filename for RDB")
	replicaOf := flag.String("replicaof", "master", "set replica")
	port := flag.Int("port", 6379, "port to connect")

	flag.Parse()

	isAReplica := false

	randomSeed := randomReplciateSeedId()

	fmt.Printf("Run Server %d\n", *port)

	replica := Srv.NewReplication(randomSeed, *replicaOf, 0)

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

	s.Run(context.Background(), &replica)
}
