package engine

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func (m *MongoStructure) PoolAggregation(pool string) Package {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	poolId,_:=primitive.ObjectIDFromHex(pool)
	fmt.Println(poolId)
	matchStage := bson.D{{"$match", bson.M{"pool_id": poolId}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$event"}, {"moment",bson.M{"$first":"$timestamp"}}, {"count",bson.M{"$sum": 1}}}}}
	sortStage := bson.D{{"$sort", bson.D{{"moment", 1}}}}
	showInfoCursor, err := m.collection("test").Aggregate(ctx, mongo.Pipeline{matchStage,groupStage,sortStage})
	if err != nil {
		log.Println("When aggregate ->", err)
		return Package{}
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		log.Println("When showInfoCursor ->", err)
		return Package{}
	}

	var pack Package

	//Now preparing json
	for k,v:=range showsWithInfo{
		var operator Operator
		operator.Id = k
		operator.Top = k*100
		operator.Info.Numb = v["count"].(int32)
		operator.Info.Description = v["_id"].(string)
		operator.Properties.Title = v["_id"].(string)

		pack.Operators = append(pack.Operators, operator)
	}

	return pack
}

func (m *MongoStructure) InsertRecord(pool string, object string, event string, tme string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, tme)
	if err != nil {
		fmt.Println(err)
		return
	}
	poolId,_:=primitive.ObjectIDFromHex(pool)
	doc:=bson.M{"pool_id":poolId, "object":object, "event":event, "timestamp": t}
	m.collection("test").InsertOne(ctx,doc, nil )
}

func (m *MongoStructure) AddPool(name string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cnt,_ := m.collection("pools").CountDocuments(ctx, bson.M{"name":name})
	if cnt>0{
		return false
	}
	doc:=bson.M{"name":name}
	m.collection("pools").InsertOne(ctx,doc)
	return true
}

func (m *MongoStructure) RemovePool(id string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	idh,_:=primitive.ObjectIDFromHex(id)
	m.collection("pools").DeleteOne(ctx, bson.M{"_id":idh})
	return true
}