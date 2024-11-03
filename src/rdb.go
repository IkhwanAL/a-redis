package src

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func StoreRDB(srv *Server) {
	ticker := time.NewTicker(time.Duration(4) * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
			case t := <-ticker.C:
				fmt.Println("Tick at:", t)
				Store(srv)
			}
		}
	}()
}

func Retreive(srv *Server) {
	dir := *srv.Config["dir"].(*string)

	dbfilename := *srv.Config["dbfilename"].(*string)

	filePath := filepath.Join(".", dir, dbfilename)

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0777)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer file.Close()

	buffer := make([]byte, 1024*4)

	readLength, err := file.Read(buffer)

	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
			return
		}
	}

	var data map[string]interface{}

	err = json.Unmarshal(buffer[:readLength], &data)

	if err != nil {
		log.Fatal(err)
	}

	srv.Database.Data = data
}

func Store(srv *Server) {
	// Convert Map into binary little endian

	test := map[string]interface{}{
		"HSR":     "asd",
		"HJS":     "DS",
		"HKA:YUI": "TYR",
	}

	value, err := json.Marshal(test)

	if err != nil {
		log.Fatal(err)
		// Panic / Restart
		return
	}

	dir := *srv.Config["dir"].(*string)

	dirPath := filepath.Join(".", dir)

	err = os.MkdirAll(dirPath, os.ModePerm)

	if err != nil {
		log.Fatal(err)
		return
	}

	// Open File And Store The Value

	dbfilename := *srv.Config["dbfilename"].(*string)

	filePath := filepath.Join(".", dir, dbfilename)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer file.Close()

	_, err = file.Write(value)

	if err != nil {
		log.Fatal(err)
		return
	}
}

//
