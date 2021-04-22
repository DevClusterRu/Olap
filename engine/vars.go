package engine

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo MongoStructure

type Operator struct {
	Id int `json:"id"`
	Top int `json:"top"`
	Left int `json:"left"`
	Info struct{
		Numb int32 `json:"numb"`
		Description string `json:"description"`
	} `json:"info"`
	Properties struct{
		Title string `json:"title"`
		Outputs struct{
			Outs struct{
				Label string `json:"label"`
			} `json:"outs"`
		} `json:"outputs"`
		Inputs struct{
			Ins struct{
				Label string `json:"label"`
			} `json:"ins"`
		} `json:"inputs"`
	} `json:"properties"`
}

type Link struct {
	FromOperator string `json:"fromOperator"`
	FromConnector string `json:"fromConnector"`
	ToOperator string `json:"toOperator"`
	ToConnector string `json:"toConnector"`
}


type Package struct {
	Operators []Operator `json:"operators"`
	Links []Link `json:"links"`
}

func Init()  {
	Mongo.DB = "olap"
	Mongo.Opts = options.Find()
	Mongo.OptsDistinct = options.Distinct()
	Mongo.ClientInit()
}


