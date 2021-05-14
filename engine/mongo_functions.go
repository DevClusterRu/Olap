package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)


func (m *MongoStructure) InsertRecord(document bson.M) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	poolId, _ := primitive.ObjectIDFromHex(document["pool_id"].(string))
	document["pool_id"] = poolId
	//doc := bson.M{"pool_id": poolId, "object": object, "event": event, "timestamp": StrToDate(tme)}
	m.collection("test").InsertOne(ctx, document, nil)
	//pool string, object string, event string, tme string
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
