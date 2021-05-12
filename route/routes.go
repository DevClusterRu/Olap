package route

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"olap/engine"
	"olap/files"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
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

		if req.FormValue("poolId") == "" {
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
			go files.FileParsing(fname, req.FormValue("poolId"))
			fmt.Fprintf(w, "OK")
		}

	}

}

func PoolListHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-Type", "application/json")
	fmt.Fprintf(w, engine.Mongo.PoolList())
}

func PoolCreateHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	poolId := "ERROR"
	if req.Form.Get("name") != "" && req.Form.Get("description") != "" {
		poolId = engine.Mongo.AddPool(req.Form.Get("name"), req.Form.Get("description"))
	}
	fmt.Fprintf(w, poolId)
}

func PoolAggregateHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if req.Form.Get("poolId") != "" {
		aggregation := engine.DoGraph(engine.Mongo.GraphPoolAggregation(req.URL.Query()))
		w.Header().Set("content-Type", "application/json")
		fmt.Fprintf(w, aggregation)
	}
}

func PoolRemoveHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if req.Form.Get("poolId") != "" {
		res, err := engine.Mongo.RemovePool(req.Form.Get("poolId"))
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
		if engine.Mongo.CreateLink(req.Form.Get("poolId"), req.Form.Get("from"), req.Form.Get("to")) == true {
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
		res, err := engine.Mongo.RemoveLink(req.Form.Get("linkId"))
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
		var p bson.M
		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		poolId, _ := primitive.ObjectIDFromHex(p["poolId"].(string))
		p["poolId"] = poolId
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		engine.Mongo.Client.Database(engine.Mongo.DB).Collection("test").InsertOne(ctx, p)
		fmt.Fprintf(w, "OK")
	}
}

func EmptyHandler(w http.ResponseWriter, req *http.Request) {

	var validID = regexp.MustCompile(`\/.*$`)
	fmt.Println(validID.FindString(req.URL.Path))
	fmt.Fprintf(w, "U on root")
}