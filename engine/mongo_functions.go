package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/url"
	"time"
)

func (m *MongoStructure) GraphPoolAggregation(query url.Values) AggregationResult {

	var dateStart, dateEnd time.Time
	var err error


	poolId:=""
	if len(query["poolId"])>0 {
		poolId = query["poolId"][0]
	}
	if len(query["dataRangeStart"])>0 {
		dateStart = StrToDate(query["dataRangeStart"][0])
	}
	if len(query["dataRangeEnd"])>0 {
		dateEnd = StrToDate(query["dataRangeEnd"][0])
	}

	fmt.Println(dateStart, dateEnd)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	poolIdObject, _ := primitive.ObjectIDFromHex(poolId)
	var aggregationResult AggregationResult

	matchStage := bson.D{{"$match", bson.D{{"pool_id",poolIdObject}}}}
	//Get total distinct objects
	groupStage := bson.D{{"$group", bson.D{{"_id", "$object"}}}}
	sortStage := bson.D{{"$sort", bson.D{{"_id", 1}}}}
	distinctObjects, err := m.collection("test").Aggregate(ctx, mongo.Pipeline{matchStage,groupStage,sortStage})
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
	distinctEvents, err := m.collection("test").Aggregate(ctx, mongo.Pipeline{matchStage,groupStage,sortStage})
	var distinctEv []bson.M
	if err = distinctEvents.All(ctx, &distinctEv); err != nil {
		log.Println("When distinctEvents ->", err)
	}

	//GetMin&Max date
	groupStage = bson.D{{"$group", bson.D{{"_id", bsontype.Null},{"minDate", bson.M{ "$min": "$timestamp" }},{"maxDate", bson.M{ "$max": "$timestamp" }}}}}
	dates, err := m.collection("test").Aggregate(ctx, mongo.Pipeline{matchStage,groupStage,sortStage})
	var datesM []bson.M
	if err = dates.All(ctx, &datesM); err != nil {
		log.Println("When dates ->", err)
	}



	match:=bson.D{{"pool_id",poolIdObject}}
	//Events OR append

	var ifaceE []interface{}
	for _,v:=range query["events[]"]{
		ifaceE = append(ifaceE,bson.D{{"event", v}})
	}
	if len(ifaceE)>0 {
		match = append(match, bson.E{"$or", ifaceE})
	}

	//Cases OR append
	var ifaceC []interface{}
	for _,v:=range query["cases[]"]{
		ifaceC = append(ifaceC,bson.D{{"object", v}})
	}
	if len(ifaceC)>0 {
		match = append(match, bson.E{"$or", ifaceC})
	}

	//Date append
	if len(query["dataRangeStart"])>0 && len(query["dataRangeEnd"])>0 {
		match = append(match, bson.E{"timestamp", bson.M{"$gte": dateStart, "$lte": dateEnd}})
	}

	nodes, err := m.collection("test").Find(ctx, match, &options.FindOptions{Sort: bson.D{{"object",1}, {"timestamp",1}}})
	if err!=nil{
		fmt.Println(err)
	}
	var data []bson.M
	if err = nodes.All(ctx, &data); err != nil {
		log.Println("When showInfoCursor ->", err)
	}

	fmt.Println("Found ",len(data)," items")

	aggregationResult.Nodes = data
	aggregationResult.NodesCount = int64(len(data))
	aggregationResult.AllObjects = distinctObj
	aggregationResult.AllEvents = distinctEv
	if len(datesM)>0 {
		aggregationResult.MinDate = datesM[0]["minDate"].(primitive.DateTime)
		aggregationResult.MaxDate = datesM[0]["maxDate"].(primitive.DateTime)
	}


	return BaseCalculation(aggregationResult)

}

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
