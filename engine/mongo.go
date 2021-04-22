package engine

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoStructure struct {
	DB     string
	Opts   *options.FindOptions
	OptsDistinct *options.DistinctOptions
	Client *mongo.Client
}

func (m *MongoStructure) ClientInit() {
	var err error
	m.Client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err!=nil{
		log.Println("When try client create -> ", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = m.Client.Connect(ctx)
	fmt.Print("#")
	if err!=nil{
		log.Println("When try client connect -> ", err)
	}
}

func (m *MongoStructure) collection(c string) *mongo.Collection {
	return m.Client.Database(m.DB).Collection(c)
}
