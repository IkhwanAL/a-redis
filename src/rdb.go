package src

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc64"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var REDIS_MAGIC_NUMBER []byte = []byte{0x52, 0x45, 0x44, 0x49, 0x53} // Redis
var REDIS_VERSION_NUMBER []byte = []byte{0x30, 0x30, 0x31, 0x31}     // 0011 -> 3

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

func StoreRDBFormat(srv *Server) {
	srv.Database.Data = map[string]map[string]interface{}{
		"HSD": {
			"VALUE": "HSD",
		},
		"IIO": {
			"VALUE": "PPO",
		},
		"JKL": {
			"VALUE": "UIO",
		},
	}

	dir := *srv.Config["dir"].(*string)

	dirPath := filepath.Join(".", dir)

	err := os.MkdirAll(dirPath, os.ModePerm)

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

	var buffer strings.Builder

	// -- Header Section
	buffer.WriteString(hex.EncodeToString(REDIS_MAGIC_NUMBER))
	buffer.WriteString(hex.EncodeToString(REDIS_VERSION_NUMBER))

	// -- Metadata Section (We Can Add more in the future)
	buffer.WriteByte(byte(START_METADATA))
	buffer.WriteString(hex.EncodeToString([]byte("redis-ver")))
	buffer.WriteString(hex.EncodeToString([]byte("6.0.16")))

	// -- Database Section Where Data is Exists
	buffer.WriteByte(byte(START_DB_SECTION))
	buffer.WriteByte(byte(00))

	buffer.WriteByte(byte(START_HASHTABEL_INFO))
	buffer.WriteByte(byte(srv.HashTableInfo["keyValue"]))
	buffer.WriteByte(byte(srv.HashTableInfo["withPx"]))

	buffer.WriteByte(byte(STRING_DATATYPE_FLAG))

	var keys string
	var values string

	for key, value := range srv.Database.Data {
		keys += key

		v := value["VALUE"].(string)
		values += v
	}

	buffer.WriteString(hex.EncodeToString([]byte(keys)))
	buffer.WriteString(hex.EncodeToString([]byte(values)))

	// Skip The Key Expire Thing First

	for key, value := range srv.Database.Data {
		buffer.WriteByte(byte(STRING_DATATYPE_FLAG))

		// No Validation because value is always exists
		// if not there's a bug in store mechanism
		v := value["VALUE"].(string)

		buffer.WriteString(hex.EncodeToString([]byte(key)))
		buffer.WriteString(hex.EncodeToString([]byte(v)))

	}
	buffer.WriteByte(byte(EOF))

	table := crc64.MakeTable(crc64.ECMA)

	checksum := crc64.Checksum([]byte(buffer.String()), table)

	chs := make([]byte, 8)
	binary.LittleEndian.PutUint64(chs, checksum)

	buffer.WriteString(hex.EncodeToString(chs))

	bytesBuffer := new(bytes.Buffer)

	binary.Write(bytesBuffer, binary.LittleEndian, []byte(buffer.String()))

	_, err = bytesBuffer.WriteTo(file)

	if err != nil {
		log.Fatal(err)
		return
	}
}

func Store(srv *Server) {
	// Convert Map into binary little endian

	// valueToSave := srv.Database.Data
	valueToSave := map[string]map[string]interface{}{
		"ASD":  {"VALUE": "asd"},
		"HASD": {"VALUE": "hasd"},
		"JJK":  {"VALUE": "jjk"},
	}

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
		if err == io.EOF {
			fmt.Printf("end of file")
			return
		} else {
			log.Fatal(err)
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
