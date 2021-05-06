package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"time"
)

func (m *MongoStructure) GraphPoolAggregation(poolId string) AggregationResult {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	poolIdObject, _ := primitive.ObjectIDFromHex(poolId)
	var aggregationResult AggregationResult

	//Get total distinct objects
	groupStage := bson.D{{"$group", bson.D{{"_id", "$object"}}}}
	sortStage := bson.D{{"$sort", bson.D{{"_id", 1}}}}
	distinctObjects, err := m.collection("test").Aggregate(ctx, mongo.Pipeline{groupStage,sortStage})
	if err!=nil{
		log.Println(err)
	}
	var distinctObj []bson.M
	if err = distinctObjects.All(ctx, &distinctObj); err != nil {
		log.Println("When distinctObjects ->", err)
	}

	//Get total events
	groupStage = bson.D{{"$group", bson.D{{"_id", "$event"}}}}
	sortStage = bson.D{{"$sort", bson.D{{"_id", 1}}}}
	distinctEvents, err := m.collection("test").Aggregate(ctx, mongo.Pipeline{groupStage,sortStage})
	var distinctEv []bson.M
	if err = distinctEvents.All(ctx, &distinctEv); err != nil {
		log.Println("When distinctEvents ->", err)
	}

	nodes, err := m.collection("test").Find(ctx, bson.M{"pool_id":poolIdObject}, &options.FindOptions{Sort: bson.M{"_id": 1}})
	var data []bson.M
	if err = nodes.All(ctx, &data); err != nil {
		log.Println("When showInfoCursor ->", err)
	}

	aggregationResult.Nodes = data
	aggregationResult.NodesCount = int64(len(data))
	aggregationResult.AllObjects = distinctObj
	aggregationResult.AllEvents = distinctEv

	return BaseCalculation(aggregationResult)

}

func (m *MongoStructure) InsertRecord(pool string, object string, event string, tme string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	layout := "2006-01-02 15:04:05"

	fmt.Println("TME: " + tme)

	t, err := time.Parse(layout, strings.TrimSpace(tme))
	if err != nil {
		fmt.Println(err)
		return
	}
	poolId, _ := primitive.ObjectIDFromHex(pool)
	doc := bson.M{"pool_id": poolId, "object": object, "event": event, "timestamp": t}
	m.collection("test").InsertOne(ctx, doc, nil)
}

func (m *MongoStructure) AddPool(name string, description string) string {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cnt, _ := m.collection("pools").CountDocuments(ctx, bson.M{"name": name})
	if cnt > 0 {
		return ""
	}
	doc := bson.M{"name": name, "description": description, "status": "active"}
	res, er := m.collection("pools").InsertOne(ctx, doc)
	if er != nil {
		return "Error: " + er.Error()
	} else {
		return res.InsertedID.(primitive.ObjectID).Hex()
	}

}

func (m *MongoStructure) PoolStatusChange(id string, status string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	poolId, _ := primitive.ObjectIDFromHex(id)
	result, _ := m.collection("pools").UpdateOne(ctx, bson.M{"_id": poolId}, bson.M{"status": status})
	if result == nil {
		return false
	} else {
		return true
	}
}

func (m *MongoStructure) PoolList() string {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, _ := m.collection("pools").Find(ctx, bson.M{})
	var result []bson.M
	err := cursor.All(ctx, &result)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	b, _ := json.Marshal(result)
	return string(b)
}

func (m *MongoStructure) RemovePool(id string) (*mongo.DeleteResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	idh, _ := primitive.ObjectIDFromHex(id)
	return m.collection("pools").DeleteOne(ctx, bson.M{"_id": idh})
	//return true
}

func (m *MongoStructure) CreateLink(poolId string, from string, to string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	idh, _ := primitive.ObjectIDFromHex(poolId)
	doc := bson.M{"pool_id": idh, "link_from": from, "link_to": to}
	_, err := m.collection("links").InsertOne(ctx, doc)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (m *MongoStructure) RemoveLink(linkId string) (*mongo.DeleteResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	idl, _ := primitive.ObjectIDFromHex(linkId)
	return m.collection("links").DeleteOne(ctx, bson.M{"_id": idl})

}
