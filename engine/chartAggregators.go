package engine

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

func (m *MongoStructure) GraphPoolAggregation(query url.Values) AggregationResult {

	var dateStart, dateEnd time.Time
	var err error

	poolId := ""
	if len(query["poolId"]) > 0 {
		poolId = query["poolId"][0]
	}
	if len(query["dataRangeStart"]) > 0 {
		dateStart = StrToDate(query["dataRangeStart"][0])
	}
	if len(query["dataRangeEnd"]) > 0 {
		dateEnd = StrToDate(query["dataRangeEnd"][0])
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var aggregationResult AggregationResult

	groupStage := bson.D{{"$group", bson.D{{"_id", "$object"}}}}
	sortStage := bson.D{{"$sort", bson.D{{"_id", 1}}}}
	distinctObjects, err := m.collection(poolId).Aggregate(ctx, mongo.Pipeline{groupStage, sortStage})
	if err != nil {
		log.Println(err)
	}
	var distinctObj []bson.M
	if err = distinctObjects.All(ctx, &distinctObj); err != nil {
		log.Println("When distinctObjects ->", err)
	}

	//Get total events
	groupStage = bson.D{{"$group", bson.D{{"_id", "$event"}}}}
	sortStage = bson.D{{"$sort", bson.D{{"_id", 1}}}}
	distinctEvents, err := m.collection(poolId).Aggregate(ctx, mongo.Pipeline{groupStage, sortStage})
	if err != nil {
		log.Println("When distinctEvents1 ->", err)
	}

	var distinctEv []bson.M
	if err = distinctEvents.All(ctx, &distinctEv); err != nil {
		log.Println("When distinctEvents2 ->", err)
	}

	//GetMin&Max date
	groupStage = bson.D{{"$group", bson.D{{"_id", bsontype.Null}, {"minDate", bson.M{"$min": "$timestamp"}}, {"maxDate", bson.M{"$max": "$timestamp"}}}}}
	dates, err := m.collection(poolId).Aggregate(ctx, mongo.Pipeline{groupStage, sortStage})
	var datesM []bson.M
	if err = dates.All(ctx, &datesM); err != nil {
		log.Println("When dates ->", err)
	}

	match := bson.D{{"_id", bson.D{{"$exists", 1}}}}

	//Events OR append

	var ifaceE []interface{}
	for _, v := range query["events[]"] {
		ifaceE = append(ifaceE, bson.D{{"event", v}})
	}
	if len(ifaceE) > 0 {
		match = append(match, bson.E{"$or", ifaceE})
	}

	//Cases OR append
	var ifaceC []interface{}
	for _, v := range query["cases[]"] {
		ifaceC = append(ifaceC, bson.D{{"object", v}})
	}
	if len(ifaceC) > 0 {
		match = append(match, bson.E{"$or", ifaceC})
	}

	//Date append
	if len(query["dataRangeStart"]) > 0 && len(query["dataRangeEnd"]) > 0 {
		match = append(match, bson.E{"timestamp", bson.M{"$gte": dateStart, "$lte": dateEnd}})
	}

	nodes, err := m.collection(poolId).Find(ctx, match, &options.FindOptions{Sort: bson.D{{"object", 1}, {"timestamp", 1}}})
	if err != nil {
		fmt.Println(err)
	}
	var data []bson.M
	if err = nodes.All(ctx, &data); err != nil {
		log.Println("When showInfoCursor ->", err)
	}

	fmt.Println("Found ", len(data), " items")

	aggregationResult.Nodes = data
	aggregationResult.NodesCount = int64(len(data))
	aggregationResult.AllObjects = distinctObj
	aggregationResult.AllEvents = distinctEv
	if len(datesM) > 0 {
		aggregationResult.MinDate = datesM[0]["minDate"].(primitive.DateTime)
		aggregationResult.MaxDate = datesM[0]["maxDate"].(primitive.DateTime)
	}

	return BaseCalculation(aggregationResult)

}

type AnswPeriod struct {
	Period []int64  `bson:"_id"`
	Count  []int64 `bson:"total"`
	CasesCount []int64 `bson:"cases"`
}

func prepareForFront(data []bson.M) AnswPeriod {
	var period []int64
	var count []int64
	var cases []int64
	for _, v := range data {

		//Через рефлексию определим типы полученных чисел и инкрементируем счетчики

		typeV:= reflect.TypeOf(v["_id"])
		if typeV==nil{
			period = append(period, 0)
		} else {
			if typeV.Kind() == reflect.Int32 {
				period = append(period, int64(v["_id"].(int32)))
			} else {
				period = append(period, v["_id"].(int64))
			}
		}

		typeV = reflect.TypeOf(v["total"])
		if typeV==nil{
			count = append(count, 0)
		} else {
			if typeV.Kind() == reflect.Int32 {
				count = append(count, int64(v["total"].(int32)))
			} else {
				count = append(count, v["total"].(int64))
			}
		}

		typeV = reflect.TypeOf(v["sizeCases"])
		if typeV==nil{
			cases = append(cases, 0)
		} else {
			if typeV.Kind() == reflect.Int32 {
				cases = append(cases, int64(v["sizeCases"].(int32)))
			} else {
				cases = append(cases, v["sizeCases"].(int64))
			}
		}

	}
	prepared := AnswPeriod{
		Period: period,
		Count:  count,
		CasesCount: cases,
	}
	return prepared
}

func (m *MongoStructure) CubeAggregator(poolId string, year int, month int) AnswPeriod {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	matchStage := bson.D{{"$match", bson.D{{"_id", bson.D{{"$exists", 1}}}}}}

	keyDate:="$year"

	if year!=0 && month==0 {
		y1 := strconv.Itoa(year)
		y2 := strconv.Itoa(year + 1)
		matchStage = bson.D{{"$match", bson.D{{"timestamp", bson.D{{"$gte", StrToDate(y1 + "-01-01")}, {"$lt", StrToDate(y2 + "-01-01")}}}}}}
		keyDate = "$month"
	}

	if year!=0 && month!=0 {
		dateFrom := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-01"
		monthTo := month + 1
		if monthTo > 12 {
			monthTo = 1
			year++
		}
		dateTo := strconv.Itoa(year) + "-" + strconv.Itoa(monthTo) + "-01"
		matchStage = bson.D{{"$match", bson.D{{"timestamp", bson.D{{"$gte", StrToDate(dateFrom)}, {"$lt", StrToDate(dateTo)}}}}}}
		keyDate = "$dayOfMonth"
	}


	groupStage := bson.D{{"$group",
		bson.D{
			{"_id", bson.D{{keyDate, "$timestamp"}}},
			{"total", bson.D{{"$sum", "$numeric"}}},
			{"cases", bson.D{{"$addToSet", "$object"}}},
		},
	}}
	projectStage:=bson.D{{"$project",
		bson.D{
			{"sizeCases", bson.D{{"$size", "$cases"}}},
			{"total", "$total"},
		},
	}}
	sortStage := bson.D{{"$sort", bson.D{{"_id", 1}}}}
	agg, err := m.collection(poolId).Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage, sortStage})

	if err != nil {
		log.Println(err)
	}
	var result []bson.M
	if err = agg.All(ctx, &result); err != nil {
		log.Println("When distinctObjects ->", err)
	}

	fmt.Println(poolId)

	return prepareForFront(result)
}
