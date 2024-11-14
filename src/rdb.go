package src

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func StoreRDB(srv *Server) {
	ticker := time.NewTicker(time.Duration(1) * time.Hour)
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

func Store(srv *Server) {
	// Convert Map into binary little endian

	valueToSave := srv.Database.Data

	if len(valueToSave) == 0 {
		return
	}

	value, err := json.Marshal(valueToSave)

	if err != nil {
		log.Fatal(err)
		return
	}

	dir := *srv.Config["dir"].(*string)

	dirPath := filepath.Join(".", dir)

	err = os.MkdirAll(dirPath, os.ModePerm)

	if err != nil {
		log.Fatal(err)
		return
	}

	dbfilename := *srv.Config["dbfilename"].(*string)

	filePath := filepath.Join(".", dir, dbfilename)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer file.Close()

	bytesBuffer := new(bytes.Buffer)

	binary.Write(bytesBuffer, binary.LittleEndian, value)

	_, err = bytesBuffer.WriteTo(file)

	if err != nil {
		log.Fatal(err)
		return
	}
}

func Retreive(srv *Server) {
	dir := *srv.Config["dir"].(*string)

	dbfilename := *srv.Config["dbfilename"].(*string)

	filePath := filepath.Join(".", dir, dbfilename)

	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
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

	var data map[string]map[string]interface{}

	err = json.Unmarshal(buffer[:readLength], &data)

	if err != nil {
		log.Fatal(err)
	}

	// TODO: Update Expire Time Value

	srv.Database.Data = data
}
