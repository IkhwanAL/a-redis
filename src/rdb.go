package src

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc64"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
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

const MAX_UINT_6BITS = 1<<6 - 1
const MAX_UINT_14BITS = 1<<14 - 1

const MaskWith14BitsLengthHeader = 0b01 << 6

func sizeBitMask(size int) []byte {
	var buf []byte

	if size < MAX_UINT_6BITS {
		buf = make([]byte, 1)
		buf[0] = byte(size)
	}

	if size >= MAX_UINT_6BITS && size < MAX_UINT_14BITS {
		buf = make([]byte, 2)
		buf[0] = byte(size>>8 | MaskWith14BitsLengthHeader)
		buf[1] = byte(size)
	}

	if size >= MAX_UINT_14BITS && size < math.MaxInt32 {
		buf = make([]byte, 5)
		buf[0] = 0b10 << 6
		binary.BigEndian.PutUint32(buf[1:], uint32(size))
	}

	return buf
}

func tryCompressIntegerOfString(s string) ([]byte, bool) {
	value, err := IsStringCanBeEncodedAsUInteger(s)

	if err != nil {
		return nil, false
	}

	var buf []byte

	// Start Special Encoding For String

	if value > math.MinInt8 && value <= math.MaxInt8 {
		encodedLength := byte(0b11 << 6)
		encodedValue := byte(value)
		buf = []byte{encodedLength, encodedValue}
	}

	// if value > math.MinInt16

	return buf, true
}

// func compressString(s string) []byte {
// 	tryCompressIntegerOfString()
// }

func binaryWriteLengthEncoding(dst *bytes.Buffer, size int) {
	val := sizeBitMask(size)
	dst.Write(val)
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

	var buffer bytes.Buffer

	// -- Header Section
	buffer.Write([]byte(REDIS_MAGIC_NUMBER))
	buffer.Write([]byte(REDIS_VERSION_NUMBER))

	// -- Metadata Section (We Can Add more in the future)
	buffer.WriteByte(byte(START_METADATA))

	binaryWriteLengthEncoding(&buffer, len("redis-ver"))
	buffer.Write([]byte("redis-ver"))

	binaryWriteLengthEncoding(&buffer, len("6.0.16"))
	buffer.Write([]byte("6.0.16"))

	// -- Database Section Where Data is Exists
	buffer.WriteByte(byte(START_DB_SECTION))
	binaryWriteLengthEncoding(&buffer, 0)

	buffer.WriteByte(byte(START_HASHTABEL_INFO))
	binaryWriteLengthEncoding(&buffer, srv.HashTableInfo["keyValue"])
	binaryWriteLengthEncoding(&buffer, srv.HashTableInfo["withPx"])

	// Skip The Key Expire Thing First

	for key, value := range srv.Database.Data {
		buffer.WriteByte(byte(STRING_DATATYPE_FLAG))

		ttl, ok := value["TTL"].(int64)

		if ok {
			buffer.WriteByte(byte(EXPIRETIMEMS))
			timestampByte := make([]byte, 8)

			binary.LittleEndian.PutUint64(timestampByte, uint64(ttl))

			buffer.Write(timestampByte)
		}

		// No Validation because value is always exists
		// if not there's a bug in store mechanism
		v := value["VALUE"].(string)

		// Key
		binaryWriteLengthEncoding(&buffer, len(key))
		buffer.Write([]byte(key))

		binaryWriteLengthEncoding(&buffer, len(v))
		buffer.Write([]byte(v))

	}
	buffer.WriteByte(byte(EOF))

	table := crc64.MakeTable(crc64.ECMA)

	checksum := crc64.Checksum(buffer.Bytes(), table)

	chs := make([]byte, 8)
	binary.LittleEndian.PutUint64(chs, checksum)

	buffer.Write(chs)

	bytesBuffer := new(bytes.Buffer)

	binary.Write(bytesBuffer, binary.LittleEndian, buffer.Bytes())

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
