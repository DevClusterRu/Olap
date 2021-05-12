package engine

import (
	"github.com/araddon/dateparse"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func StrToDate(s string) time.Time  {
	t, _ := dateparse.ParseAny(s)
	return t
}

var Mongo MongoStructure

var Way string

type AggregationResult struct {
	Nodes        []bson.M
	NodesCount   int64
	AllObjects   []bson.M
	AllEvents    []bson.M
	NodesMap     map[string]int
	LinksMap     map[string]int
	NodesSorting []string
	MinDate      primitive.DateTime
	MaxDate      primitive.DateTime
}

type Stream struct {
	Pool string
	Object string
	Event string
	Timestamp string
}

type JsonReturn struct {
	Filename string
	Objects  []bson.M
	Events   []bson.M
	MinTime  primitive.DateTime
	MaxTime  primitive.DateTime
}

func Init()  {
	Way = "/var/www/olap/svg/"
	Mongo.DB = "olap"
	Mongo.Opts = options.Find()
	Mongo.OptsDistinct = options.Distinct()
	Mongo.ClientInit()
}
