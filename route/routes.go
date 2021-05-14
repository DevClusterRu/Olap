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
	"strconv"
	"strings"
	"time"
)

func getSegments(s string) []string {
	var response []string
	res := strings.Split(s, "/")
	for _, v := range res {
		if v != "" {
			response = append(response, v)
		}
	}
	return response
}

func checkMethod(w http.ResponseWriter, req *http.Request) bool {
	if req.Method == "POST" {
		return true
	} else {
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



func ChartAggregateByYearHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	result := engine.Mongo.CubeAggregator(req.Form.Get("poolId"),0,0)
	b, _ := json.Marshal(result)
	w.Header().Set("content-Type", "application/json")
	fmt.Fprintf(w, string(b))
}

func ChartAggregateByMonthHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	y, e := strconv.Atoi(req.Form.Get("year"))
	if e != nil {
		fmt.Fprintf(w, "Year error")
		return
	}
	result := engine.Mongo.CubeAggregator(req.Form.Get("poolId"), y,0)
	b, _ := json.Marshal(result)
	w.Header().Set("content-Type", "application/json")
	fmt.Fprintf(w, string(b))
}

func ChartAggregateByDayHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	y, e := strconv.Atoi(req.Form.Get("year"))
	if e != nil {
		fmt.Fprintf(w, "Year error")
		return
	}
	m, e := strconv.Atoi(req.Form.Get("month"))
	if e != nil {
		fmt.Fprintf(w, "Month error")
		return
	}
	result := engine.Mongo.CubeAggregator(req.Form.Get("poolId"), y, m)
	b, _ := json.Marshal(result)
	w.Header().Set("content-Type", "application/json")
	fmt.Fprintf(w, string(b))
}

func stream(w http.ResponseWriter, req *http.Request, poolId string) {
	if checkMethod(w, req) {
		var p bson.M
		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		poolId, _ := primitive.ObjectIDFromHex(poolId)
		p["pool_id"] = poolId

		p["timestamp"] = engine.StrToDate(p["timestamp"].(string))

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		engine.Mongo.Client.Database(engine.Mongo.DB).Collection("test").InsertOne(ctx, p)
		fmt.Fprintf(w, "OK")
	}
}

func EmptyHandler(w http.ResponseWriter, req *http.Request) {

	seg := getSegments(req.URL.Path)
	if len(getSegments(req.URL.Path)) == 3 && seg[0] == "api" && seg[1] == "stream" {
		stream(w, req, seg[2])
	}

	fmt.Fprintf(w, "U on root")
}
