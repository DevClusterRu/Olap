package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	. "olap/engine"
	"olap/files"
	"os"
	"path/filepath"
)

func CsvHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {

		if req.FormValue("pool")=="" {
			http.Error(w, "No pool", 500)
			return
		}

		src, hdr, err := req.FormFile("filename")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer src.Close()

		//dst, err := os.Create(filepath.Join(os.TempDir(), hdr.Filename))
		fname:=filepath.Join("./", hdr.Filename)
		dst, err := os.Create(fname)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err!=nil{
			http.Error(w, err.Error(), 500)
			return
		} else {
			go files.FileParsing(fname, req.Form.Get("pool"))
			fmt.Fprintf(w, "OK")
		}

	}

}

func PoolCreateHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if req.Form.Get("poolname")!=""{
		Mongo.AddPool(req.Form.Get("poolname"))
	}
	fmt.Fprintf(w, "OK")
}

func PoolAggregateHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if req.Form.Get("idpool")!=""{
		packReady:=Mongo.PoolAggregation(req.Form.Get("idpool"))
		b,_:=json.Marshal(packReady)
		fmt.Fprintf(w, string(b))
	}

}

func PoolRemoveHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if req.Form.Get("idpool")!=""{
		Mongo.AddPool(req.Form.Get("idpool"))
	}
	fmt.Fprintf(w, "OK")
}


func main() {

	Init()
	http.HandleFunc("/csvPush", CsvHandler)
	http.HandleFunc("/poolCreate", PoolCreateHandler)
	http.HandleFunc("/poolRemove", PoolRemoveHandler)
	http.HandleFunc("/poolAggregate", PoolAggregateHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
