package engine

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

type JsonReturn struct {
	Filename string
	Objects  []bson.M
	Events   []bson.M
	MinTime  primitive.DateTime
	MaxTime  primitive.DateTime
}

type Operator struct {
	Id   int `json:"id"`
	Top  int `json:"top"`
	Left int `json:"left"`
	Info struct {
		Numb        int32  `json:"numb"`
		Description string `json:"description"`
	} `json:"info"`
	Properties struct {
		Title   string `json:"title"`
		Outputs struct {
			Outs struct {
				Label string `json:"label"`
			} `json:"outs"`
		} `json:"outputs"`
		Inputs struct {
			Ins struct {
				Label string `json:"label"`
			} `json:"ins"`
		} `json:"inputs"`
	} `json:"properties"`
}

type Link struct {
	FromOperator  int    `json:"fromOperator"`
	FromConnector string `json:"fromConnector"`
	ToOperator    int    `json:"toOperator"`
	ToConnector   string `json:"toConnector"`
	LinkId        string `json:"_id"`
}

type Package struct {
	Operators []Operator `json:"operators"`
	Links     []Link     `json:"links"`
}

func Init() {
	Way = "/var/www/olap/svg/"
	Mongo.DB = "olap"
	Mongo.Opts = options.Find()
	Mongo.OptsDistinct = options.Distinct()
	Mongo.ClientInit()
}
