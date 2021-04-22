package main

import (
	"fmt"
	"olap/engine"
	"olap/files"
	"testing"
)




//func TestGen(t *testing.T)  {
//	var fld [12]string
//	var tme [12]time.Time
//
//	fld[0] = "Device on"
//	fld[1] = "Device preparing for work"
//	fld[2] = "Preparing complete"
//	fld[3] = "Starting work"
//	fld[4] = "Work complete"
//	fld[5] = "Starting work"
//	fld[6] = "Work complete"
//	fld[7] = "Starting work"
//	fld[8] = "Work complete"
//	fld[9] = "Starting work"
//	fld[10] = "Work complete"
//	fld[11] = "Device off"
//
//	Init()
//
//	var i int64
//
//	for i=0; i<100000; i++{
//	tme[0] = time.Now().Add(time.Minute*time.Duration(i+1))
//	tme[1] = time.Now().Add(time.Minute*time.Duration(i+2))
//	tme[2] = time.Now().Add(time.Minute*time.Duration(i+3))
//	tme[3] = time.Now().Add(time.Minute*time.Duration(i+4))
//	tme[4] = time.Now().Add(time.Minute*time.Duration(i+5))
//	tme[5] = time.Now().Add(time.Minute*time.Duration(i+6))
//	tme[6] = time.Now().Add(time.Minute*time.Duration(i+7))
//	tme[7] = time.Now().Add(time.Minute*time.Duration(i+8))
//	tme[8] = time.Now().Add(time.Minute*time.Duration(i+9))
//	tme[9] = time.Now().Add(time.Minute*time.Duration(i+10))
//	tme[10] = time.Now().Add(time.Minute*time.Duration(i+11))
//	tme[11] = time.Now().Add(time.Minute*time.Duration(i+12))
//
//		Mongo.InsertRecord("SampleObject", fld[0],tme[0])
//		Mongo.InsertRecord("SampleObject", fld[1],tme[1])
//		Mongo.InsertRecord("SampleObject", fld[2],tme[2])
//		Mongo.InsertRecord("SampleObject", fld[3],tme[3])
//		Mongo.InsertRecord("SampleObject", fld[4],tme[4])
//		Mongo.InsertRecord("SampleObject", fld[5],tme[5])
//		Mongo.InsertRecord("SampleObject", fld[6],tme[6])
//		Mongo.InsertRecord("SampleObject", fld[7],tme[7])
//		Mongo.InsertRecord("SampleObject", fld[8],tme[8])
//		Mongo.InsertRecord("SampleObject", fld[9],tme[9])
//		Mongo.InsertRecord("SampleObject", fld[10],tme[10])
//		Mongo.InsertRecord("SampleObject", fld[11],tme[11])
//
//	}
//
//
//
//
//}
//


func TestAggregate(t *testing.T)  {
	engine.Init()
	x:=engine.Mongo.PoolAggregation("607fe9aadf1560b2a03d776d")
	fmt.Println(x)
}

func TestCsvParse(t *testing.T) {
	files.FileParsing("example.csv", "607fe9aadf1560b2a03d776d")
}

func TestCreatePool(t *testing.T)  {
	engine.Init()
	engine.Mongo.AddPool("TestPool")
}

func TestRemovePool(t *testing.T)  {
	engine.Init()
	engine.Mongo.RemovePool("607fe42e80ffbdd30b298c14")
}

func TestStructTest(t *testing.T)  {



}