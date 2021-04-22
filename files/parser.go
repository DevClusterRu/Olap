package files

import (
	"fmt"
	"io"
	"olap/engine"
	"os"
	"strings"
)

func FileParsing(fname string, pool string)  {

	engine.Init()

	file, err := os.Open(fname)
	if err != nil{
		fmt.Println(err)
	}
	defer file.Close()

	data := make([]byte, 1)

	//var pointer int64 = 0
	firstRow := true
	var accum string = ""

	for{
		//file.Seek(pointer, 0)

		n, err := file.Read(data)
		if n == 1 {
			if string(data)!="\n" {
				accum += string(data)
			}
		} else {
			break
		}
		if err == io.EOF{   // если конец файла
			extractFields(accum, pool)
			accum = ""
			break           // выходим из цикла
		}
		if string(data)=="\n"{
			if (firstRow) {
				firstRow = false
				accum = ""
			} else {
				extractFields(accum, pool)
				accum = ""
			}
		}

	}

}

func extractFields(str string, pool string) bool {
	fields := strings.Split(str,",")
	if len(fields)<3{
		return false
	}
	engine.Mongo.InsertRecord(pool, fields[0],fields[1],fields[2])

	return true
}