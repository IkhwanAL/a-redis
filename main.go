package main

import (
	"context"
	"flag"
	"fmt"

	Srv "github.com/IkhwanAL/a-redis/src"
)

func main() {
	fmt.Print("Run Server 6379\n")

	dir := flag.String("dir", "", "the directory where RDB Store")
	dbfilename := flag.String("dbfilename", "", "the filename for RDB Store")

	flag.Parse()

	s := Srv.Server{
		Port: 6379,
		Database: Srv.Database{
			Data: make(map[string]map[string]interface{}),
		},
		Config: map[string]interface{}{
			"dir":        dir,
			"dbfilename": dbfilename,
		},
	}

	Srv.Retreive(&s)

	go Srv.StoreRDB(&s)

	s.Run(context.Background())
}
