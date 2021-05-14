package main

import (
	"fmt"
	"log"
	"net/http"
	. "olap/engine"
	"olap/route"
)

func main() {

	fmt.Println("Start v1.2")

	Init()
	http.HandleFunc("/", route.EmptyHandler)
	http.HandleFunc("/api/csvPush", route.CsvHandler)
	http.HandleFunc("/api/poolCreate", route.PoolCreateHandler)
	http.HandleFunc("/api/poolList", route.PoolListHandler)
	http.HandleFunc("/api/poolRemove", route.PoolRemoveHandler)
	http.HandleFunc("/api/poolAggregate", route.PoolAggregateHandler)
	http.HandleFunc("/api/chartByYear", route.ChartAggregateByYearHandler)
	http.HandleFunc("/api/chartByMonth", route.ChartAggregateByMonthHandler)
	http.HandleFunc("/api/chartByDay", route.ChartAggregateByDayHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
