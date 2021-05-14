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


func (m *MongoStructure) InsertRecord(poolId string, document bson.M) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	m.collection(poolId).InsertOne(ctx, document, nil)
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

func (m *MongoStructure) RemovePool(poolId string) (*mongo.DeleteResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	poolIdH, _ := primitive.ObjectIDFromHex(poolId)
	m.collection(poolId).Drop(ctx)
	return m.collection("pools").DeleteOne(ctx, bson.M{"_id": poolIdH})
}
