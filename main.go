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
	"regexp"
	"strconv"
)

func checkMethod(w http.ResponseWriter, req *http.Request) bool {
	if req.Method=="POST"{
		return true
	} else{
		fmt.Fprintf(w, "Sorry, only POST method supported.")
	}
	return false
}

func CsvHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {

		if req.FormValue("pool") == "" {
			http.Error(w, "No pool", 500)
			return
		}

		src, hdr, err := req.FormFile("filename")
		if err != nil {
			http.Error(w, "filename!!"+err.Error(), 500)
			return
		}
		defer src.Close()

		//dst, err := os.Create(filepath.Join(os.TempDir(), hdr.Filename))
		fname := filepath.Join("./", hdr.Filename)
		dst, err := os.Create(fname)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		} else {
			go files.FileParsing(fname, req.Form.Get("pool"))
			fmt.Fprintf(w, "OK")
		}

	}

}

func PoolListHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-Type", "application/json")
	fmt.Fprintf(w, Mongo.PoolList())
}

func PoolCreateHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	poolId := "ERROR"
	if req.Form.Get("name") != "" && req.Form.Get("description") != "" {
		poolId = Mongo.AddPool(req.Form.Get("name"), req.Form.Get("description"))
	}
	fmt.Fprintf(w, poolId)
}

func PoolAggregateHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if req.Form.Get("poolId") != "" {
		aggregation := DoGraph(Mongo.GraphPoolAggregation(req.URL.Query()))
		w.Header().Set("content-Type", "application/json")
		fmt.Fprintf(w, aggregation)
	}
}

func PoolRemoveHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if req.Form.Get("poolId") != "" {
		res, err := Mongo.RemovePool(req.Form.Get("poolId"))
		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			fmt.Fprintf(w, strconv.Itoa(int(res.DeletedCount)))
		}
	} else {
		fmt.Fprintf(w, "ERROR")
	}

}

func LinkCreateHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if req.Form.Get("poolId") != "" {
		if Mongo.CreateLink(req.Form.Get("poolId"), req.Form.Get("from"), req.Form.Get("to")) == true {
			fmt.Fprintf(w, "OK")
		} else {
			fmt.Fprintf(w, "ERROR")
		}
	} else {
		fmt.Fprintf(w, "Dont see poolId")
	}
}

func LinkRemoveHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if req.Form.Get("linkId") != "" {
		res, err := Mongo.RemoveLink(req.Form.Get("linkId"))
		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			fmt.Fprintf(w, strconv.Itoa(int(res.DeletedCount)))
		}
	} else {
		fmt.Fprintf(w, "ERROR")
	}
}

func StreamHandler(w http.ResponseWriter, req *http.Request) {
	if checkMethod(w, req) {
		var p Stream
		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		Mongo.InsertRecord(p.Pool,p.Object,p.Event, p.Timestamp)

		fmt.Fprintf(w, "OK")
	}
}

func EmptyHandler(w http.ResponseWriter, req *http.Request) {

	var validID = regexp.MustCompile(`\/.*$`)
	fmt.Println(validID.FindString(req.URL.Path))

	fmt.Fprintf(w, "U on root")
}

func main() {

	fmt.Println("Start v1.2")

	Init()
	http.HandleFunc("/", EmptyHandler)
	http.HandleFunc("/api/csvPush", CsvHandler)
	http.HandleFunc("/api/linkCreate", LinkCreateHandler)
	http.HandleFunc("/api/linkRemove", LinkRemoveHandler)
	http.HandleFunc("/api/poolCreate", PoolCreateHandler)
	http.HandleFunc("/api/poolList", PoolListHandler)
	http.HandleFunc("/api/poolRemove", PoolRemoveHandler)
	http.HandleFunc("/api/poolAggregate", PoolAggregateHandler)
	http.HandleFunc("/api/stream", StreamHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
