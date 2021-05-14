package files

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	. "olap/engine"
	"os"
	"strconv"
	"strings"
)

func FileParsing(fname string, pool string) {

	Mongo.PoolStatusChange(pool, "busy")

	file, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data := make([]byte, 1)

	//var pointer int64 = 0
	firstRow := true
	var accum string = ""
	var ref []string
	for {
		//file.Seek(pointer, 0)

		n, err := file.Read(data)
		if n == 1 {
			if string(data) != "\n" {
				accum += string(data)
			}
		} else {
			fmt.Println("No bytes read")
			break
		}
		if err == io.EOF { // если конец файла
			if ref != nil {
				fmt.Println("Ref==nil")
				extractFields(accum, pool, ref)
				accum = ""
			}
			break // выходим из цикла
		}
		if string(data) == "\n" {
			if firstRow {
				firstRow = false
				ref = detectingColumns(accum)
				if ref == nil {
					fmt.Println("Ref error")
					break
				}
				accum = ""
			} else {
				extractFields(accum, pool, ref)
				accum = ""
			}
		}

	}

	Mongo.PoolStatusChange(pool, "active")

}

func extractFields(str string, poolId string, ref []string) bool {
	fields := strings.Split(str, ",")
	if len(fields) < 3 {
		fmt.Println("Low")
		return false
	}
	document := make(bson.M)

	paramNumber := 1

	for k, v := range fields {
		switch ref[k] {
		case "object":
			document["object"] = v
			break
		case "numeric":
			n, err := strconv.ParseInt(v, 10, 64)
			if err!=nil{
				n = 0
			}
			document["numeric"] = n
			break
		case "event":
			document["event"] = v
			break
		case "timestamp":
			document["timestamp"] = StrToDate(v)
			break
		default:
			document["param"+strconv.Itoa(paramNumber)] = v
			paramNumber++
		}
	}

	Mongo.InsertRecord(poolId, document)
	return true
}

func detectingColumns(str string) []string {
	columns := strings.Split(str, ",")

	var ref = make([]string, len(columns))

	allPresent := 0

	for k, v := range columns {
		if strings.Contains(strings.ToUpper(v), "NUMERIC") {
			ref[k] = "numeric"
			allPresent++
		}
		if strings.Contains(strings.ToUpper(v), "EVENT") {
			ref[k] = "event"
			allPresent++
		}
		if strings.Contains(strings.ToUpper(v), "DATE") || strings.Contains(strings.ToUpper(v), "TIME") {
			ref[k] = "timestamp"
			allPresent++
		}
		if strings.Contains(strings.ToUpper(v), "CASE") || strings.Contains(strings.ToUpper(v), "OBJECT") {
			ref[k] = "object"
			allPresent++
		}

	}
	if allPresent < 3 {
		return nil
	}
	return ref
}
